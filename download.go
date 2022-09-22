package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/grafov/m3u8"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"
)

type DownloadTask struct {
	Req    *http.Request
	Idx    uint64
	Dst    string
	cancel context.CancelFunc
}

func (t *DownloadTask) Start(g *sync.WaitGroup) {
	defer g.Done()
	out, err := os.Create(t.Dst)
	if err != nil {
		runtime.LogError(SharedApp.ctx, err.Error())
		return
	}
	defer func(out *os.File) {
		err = out.Close()
		if err != nil {
			runtime.LogError(SharedApp.ctx, err.Error())
		}
	}(out)

	ctx, cancel := context.WithCancel(SharedApp.ctx)
	req := t.Req.WithContext(ctx)
	t.cancel = cancel

	resp, err := SharedApp.client.Do(req)
	if err != nil {
		runtime.LogError(SharedApp.ctx, err.Error())
		return
	}

	if resp.StatusCode != 200 {
		runtime.LogError(SharedApp.ctx, fmt.Sprintf("Received HTTP %v for %v", resp.StatusCode, t.Req.URL.String()))
		return
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		runtime.LogError(SharedApp.ctx, fmt.Sprintf("Received HTTP %v for %v", resp.StatusCode, t.Req.URL.String()))
		return
	}

	err = resp.Body.Close()
	if err != nil {
		runtime.LogError(SharedApp.ctx, err.Error())
	}
}

func (t *DownloadTask) Stop() {
	if t.cancel != nil {
		t.cancel()
	}
}

type DownloadQueue struct {
	tasks         []*DownloadTask
	TotalDuration float64
	ctx           context.Context
	myKeyIV       string
	Done          chan bool
	DownloadDir   string
}

func (q *DownloadQueue) startDownloadVOD(config *ParserTask, list *m3u8.MediaPlaylist) {
	q.tasks = nil

	for _, seg := range list.Segments {
		if seg != nil {
			q.TotalDuration += seg.Duration
			req, err := config.BuildReq(seg.URI)
			if err != nil {
				runtime.LogError(SharedApp.ctx, fmt.Sprintf("生成 Segments 请求出粗：%v", err))
				continue
			}

			if seg.Key != nil && len(seg.Key.Method) > 0 {

			}

			dst := path.Base(req.URL.Path)
			dst = filepath.Join(q.DownloadDir, dst)
			q.tasks = append(q.tasks, &DownloadTask{
				Req: req,
				Idx: seg.SeqId,
				Dst: dst,
			})
		}
	}

	wg := &sync.WaitGroup{}
	wg.Add(int(list.Count()))

	for _, task := range q.tasks {
		go task.Start(wg)
	}

	wg.Wait()
	q.Done <- true
}

func (q *DownloadQueue) StartDownload(config *ParserTask, list *m3u8.MediaPlaylist) {
	name := config.TaskName
	if len(name) == 0 {
		name = fmt.Sprintf("%v", time.Now().Unix())
	}

	q.DownloadDir = filepath.Join(SharedApp.config.PathDownloader, name)
	if _, err := os.Stat(q.DownloadDir); errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(q.DownloadDir, os.ModePerm)
	}

	if len(q.tasks) > 0 {
		for _, task := range q.tasks {
			task.Stop()
		}
	}

	if list.Closed {
		q.startDownloadVOD(config, list)
	}
}
