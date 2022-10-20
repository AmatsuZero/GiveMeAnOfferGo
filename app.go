package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/grafov/m3u8"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/skratchdot/open-golang/open"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var configFilePath string

type AppCtxKey string

func init() {
	configFilePath, _ = os.UserConfigDir()
	if len(configFilePath) == 0 {
		configFilePath = os.Getenv("APPDATA")
	}
	configFilePath = filepath.Join(configFilePath, "M3U8-Downloader-GO")

	if _, err := os.Stat(configFilePath); errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(configFilePath, os.ModePerm)
		if err != nil {
			return
		}
	}

	configFilePath = filepath.Join(configFilePath, "config.json")
}

// App struct
type App struct {
	config         *UserConfig
	ctx            context.Context
	client         *http.Client
	stopTasks      context.CancelFunc
	sniffer        *Sniffer
	concurrentLock chan struct{}

	tasks  map[string]*DownloadTaskUIItem
	queues map[string]StoppableTask
}

// NewApp creates a new App application struct
func NewApp() *App {
	tr := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   0,
		ExpectContinueTimeout: 1 * time.Second,
	}

	return &App{
		client: &http.Client{
			Transport: tr,
		},
	}
}

func (a *App) startup(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	a.ctx, a.stopTasks = ctx, cancel

	config, err := NewConfig(configFilePath)
	if err != nil {
		a.logError(err.Error())
	} else if config.ConfigProxy != nil {
		// 写入代理配置
		err = os.Setenv("HTTP_PROXY", config.ConfigProxy.http)
		if err != nil {
			a.logError(err.Error())
		}
		err = os.Setenv("HTTPS_PROXY", config.ConfigProxy.https)
		if err != nil {
			a.logError(err.Error())
		}
	}
	a.config = config
	a.concurrentLock = make(chan struct{}, config.ConCurrentCnt)
	a.tasks = map[string]*DownloadTaskUIItem{}
	a.queues = map[string]StoppableTask{}
}

func (a *App) shutdown(ctx context.Context) {
	err := a.config.Save()
	if err != nil {
		runtime.LogError(ctx, err.Error())
	}
	a.stopTasks()
}

func (a *App) OpenSelectM3U8() (string, error) {
	return runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "请选择一个M3U8文件",
	})
}

func (a *App) OpenSelectTsDir(dir string) ([]string, error) {
	if len(dir) > 0 {
		var files []string
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			ext := filepath.Ext(path)
			if ext == ".ts" || ext == ".TS" {
				files = append(files, path)
			}
			return nil
		})

		if err != nil {
			return nil, err
		}
	}

	return runtime.OpenMultipleFilesDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "请选择欲合并的TS文件",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "视频片段", Pattern: "*.ts",
			},
			{
				DisplayName: "所有文件", Pattern: "*",
			},
		},
	})
}

func (a *App) StartMergeTs(config MergeFilesConfig) error {
	_, e := config.Merge()
	return e
}

func (a *App) OpenConfigDir() (string, error) {
	defaultDir := a.config.PathDownloader
	if len(defaultDir) == 0 {
		base, _ := os.UserHomeDir()
		defaultDir = filepath.Join(base, "Downloads")
	}

	dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title:                "请选择文件夹",
		DefaultDirectory:     defaultDir,
		CanCreateDirectories: true,
	})

	if err != nil {
		return "", err
	}

	a.config.PathDownloader = dir
	return dir, err
}

func (a *App) TaskAdd(task *ParserTask) error {
	ret, err := task.Parse()
	if err != nil {
		return err
	}
	switch ret.Type {
	case TaskTypeCommon:
		err = a.handleCommonTask(task, ret)
	case TaskTypeChinaAACC:
		err = a.handleAACCTask(ret)
	case TaskTypeM3U8:
		err = a.handleM3U8Task(task, ret)
	case TaskTypeM3U:
		err = a.handleM3UTask(ret)
	case TaskTypeBilibili:
		err = a.handleBilibiliTask(task, ret)
	}
	return err
}

func (a *App) handleBilibiliTask(task *ParserTask, result *ParseResult) (err error) {
	info := result.Data.(*videoInfoResp)

	task.TaskName = info.Data.Title

	u := baseURl.JoinPath(videoStream)
	values := u.Query()
	if len(info.Data.Bvid) > 0 {
		values.Add("bvid", info.Data.Bvid)
	} else {
		values.Add("aid", strconv.Itoa(int(info.Data.aid)))
	}
	u.RawQuery = values.Encode()

	if len(info.Data.Pages) == 0 {
		return errors.New("no page data")
	}

	res, err := info.Data.Pages[0].selectResolution(u)
	if err != nil {
		return err
	}

	item := DownloadTaskUIItem{
		TaskName: task.TaskName,
		Time:     time.Now().Format("2006-01-02 15:04:05"),
		Status:   "初始化...",
		Url:      task.Url,
	}

	a.eventsEmit(TaskNotifyCreate, item)

	go func() {
		wg := &sync.WaitGroup{}
		wg.Add(int(info.Data.Videos))

		for _, page := range info.Data.Pages {
			go func(p *videoPageData) {
				tmp, e := u.Parse(u.String())
				if e != nil {
					return
				}
				vars := tmp.Query()
				vars.Add("qn", res)
				tmp.RawQuery = vars.Encode()
				err = p.download(tmp, wg, task)
				if err != nil {
					a.logInfof("B站任务下载失败：%v", err)
				}
			}(page)
		}
		wg.Wait()
	}()

	return
}

func (a *App) handleM3UTask(result *ParseResult) (err error) {
	tasks := result.Data.([]*ParserTask)
	for _, task := range tasks {
		err = a.TaskAdd(task)
		if err != nil {
			a.logErrorf("m3u 列表任务下载失败，链接：%v", task.Url)
		}
	}
	return
}

func (a *App) addTaskNotifyItem(task *ParserTask) *DownloadTaskUIItem {
	// 先从记录里面查找
	item, ok := a.tasks[task.Url]
	if ok {
		return item
	}

	// 通知前端任务列表添加任务
	item = &DownloadTaskUIItem{
		TaskName: task.TaskName,
		Time:     time.Now().Format("2006-01-02 15:04:05"),
		Status:   "初始化...",
		Url:      task.Url,
		TaskID:   len(a.tasks) + 1,
	}
	a.eventsEmit(TaskNotifyCreate, item)
	return item
}

func (a *App) RemoveTaskNotifyItem(item *DownloadTaskUIItem) error {
	if !item.IsDone {
		task, ok := a.queues[item.Url]
		if ok {
			task.Stop()
		}
	}
	err := os.Remove(item.VideoPath)
	delete(a.tasks, item.Url)
	return err
}

func (a *App) handleM3U8Task(task *ParserTask, result *ParseResult) (err error) {
	mpl := result.Data.(*m3u8.MediaPlaylist)
	queue := &M3U8DownloadQueue{}
	a.queues[task.Url] = queue

	info, cnt := "", 0
	if mpl.Closed {
		d := time.Unix(int64(queue.TotalDuration), 0).Format("15:07:51")
		info = fmt.Sprintf("点播资源解析成功，有%v个片段，时长：%v，，即将开始缓存...", cnt, d)
	} else {
		info = "直播资源解析成功，即将开始缓存..."
	}
	// 通知前端任务即将开始
	a.eventsEmit(TaskAddEvent, EventMessage{
		Code:    0,
		Message: info,
	})

	// 任务列表添加任务
	item := a.addTaskNotifyItem(task)
	queue.NotifyItem = item

	go func() {
		defer delete(a.queues, task.Url)

		// 开始下载
		a.concurrentLock <- struct{}{}
		err = queue.StartDownload(task, mpl)
		// 下载完毕，释放资源
		<-a.concurrentLock
		if err != nil {
			item.Status = "下载失败，请检查链接有效性"
			a.eventsEmit(TaskFinish, item)
			return
		}

		item.Status = "已完成，合并中..."
		a.eventsEmit(TaskFinish, item)
		a.logInfof("切片下载完成，一共%v个", len(queue.tasks))

		merger := NewMergeConfigFromDownloadQueue(queue, task.TaskName)
		output, err := merger.Merge()
		if err != nil {
			item.Status = "合并出错，请尝试手动合并"
			a.eventsEmit(TaskFinish, item)
			return
		}

		a.logInfo("切片合并完成")
		if task.DelOnComplete {
			err = os.RemoveAll(queue.DownloadDir)
			a.logInfo("切片删除完成")
		}

		item.Status = "已完成"
		item.IsDone = true
		item.VideoPath = output
		a.eventsEmit(TaskFinish, item)
	}()
	return
}

func (a *App) handleAACCTask(result *ParseResult) error {
	ch := result.Data.(chan *ParserTask)
	for parserTask := range ch {
		err := a.TaskAdd(parserTask)
		if err != nil {
			a.logErrorf("正保网校课程下载失败，任务名：%v，链接：%v", parserTask.TaskName, parserTask.Url)
		}
	}
	return nil
}

func (a *App) handleCommonTask(task *ParserTask, result *ParseResult) (err error) {
	item := a.addTaskNotifyItem(task)

	go func() {
		a.concurrentLock <- struct{}{}
		q := &CommonDownloader{}
		a.queues[task.Url] = q
		q.NotifyItem = item
		err = q.StartDownload(task, result.Data.([]string))
		<-a.concurrentLock
	}()

	return
}

func (a *App) TaskAddMuti(tasks []*ParserTask) error {
	var wg sync.WaitGroup
	for _, task := range tasks {
		wg.Add(1)
		go func(t *ParserTask) {
			defer wg.Done()
			e := a.TaskAdd(t)
			if e != nil {
				a.logErrorf("下载任务失败:%v, 原因：%v", t.Url, e.Error())
			}

		}(task)
	}
	wg.Wait()
	return nil
}

func (a *App) SniffLinks(u string) ([]string, error) {
	if a.sniffer != nil && a.sniffer.Cancel != nil {
		a.sniffer.Cancel()
	}
	a.sniffer = NewSniffer(u)
	return a.sniffer.GetLinks()
}

func (a *App) Open(link string) error {
	if len(link) == 0 {
		link = a.config.PathDownloader
	}
	return open.Run(link)
}

func (a *App) Play(file string) error {
	msg, err := Cmd("ffplay", []string{file})
	if err == nil {
		a.logInfof("播放文件 %v \n %v", file, msg)
	}
	return err
}

func (a *App) getCli() *Cli {
	val := a.ctx.Value(CliKey)
	if val == nil {
		return nil
	}
	return val.(*Cli)
}

func (a *App) eventsEmit(eventName string, optionalData ...interface{}) {
	cli := a.getCli()
	if cli != nil {
		return
	}
	runtime.EventsEmit(a.ctx, eventName, optionalData...)
}

func (a *App) eventsOnce(eventName string, callback func(optionalData ...interface{})) {
	cli := a.getCli()
	if cli != nil {
		return
	}
	runtime.EventsOnce(a.ctx, eventName, callback)
}

func (a *App) messageDialog(dialogOptions runtime.MessageDialogOptions) (string, error) {
	cli := a.ctx.Value(CliKey).(*Cli)
	if cli != nil {
		return "", nil
	}
	return runtime.MessageDialog(a.ctx, dialogOptions)
}

func (a *App) eventsOn(eventName string, callback func(optionalData ...interface{})) {
	cli := a.getCli()
	if cli != nil {
		return
	}
	runtime.EventsOn(a.ctx, eventName, callback)
}

func (a *App) logInfof(format string, args ...interface{}) {
	cli := a.getCli()
	if cli != nil {
		if *cli.verbose {
			fmt.Printf("INFO | "+format+"\n", args...)
		}
	} else {
		runtime.LogInfof(a.ctx, format, args...)
	}
}

func (a *App) logInfo(message string) {
	cli := a.getCli()
	if cli != nil {
		if *cli.verbose {
			fmt.Println("INFO | " + message)
		}
	} else {
		runtime.LogInfo(a.ctx, message)
	}
}

func (a *App) logError(message string) {
	cli := a.getCli()
	if cli != nil {
		fmt.Println("ERR | " + message)
	} else {
		runtime.LogError(a.ctx, message)
	}
}

func (a *App) logErrorf(format string, args ...interface{}) {
	cli := a.getCli()
	if cli != nil {
		fmt.Printf("ERR | "+format+"\n", args...)
	} else {
		runtime.LogErrorf(a.ctx, format, args...)
	}
}
