package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/grafov/m3u8"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type DownloadTask struct {
	Req *http.Request
	Idx int
}

func (t *DownloadTask) Start(c context.Context) {

}

func (t *DownloadTask) Stop(c context.Context) {

}

type DownloadQueue struct {
	tasks         []*DownloadTask
	TotalDuration float64
	ctx           context.Context
	myKeyIV       string
	Done          chan bool
}

func (q *DownloadQueue) StartDownloadVOD(config *ParserTask, list *m3u8.MediaPlaylist) {
	name := config.TaskName
	if len(name) == 0 {
		name = fmt.Sprintf("%v", time.Now().Unix())
	}

	downloadDir := filepath.Join(SharedApp.config.PathDownloader, name)
	if _, err := os.Stat(downloadDir); errors.Is(err, os.ErrNotExist) {
		os.Mkdir(downloadDir, os.ModePerm)
	}

	if len(q.tasks) > 0 {
		for _, task := range q.tasks {
			task.Stop(q.ctx)
		}
	}

	q.tasks = nil
	for i, seg := range list.Segments {
		if seg != nil {
			q.TotalDuration += seg.Duration
			req, err := config.BuildReq(seg.URI)
			if err != nil {
				runtime.LogError(SharedApp.ctx, fmt.Sprintf("生成 Segments 请求出粗：%v", err))
				continue
			}

			q.tasks = append(q.tasks, &DownloadTask{
				Req: req,
				Idx: i,
			})
		}
	}

	for _, task := range q.tasks {
		println(task)
	}

	q.Done <- true
}
