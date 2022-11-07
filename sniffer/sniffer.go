package sniffer

import (
	"GiveMeAnOffer/eventbus"
	"GiveMeAnOffer/logger"
	"context"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"strings"
	"time"
)

type Sniffer struct {
	Link          string
	resourceLinks map[string]bool
	Cancel        context.CancelFunc

	ctx     context.Context
	logger  logger.AppLogger
	handler eventbus.RuntimeHandler
}

func NewSniffer(link string) *Sniffer {
	return &Sniffer{
		Link: link,
	}
}

func (s *Sniffer) SetupLogger(l logger.AppLogger) {
	s.logger = l
}

func (s *Sniffer) SetupRuntimeHandler(h eventbus.RuntimeHandler) {
	s.handler = h
}

func (s *Sniffer) GetLinks() ([]string, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Flag("headless", true),
		chromedp.Flag("ignore-certificate-custom_error", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(s.ctx, opts...)
	defer cancel()

	taskCtx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(func(ss string, i ...interface{}) {
		s.logger.LogInfof(ss, i)
	}))

	defer cancel()
	s.Cancel = cancel

	// create a timeout
	taskCtx, cancel = context.WithTimeout(taskCtx, 100*time.Second)
	defer cancel()

	// ensure that the browser process is started
	if err := chromedp.Run(taskCtx); err != nil {
		return nil, err
	}

	chromedp.ListenTarget(taskCtx, s.interceptResource(taskCtx))

	s.logger.LogInfof("开始嗅探 URL：", s.Link)

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
	for l := range s.resourceLinks {
		links = append(links, l)
	}
	return links, nil
}

func (s *Sniffer) interceptResource(_ context.Context) func(interface{}) {
	s.resourceLinks = make(map[string]bool)
	suffixes := []string{".m3u8", ".mp4", ".flv", ".mp3", ".mpd", "wav"}
	return func(event interface{}) {
		switch ev := event.(type) {
		case *network.EventResponseReceived:
			s.handler.EventsEmit("intercept-url", ev.Response)
			for _, suffix := range suffixes {
				if strings.Contains(ev.Response.URL, suffix) {
					s.resourceLinks[ev.Response.URL] = true
				}
			}
		default:
			return
		}
	}
}
