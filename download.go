package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/flytam/filenamify"
	"github.com/grafov/m3u8"
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

type NetworkError struct {
	Code int
	URL  string
}

func (e NetworkError) Error() string {
	return fmt.Sprintf("下载失败：Received HTTP %v for %v", e.Code, e.URL)
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
	defer func() {
		t.Done = true
	}()

	// 解密
	if t.decrypt != nil {
		t.decrypt.Ctx = t.Req.Context()
		err := t.decrypt.Generate()
		if err != nil {
			SharedApp.LogError(fmt.Sprintf("创建解密信息失败：%v", err))
			return err
		}
	}

	out, err := os.Create(t.Dst)
	if err != nil {
		SharedApp.LogError(err.Error())
		return err
	}
	defer func(out *os.File) {
		err = out.Close()
		if err != nil {
			SharedApp.LogError(err.Error())
		}
	}(out)

	resp, err := SharedApp.client.Do(t.Req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			SharedApp.LogError(err.Error())
		}
	}(resp.Body)

	if resp.StatusCode != 200 {
		return NetworkError{
			Code: resp.StatusCode,
			URL:  t.Req.URL.String(),
		}
	}

	if t.decrypt != nil {
		buffer, e := t.decrypt.Decrypt(resp.Body)
		if e != nil {
			return e
		}
		t.Reader = buffer
		t.Total = int64(buffer.Len())
		_, err = io.Copy(out, buffer)
	} else {
		t.Reader = resp.Body
		t.Total = resp.ContentLength
		_, err = io.Copy(out, t)
	}
	SharedApp.LogInfof("下载完成: %v", t.Req.URL.String())
	return err
}

func (t *DownloadTask) Start(g *sync.WaitGroup) error {
	ctx, cancel := context.WithCancel(SharedApp.ctx)
	t.Req = t.Req.WithContext(ctx)
	t.cancel = cancel
	err := retry.Do(t.download,
		retry.Context(ctx),
		retry.RetryIf(func(err error) bool {
			_, ok := err.(NetworkError)
			if ok {
				e := err.(NetworkError)
				return e.Code < 400 || e.Code >= 500
			}
			return false
		}),
		retry.Attempts(3), // 重试三次
		retry.DelayType(func(n uint, err error, config *retry.Config) time.Duration {
			return 3 * time.Second
		}))

	g.Done()
	return err
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
	concurrentCnt int
}

func (q *M3U8DownloadQueue) startDownloadVOD(config *ParserTask, list *m3u8.MediaPlaylist) error {
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
				SharedApp.LogErrorf("生成 Segments 请求出粗：%v", err)
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
			SharedApp.LogErrorf("生成解密结构体出错：%v", err)
			return err
		}
	}

	if len(q.tasks) == 0 {
		return nil
	}

	wg := &sync.WaitGroup{}

	var err error
	ch := make(chan struct{}, q.concurrentCnt)
	for _, task := range q.tasks {
		ch <- struct{}{}
		wg.Add(1)
		go func(t *DownloadTask) {
			err = t.Start(wg)
			if err != nil { // 出现了错误，直接停掉其他任务，结束
				SharedApp.LogError(err.Error())
				q.Stop()
			}
			<-ch
		}(task)
	}

	wg.Wait()
	SharedApp.LogInfof("完成了%v个切片任务的下载", len(q.tasks))
	return err
}

func (q *M3U8DownloadQueue) Stop() {
	if _, err := os.Stat(q.DownloadDir); errors.Is(err, os.ErrNotExist) {
		return
	}
	// 停掉所有任务
	for _, task := range q.tasks {
		task.Stop()
	}
	_ = os.RemoveAll(q.DownloadDir)
}

func (q *M3U8DownloadQueue) startDownloadLive(config *ParserTask, list *m3u8.MediaPlaylist) error {
	shouldStop := false
	var tasks []*DownloadTask

	SharedApp.EventsOn(StopLiveStreamDownload, func(optionalData ...interface{}) { // 收到停止直播下载的通知
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
		err := q.startDownloadVOD(config, list)
		if err != nil {
			return err
		}
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
			SharedApp.LogError("获取直播 m3u8 列表失败")
		} else {
			list = playlist.(*m3u8.MediaPlaylist)
			SharedApp.LogInfo("刷新了直播 m3u8 链接")
		}
	}
	// 将所有任务传过去, 后面文件合并用
	q.tasks = tasks

	return nil
}

func (q *M3U8DownloadQueue) preDownload(config *ParserTask) (err error) {
	name := config.TaskName
	q.concurrentCnt = 10
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

	SharedApp.EventsEmit(TaskNotifyCreate, item)

	return
}

func (q *M3U8DownloadQueue) StartDownload(config *ParserTask, list *m3u8.MediaPlaylist) error {
	err := q.preDownload(config)
	if err != nil {
		return err
	}

	if list.Closed {
		return q.startDownloadVOD(config, list)
	} else {
		return q.startDownloadLive(config, list)
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
	ch := make(chan struct{}, c.concurrentCnt)
	for _, task := range c.tasks {
		wg.Add(1)
		go func(t *DownloadTask) {
			defer wg.Done()
			err = t.Start(wg)
			if err != nil {
				SharedApp.LogError(err.Error())
			}
			<-ch
		}(task)
	}

	wg.Wait()
	return nil
}
