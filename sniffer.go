package main

import (
	"context"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"path"
	"time"
)

type Sniffer struct {
	Link          string
	resourceLinks map[string]bool
	Cancel        context.CancelFunc
}

func NewSniffer(link string) *Sniffer {
	return &Sniffer{
		Link: link,
	}
}

func (s *Sniffer) GetLinks() ([]string, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Flag("headless", true),
		chromedp.Flag("ignore-certificate-errors", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(SharedApp.ctx, opts...)
	defer cancel()

	taskCtx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()
	s.Cancel = cancel

	// create a timeout
	taskCtx, cancel = context.WithTimeout(taskCtx, 100*time.Second)
	defer cancel()

	// ensure that the browser process is started
	if err := chromedp.Run(taskCtx); err != nil {
		return nil, err
	}

	chromedp.ListenTarget(taskCtx, s.interceptM3u8(taskCtx))

	runtime.LogInfof(SharedApp.ctx, "开始嗅探 URL：", s.Link)

	err := chromedp.Run(taskCtx,
		network.Enable(),
		chromedp.Navigate(s.Link),
		chromedp.WaitVisible(`body`, chromedp.BySearch),
	)

	if err != nil {
		return nil, err
	}

	// 去重
	var links []string
	for l, _ := range s.resourceLinks {
		links = append(links, l)
	}
	return links, nil
}

func (s *Sniffer) interceptM3u8(ctx context.Context) func(interface{}) {
	s.resourceLinks = make(map[string]bool)

	return func(event interface{}) {
		switch ev := event.(type) {
		case *network.EventResponseReceived:
			ext := path.Ext(ev.Response.URL)
			if ext == ".m3u8" {
				s.resourceLinks[ev.Response.URL] = true
			}
		}
	}
}
