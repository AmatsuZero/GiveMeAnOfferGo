package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"net/http"
	"regexp"
	"strings"
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
		SharedApp.LogInfof(s, i)
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

	// 资源链接有效时间很短，每个页面单独提取，然后立刻下载
	for _, task := range tasks {
		timeoutCtx, cancel := createChromeContext()
		// 提取 m3u8 链接
		chromedp.ListenTarget(timeoutCtx, t.interceptResource(task))
		e := chromedp.Run(timeoutCtx,
			t.setCookies(),
			network.Enable(),
			chromedp.Navigate(task.Url),
			chromedp.WaitVisible("#catalog", chromedp.ByID),
		)
		if e != nil {
			SharedApp.LogErrorf("❌获取链接失败：%v", e)
		}
		cancel()
		// 下载
		e = task.Parse()
		if e != nil {
			SharedApp.LogErrorf("❌下载课程失败：%v", task.TaskName)
		}
	}
	return err
}

func (t *ChinaAACCParserTask) interceptResource(task *ParserTask) func(interface{}) {
	return func(event interface{}) {
		var keyRequestID network.RequestID

		switch ev := event.(type) {
		case *network.EventRequestWillBeSent:
			if strings.Contains(ev.Request.URL, ".m3u8") {
				task.Url = ev.Request.URL
				SharedApp.LogInfof("提取到网校 m3u8 资源链接：%v")
			} else if strings.Contains(ev.Request.URL, "getKeyForHls") { // 获取密钥的链接
				keyRequestID = ev.RequestID
			}
		case *network.EventResponseReceived:
			if len(keyRequestID) > 0 && ev.RequestID == keyRequestID {
				b, err := network.GetResponseBody(keyRequestID).MarshalJSON()
				if err != nil {
					SharedApp.LogError(err.Error())
				} else {
					task.KeyIV = string(b)
				}
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
		task := &ParserTask{
			Headers:       t.Headers,
			DelOnComplete: t.DelOnComplete,
			Prefix:        t.Prefix,
			KeyIV:         t.KeyIV,
		}
		u, ok := selection.Attr("href")
		if !ok {
			return
		}

		task.Url = "https://www.chinaacc.com" + u

		reg := regexp.MustCompile("[\u4e00-\u9fa5]+") //我们要匹配中文的匹配规则
		name := reg.FindAllString(selection.Text(), -1)
		if len(name) > 0 {
			task.TaskName = name[0]
		} else {
			task.TaskName = fmt.Sprintf("%v", time.Now().Unix())
		}
		tasks = append(tasks, task)
	})

	return
}
