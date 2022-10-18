package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/flytam/filenamify"
	"github.com/grafov/m3u8"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type DownloadTask struct {
	ProgressTracker
	Req     *http.Request
	Idx     uint64
	Dst     string
	cancel  context.CancelFunc
	decrypt *Cipher
	Done    bool
}

type ProgressTracker struct {
	io.Reader
	Total   int64
	Current int64
}

func (t *DownloadTask) Read(p []byte) (n int, err error) {
	n, err = t.Reader.Read(p)
	t.Current += int64(n)
	//runtime.LogPrint(SharedApp.ctx, fmt.Sprintf("开始下载: %v", t.Req.URL.String()))
	//runtime.LogInfof(SharedApp.ctx, "\r正在下载，下载进度：%.2f%%", float64(t.Current*10000/t.Total)/100)
	if t.Current == t.Total {
		//runtime.LogInfof(SharedApp.ctx, "\r下载完成，下载进度：%.2f%%", float64(t.Current*10000/t.Total)/100)
	}
	return
}

func (t *DownloadTask) download() error {
	// 解密
	if t.decrypt != nil {
		t.decrypt.Ctx = t.Req.Context()
		err := t.decrypt.Generate()
		if err != nil {
			runtime.LogError(SharedApp.ctx, fmt.Sprintf("创建解密信息失败：%v", err))
			return err
		}
	}

	out, err := os.Create(t.Dst)
	if err != nil {
		runtime.LogError(SharedApp.ctx, err.Error())
		return err
	}
	defer func(out *os.File) {
		err = out.Close()
		if err != nil {
			runtime.LogError(SharedApp.ctx, err.Error())
		}
	}(out)

	resp, err := SharedApp.client.Do(t.Req)
	if err != nil {
		runtime.LogError(SharedApp.ctx, fmt.Sprintf("下载失败：%v", err))
		return err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			runtime.LogError(SharedApp.ctx, err.Error())
		}
	}(resp.Body)

	if resp.StatusCode != 200 {
		info := fmt.Sprintf("下载失败：Received HTTP %v for %v", resp.StatusCode, t.Req.URL.String())
		runtime.LogError(SharedApp.ctx, info)
		return fmt.Errorf(info)
	}

	if t.decrypt != nil {
		buffer, e := t.decrypt.Decrypt(resp.Body)
		if e != nil {
			runtime.LogErrorf(SharedApp.ctx, "解密失败: %v", e.Error())
			return err
		}
		t.Reader = buffer
		t.Total = int64(buffer.Len())
		_, err = io.Copy(out, buffer)
	} else {
		t.Reader = resp.Body
		t.Total = resp.ContentLength
		_, err = io.Copy(out, t)
	}
	runtime.LogPrint(SharedApp.ctx, fmt.Sprintf("下载完成: %v", t.Req.URL.String()))
	if err != nil {
		runtime.LogErrorf(SharedApp.ctx, "Received HTTP %v for %v", resp.StatusCode, t.Req.URL.String())
		return err
	}

	t.Done = true
	return nil
}

func (t *DownloadTask) Start(g *sync.WaitGroup) {
	ctx, cancel := context.WithCancel(SharedApp.ctx)
	t.Req = t.Req.WithContext(ctx)
	t.cancel = cancel
	_ = retry.Do(t.download, retry.Context(ctx), retry.DelayType(func(n uint, config *retry.Config) time.Duration {
		return retry.BackOffDelay(n, config)
	}))

	g.Done()
}

func (t *DownloadTask) Stop() {
	if !t.Done && t.cancel != nil {
		t.cancel()
	}
	t.Done = true
}

type M3U8DownloadQueue struct {
	tasks         []*DownloadTask
	TotalDuration float64
	ctx           context.Context
	DownloadDir   string
	keys          map[string][]byte
	tasksSet      map[uint64]bool
}

func (q *M3U8DownloadQueue) startDownloadVOD(config *ParserTask, list *m3u8.MediaPlaylist) {
	q.tasks = nil

	var cipher *Cipher
	keys := map[string][]byte{}
	var mutex sync.RWMutex

	queryKey := func(u string) ([]byte, bool) {
		mutex.RLock()
		b, ok := keys[u]
		mutex.RUnlock()
		return b, ok
	}
	setKey := func(u string, key []byte) {
		mutex.Lock()
		keys[u] = key
		mutex.Unlock()
	}

	for _, seg := range list.Segments {
		if seg != nil && !q.tasksSet[seg.SeqId] {
			q.TotalDuration += seg.Duration
			req, err := config.BuildReq(seg.URI)
			if err != nil {
				runtime.LogError(SharedApp.ctx, fmt.Sprintf("生成 Segments 请求出粗：%v", err))
				continue
			}

			dst := strconv.Itoa(int(seg.SeqId))
			dst += path.Ext(req.URL.Path)
			dst = filepath.Join(q.DownloadDir, dst)
			task := &DownloadTask{
				Req: req,
				Idx: seg.SeqId,
				Dst: dst,
			}

			decrypt, err := NewCipherFromKey(config, seg.Key, queryKey, setKey)
			if decrypt != nil {
				task.decrypt = decrypt
				if cipher == nil {
					cipher = decrypt
				}
			} else if cipher != nil {
				task.decrypt = cipher
			}
			q.tasks = append(q.tasks, task)
			q.tasksSet[seg.SeqId] = true // 记录下载任务
		}
	}

	// 创建解密
	if cipher != nil {
		err := cipher.Generate()
		if err != nil {
			runtime.LogError(SharedApp.ctx, fmt.Sprintf("生成解密结构体出错：%v", err))
		}
	}

	if len(q.tasks) == 0 {
		return
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(q.tasks))

	for _, task := range q.tasks {
		go task.Start(wg)
	}

	wg.Wait()
	runtime.LogInfof(SharedApp.ctx, "完成了%v个切片任务的下载", len(q.tasks))
}

func (q *M3U8DownloadQueue) startDownloadLive(config *ParserTask, list *m3u8.MediaPlaylist) {
	shouldStop := false
	var tasks []*DownloadTask

	runtime.EventsOn(SharedApp.ctx, StopLiveStreamDownload, func(optionalData ...interface{}) { // 收到停止直播下载的通知
		u := optionalData[0].(string)
		if u != config.Url { // 只停止自己的任务
			return
		}
		shouldStop = true
		for _, task := range q.tasks {
			task.Stop()
		}
	})

	var interval time.Duration = 1 // 重试间隔
	// 直播链接就是不停的分段下载
	for !(shouldStop || list.Closed) {
		q.startDownloadVOD(config, list)
		needWait := len(q.tasks) == 0
		tasks = append(tasks, q.tasks...)
		if needWait { // 如果本次没有新增任务，睡一小会儿
			time.Sleep(interval * time.Second)
			interval += 2
		} else {
			interval = 1
		}
		playlist, _, err := config.retrieveM3U8List()
		if err != nil {
			runtime.LogError(SharedApp.ctx, "获取直播 m3u8 列表失败")
		} else {
			list = playlist.(*m3u8.MediaPlaylist)
			runtime.LogInfof(SharedApp.ctx, "刷新了直播 m3u8 链接")
		}
	}
	// 将所有任务传过去, 后面文件合并用
	q.tasks = tasks
}

func (q *M3U8DownloadQueue) preDownload(config *ParserTask) (err error) {
	name := config.TaskName
	if len(name) == 0 {
		name = fmt.Sprintf("%v", time.Now().Unix())
	}
	q.tasksSet = map[uint64]bool{}
	// 处理非法文件名
	output, err := filenamify.Filenamify(name, filenamify.Options{})
	if err != nil {
		return err
	}
	q.DownloadDir = filepath.Join(SharedApp.config.PathDownloader, output)

	if _, err = os.Stat(q.DownloadDir); errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(q.DownloadDir, os.ModePerm)
	}

	if err != nil {
		return err
	}

	if len(q.tasks) > 0 {
		for _, task := range q.tasks {
			task.Stop()
		}
	}

	item := DownloadTaskUIItem{
		TaskName: config.TaskName,
		Time:     time.Now().Format("2006-01-02 15:04:05"),
		Status:   "初始化...",
		Url:      config.Url,
	}

	runtime.EventsEmit(SharedApp.ctx, TaskNotifyCreate, item)

	return
}

func (q *M3U8DownloadQueue) StartDownload(config *ParserTask, list *m3u8.MediaPlaylist) {
	err := q.preDownload(config)
	if err != nil {
		runtime.LogError(SharedApp.ctx, err.Error())
		return
	}

	if list.Closed {
		q.startDownloadVOD(config, list)
	} else {
		q.startDownloadLive(config, list)
	}
}

type CommonDownloader struct {
	M3U8DownloadQueue
}

type DownloadTaskUIItem struct {
	TaskName string `json:"taskName"`
	Time     string `json:"time"`
	Status   string `json:"status"`
	Url      string `json:"url"`
	IsDone   bool   `json:"isDone"`
}

func (c *CommonDownloader) StartDownload(config *ParserTask, urls []string) error {
	err := c.preDownload(config)
	if err != nil {
		return err
	}

	for idx, u := range urls {
		req, e := http.NewRequest("GET", u, nil)
		for k, v := range config.Headers {
			req.Header.Add(k, v)
		}

		if e != nil {
			return err
		}

		dst := path.Base(req.URL.Path)
		dst = filepath.Join(c.DownloadDir, dst)

		task := &DownloadTask{
			Req: req,
			Idx: uint64(idx),
			Dst: dst,
		}

		c.tasks = append(c.tasks, task)
	}

	if len(c.tasks) == 0 {
		return err
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(c.tasks))

	for _, task := range c.tasks {
		go task.Start(wg)
	}

	wg.Wait()

	return nil
}
