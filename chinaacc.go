package main

import (
	"context"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

type ChinaAACCParserTask struct {
	*ParserTask
	cookies []*http.Cookie
}

func (t *ChinaAACCParserTask) getCookies() error {
	rawCookies, ok := t.Headers["Cookie"]
	if !ok {
		return errors.New("需要设置Cookie")
	}

	header := http.Header{}
	header.Add("Cookie", rawCookies)
	request := http.Request{Header: header}
	t.cookies = request.Cookies()

	return nil
}

func createChromeContext() (context.Context, context.CancelFunc) {
	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", true), // debug使用
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
		chromedp.NoDefaultBrowserCheck,
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36`),
	}
	//初始化参数，先传一个空的数据
	options = append(chromedp.DefaultExecAllocatorOptions[:], options...)

	c, _ := chromedp.NewExecAllocator(SharedApp.ctx, options...)

	// create context
	chromeCtx, _ := chromedp.NewContext(c, chromedp.WithLogf(func(s string, i ...interface{}) {
		runtime.LogInfof(SharedApp.ctx, s, i)
	}))

	//创建一个上下文，超时时间为40s
	return context.WithTimeout(chromeCtx, 40*time.Second)
}

func (t *ChinaAACCParserTask) Parse() error {
	// 获取cookie
	err := t.getCookies()
	if err != nil {
		return err
	}

	content, err := t.getHttpHtmlContent()
	if err != nil {
		return err
	}

	tasks, err := t.scrapeLinks(content)
	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}

	// 获取所有 m3u8 链接
	for _, task := range tasks {
		wg.Add(1)
		go func(g *sync.WaitGroup, tt *ParserTask) {
			timeoutCtx, cancel := createChromeContext()
			defer cancel()
			chromedp.ListenTarget(timeoutCtx, t.interceptResource(tt))
			e := chromedp.Run(timeoutCtx,
				t.setCookies(),
				network.Enable(),
				chromedp.Navigate(tt.Url),
				chromedp.WaitVisible("#catalog", chromedp.ByID),
			)
			if e != nil {
				runtime.LogErrorf(SharedApp.ctx, "获取链接失败：%v", e)
			}
			wg.Done()
		}(wg, task)
	}
	wg.Wait()
	return SharedApp.TaskAddMuti(tasks)
}

func (t *ChinaAACCParserTask) interceptResource(task *ParserTask) func(interface{}) {
	return func(event interface{}) {
		switch ev := event.(type) {
		case *network.EventRequestWillBeSent:
			if strings.Contains(ev.Request.URL, ".m3u8") {
				task.Url = ev.Request.URL
			}
		}
	}
}

func (t *ChinaAACCParserTask) setCookies() chromedp.ActionFunc {
	return func(ctx context.Context) error {
		for _, cookie := range t.cookies {
			err := network.SetCookie(cookie.Name, cookie.Value).
				WithPath("/").
				WithDomain(".chinaacc.com").
				Do(ctx)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

// 获取网站上爬取的数据
func (t *ChinaAACCParserTask) getHttpHtmlContent() (string, error) {
	var htmlContent string
	timeoutCtx, cancel := createChromeContext()
	defer cancel()
	err := chromedp.Run(timeoutCtx,
		t.setCookies(),
		chromedp.Navigate(t.Url),
		chromedp.WaitVisible("#catalog", chromedp.ByID),
		chromedp.OuterHTML("#list1", &htmlContent, chromedp.ByID),
	)
	if err != nil {
		return "", err
	}

	return htmlContent, nil
}

func (t *ChinaAACCParserTask) scrapeLinks(htmlContent string) (tasks []*ParserTask, err error) {
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}

	dom.Find("#list1").Children().Each(func(i int, selection *goquery.Selection) {
		task := &ParserTask{Headers: t.Headers}
		u, ok := selection.Attr("href")
		if !ok {
			return
		}

		task.Url = "https://www.chinaacc.com" + u

		reg := regexp.MustCompile("[\u4e00-\u9fa5]{1,}") //我们要匹配中文的匹配规则
		name := reg.FindAllString(selection.Text(), -1)
		if len(name) > 0 {
			task.TaskName = name[0]
		}

		tasks = append(tasks, task)
	})

	return
}
