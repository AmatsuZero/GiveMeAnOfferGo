package main

import (
	"bufio"
	"fmt"
	"github.com/grafov/m3u8"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	TaskAddEvent      = "task-add-reply"
	SelectVariant     = "select-variant"
	OnVariantSelected = "variant-selected"
)

type ParserTask struct {
	Url           string            `json:"url"`
	TaskName      string            `json:"taskName"`
	Prefix        string            `json:"prefix"`
	DelOnComplete bool              `json:"delOnComplete"`
	KeyIV         string            `json:"keyIV"`
	Headers       map[string]string `json:"headers"`
}

type EventMessage struct {
	Code    int
	Message string
	Info    interface{}
}

func (t *ParserTask) Parse() error {
	u, err := url.Parse(t.Url)
	if err != nil {
		return err
	}

	isLocal := u.Scheme == "http" || u.Scheme == "https"

	if !isLocal {
		err = t.parseLocalFile(u.String())
		if err != nil {
			return err
		}
	} else {
		err = t.getPlayerList()
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *ParserTask) parseLocalFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	p, listType, err := m3u8.DecodeFrom(bufio.NewReader(f), true)
	if err != nil {
		panic(err)
	}
	switch listType {
	case m3u8.MEDIA:
		t.handleMediaPlayList(p.(*m3u8.MediaPlaylist))
	case m3u8.MASTER:
		t.selectVariant(p.(*m3u8.MasterPlaylist))
	}
	return nil
}

func (t *ParserTask) buildSegmentsURL(u string) (string, error) {
	if strings.HasPrefix(u, "http") {
		return url.QueryUnescape(u)
	} else {
		playlistUrl, err := url.Parse(t.Url)
		msUrl, err := playlistUrl.Parse(u)
		if err != nil {
			return "", err
		}
		return url.QueryUnescape(msUrl.String())
	}
}

func (t *ParserTask) BuildReq(u string) (*http.Request, error) {
	if u != t.Url {
		var err error
		u, err = t.buildSegmentsURL(u)
		if err != nil {
			return nil, err
		}
	}
	playlistUrl, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", playlistUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	origin := t.Headers["Origin"]
	if len(origin) == 0 {
		origin = t.Headers["origin"]
	}
	if len(origin) == 0 {
		origin = playlistUrl.Host
	}
	req.Header.Set("Origin", origin)

	refer := t.Headers["Referer"]
	if len(refer) == 0 {
		refer = t.Headers["referer"]
	}
	if len(refer) == 0 {
		refer = playlistUrl.Host
	}
	req.Header.Set("Referer", refer)

	req = req.WithContext(SharedApp.ctx)
	return req, nil
}

func (t *ParserTask) selectVariant(l *m3u8.MasterPlaylist) {
	// 等待前端选择
	msg := EventMessage{
		Code:    1,
		Message: "",
	}
	playlist := map[string]int{}
	for i, variant := range l.Variants {
		playlist[variant.Resolution] = i
	}
	msg.Info = playlist
	runtime.EventsEmit(SharedApp.ctx, SelectVariant, msg)
	runtime.EventsOnce(SharedApp.ctx, OnVariantSelected, func(optionalData ...interface{}) {
		res := optionalData[0].(string)
		idx := playlist[res]
		t.handleVariant(l.Variants[idx])
	})
}

func (t *ParserTask) handleVariant(v *m3u8.Variant) {
	if v.Chunklist != nil {
		t.handleMediaPlayList(v.Chunklist)
		return
	}
	req, err := t.BuildReq(v.URI)
	if err != nil {
		runtime.LogError(SharedApp.ctx, err.Error())
		return
	}

	t.Url = req.URL.String()
	err = t.Parse()
	if err != nil {
		runtime.LogError(SharedApp.ctx, err.Error())
	}
}

func (t *ParserTask) handleMediaPlayList(mpl *m3u8.MediaPlaylist) {
	cnt := 0
	info := ""

	queue := &DownloadQueue{}

	if mpl.Closed {
		go queue.StartDownload(t, mpl)
		d := time.Unix(int64(queue.TotalDuration), 0).Format("15:07:51")
		info = fmt.Sprintf("点播资源解析成功，有%v个片段，时长：%v，，即将开始缓存...", cnt, d)
	} else {
		info = "直播资源解析成功，即将开始缓存..."
	}

	runtime.EventsEmit(SharedApp.ctx, TaskAddEvent, EventMessage{
		Code:    0,
		Message: info,
	})

	<-queue.Done
}

func (t *ParserTask) getPlayerList() error {
	req, err := t.BuildReq(t.Url)
	if err != nil {
		return err
	}

	resp, err := SharedApp.client.Do(req)
	if err != nil {
		return err
	}

	playlist, listType, err := m3u8.DecodeFrom(resp.Body, true)
	if err != nil {
		return err
	}

	err = resp.Body.Close()
	if err != nil {
		return err
	}

	if listType == m3u8.MASTER {
		t.selectVariant(playlist.(*m3u8.MasterPlaylist))
	} else {
		t.handleMediaPlayList(playlist.(*m3u8.MediaPlaylist))
	}
	return nil
}
