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
	TaskAddEvent = "task-add-reply"
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
	Code      int
	Message   string
	PlayLists []string
}

func (t *ParserTask) Parse() error {
	u, err := url.Parse(t.Url)
	if err != nil {
		return err
	}

	isLocal := u.Scheme == "http" || u.Scheme == "https"

	if !isLocal {
		t.parseLocalFile(u.String())
	} else {
		t.getPlayerList()
	}

	return nil
}

func (t *ParserTask) parseLocalFile(path string) error {
	f, err := os.Open("playlist.m3u8")
	if err != nil {
		return err
	}
	p, listType, err := m3u8.DecodeFrom(bufio.NewReader(f), true)
	if err != nil {
		panic(err)
	}
	switch listType {
	case m3u8.MEDIA:
		mediapl := p.(*m3u8.MediaPlaylist)
		fmt.Printf("%+v\n", mediapl)
	case m3u8.MASTER:
		masterpl := p.(*m3u8.MasterPlaylist)
		fmt.Printf("%+v\n", masterpl)
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

	return req, nil
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
	resp.Body.Close()

	if listType == m3u8.MASTER {

	} else {
		mpl := playlist.(*m3u8.MediaPlaylist)
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
			Code:    1,
			Message: info,
		})

		<-queue.Done
	}

	return nil
}
