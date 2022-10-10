package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
)

type ChinaAACCParserTask struct {
	*ParserTask
}

func (t *ChinaAACCParserTask) Parse() error {
	req, err := http.NewRequest("GET", t.Url, nil)
	if err != nil {
		return err
	}

	if _, ok := t.Headers["Referer"]; !ok {
		t.Headers["Referer"] = t.Url
	}

	for k, v := range t.Headers {
		req.Header.Add(k, v)
	}

	// 先获取页面
	res, err := SharedApp.client.Do(req)
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// 获取链接
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}

	t.scrapeLinks(doc)
	println(doc.Text())
	return nil
}

func (t *ChinaAACCParserTask) scrapeLinks(doc *goquery.Document) (tasks []ParserTask) {
	doc.Find("#list1").Children().Each(func(i int, selection *goquery.Selection) {
		task := ParserTask{Headers: t.Headers}
		link := selection.Find("#list1 > a:nth-child(1)")
		u, ok := link.Attr("href")
		if !ok {
			return
		}
		task.Url = u

		name := selection.Find("#list1 > a:nth-child(1) > span.fl").Text()
		task.TaskName = name

		tasks = append(tasks, task)
	})

	return
}
