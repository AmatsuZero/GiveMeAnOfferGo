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
	Req     *http.Request
	Idx     uint64
	Dst     string
	cancel  context.CancelFunc
	decrypt *Cipher
}

func (t *DownloadTask) Start(g *sync.WaitGroup) {
	defer func() {
		g.Done()
	}()

	// 解密
	if t.decrypt != nil {
		t.decrypt.Ctx = t.Req.Context()
		err := t.decrypt.Generate()
		if err != nil {
			runtime.LogError(SharedApp.ctx, fmt.Sprintf("创建解密信息失败：%v", err))
			return
		}
	}

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
		runtime.LogError(SharedApp.ctx, fmt.Sprintf("下载失败：%v", err))
		return
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			runtime.LogError(SharedApp.ctx, err.Error())
		}
	}(resp.Body)

	if resp.StatusCode != 200 {
		runtime.LogErrorf(SharedApp.ctx, "下载失败：Received HTTP %v for %v", resp.StatusCode, t.Req.URL.String())
		return
	}

	buffer, err := t.decrypt.Decrypt(resp.Body)
	if err != nil {
		runtime.LogError(SharedApp.ctx, "解密失败")
		return
	}
	_, err = io.Copy(out, buffer)
	if err != nil {
		runtime.LogErrorf(SharedApp.ctx, "Received HTTP %v for %v", resp.StatusCode, t.Req.URL.String())
		return
	}
	runtime.LogPrint(SharedApp.ctx, fmt.Sprintf("下载完成: %v", req.URL.String()))
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
	Done          chan bool
	DownloadDir   string

	keys map[string][]byte
}

func (q *DownloadQueue) startDownloadVOD(config *ParserTask, list *m3u8.MediaPlaylist) {
	q.tasks = nil

	var cipher *Cipher
	keys := map[string][]byte{}
	queryKey := func(u string) ([]byte, bool) {
		return keys[u], false
	}
	setKey := func(u string, key []byte) {
		keys[u] = key
	}

	for _, seg := range list.Segments {
		if seg != nil {
			q.TotalDuration += seg.Duration
			req, err := config.BuildReq(seg.URI)
			if err != nil {
				runtime.LogError(SharedApp.ctx, fmt.Sprintf("生成 Segments 请求出粗：%v", err))
				continue
			}

			dst := path.Base(req.URL.Path)
			dst = filepath.Join(q.DownloadDir, dst)
			task := &DownloadTask{
				Req: req,
				Idx: seg.SeqId,
				Dst: dst,
			}

			decrypt, err := NewCipherFromKey(seg.Key, config.KeyIV, queryKey, setKey)
			if decrypt != nil {
				task.decrypt = decrypt
				if cipher == nil {
					cipher = decrypt
				}
			} else if cipher != nil {
				task.decrypt = cipher
			}
			q.tasks = append(q.tasks, task)
		}
	}

	// 创建解密
	if cipher != nil {
		err := cipher.Generate()
		if err != nil {
			runtime.LogError(SharedApp.ctx, fmt.Sprintf("生成解密结构体出错：%v", err))
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

	q.Done = make(chan bool)
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
