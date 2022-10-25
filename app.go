package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/grafov/m3u8"
	"gorm.io/gorm"

	"github.com/skratchdot/open-golang/open"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var configFilePath string
var appFolder string

type AppCtxKey string

func init() {
	appFolder, _ = os.UserConfigDir()
	if len(appFolder) == 0 {
		appFolder = os.Getenv("APPDATA")
	}
	appFolder = filepath.Join(appFolder, "M3U8-Downloader-GO")
	if _, err := os.Stat(appFolder); errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(appFolder, os.ModePerm)
		if err != nil {
			return
		}
	}

	configFilePath = filepath.Join(appFolder, "config.json")
}

// App struct
type App struct {
	config         *UserConfig
	ctx            context.Context
	client         *http.Client
	stopTasks      context.CancelFunc
	sniffer        *Sniffer
	concurrentLock chan struct{}

	db     *gorm.DB
	tasks  []*DownloadTaskUIItem
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
	a.initDB()
}

func (a *App) shutdown(ctx context.Context) {
	err := a.config.Save()
	if err != nil {
		runtime.LogError(ctx, err.Error())
	}
	a.stopTasks()
	a.saveTracks()
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
		err = a.handleBilibiliTask(ret)
	}
	return err
}

func (a *App) handleBilibiliTask(result *ParseResult) error {
	tasks := result.Data.([]*BilibiliParserTask)
	for _, task := range tasks {
		go func(t *BilibiliParserTask) {
			item := &DownloadTaskUIItem{
				ParserTask: t.ParserTask,
				Status:     "初始化...",
				State:      DownloadTaskProcessing,
			}
			a.eventsEmit(TaskNotifyCreate, item)

			downloader := &CommonDownloader{}
			downloader.NotifyItem = item
			err := downloader.StartDownload(t.ParserTask, t.Urls)

			if err != nil {
				item.Status = "下载失败，请检查链接有效性"
				item.State = DownloadTaskError
				a.eventsEmit(TaskFinish, item)
				return
			}

			// 遍历下载文件夹，调整顺序
			files, err := os.ReadDir(downloader.DownloadDir)
			var fileList []string
			if err != nil {
				item.Status = "读取文件夹失败"
				item.State = DownloadTaskError
				a.eventsEmit(TaskFinish, item)
				return
			}

			for _, f := range files {
				fileList = append(fileList, filepath.Join(downloader.DownloadDir, f.Name()))
			}

			sort.Slice(fileList, func(i, j int) bool {
				lhs, rhs := path.Base(fileList[i]), path.Base(fileList[j])
				return t.OrderDict[lhs] < t.OrderDict[rhs]
			})

			merger := &MergeFilesConfig{
				Files:     fileList,
				TsName:    t.TaskName,
				MergeType: MergeTypeSpeed,
			}

			output, err := merger.Merge()
			if err != nil {
				item.Status = "合并出错，请尝试手动合并"
				item.State = DownloadTaskError
				a.eventsEmit(TaskFinish, item)
				return
			}

			if t.DelOnComplete {
				err = os.RemoveAll(downloader.DownloadDir)
				if err != nil {
					SharedApp.logErrorf("临时文件删除失败：%v", err)
				} else {
					SharedApp.logInfo("临时文件删除完成")
				}
			}

			item.Status = "已完成"
			item.IsDone = true
			item.VideoPath = output
			item.State = DownloadTaskDone
			a.eventsEmit(TaskFinish, item)
		}(task)
	}
	return nil
}

func (a *App) handleM3UTask(result *ParseResult) (err error) {
	tasks := result.Data.([]*ParserTask)
	ch := make(chan int)

	msg := &EventMessage{
		Code:    1,
		Message: "请选择要下载的链接",
		Title:   "* 片源",
	}

	for i, task := range tasks {
		msg.Info = append(msg.Info, &playListInfo{
			Desc: task.TaskName,
			Uri:  strconv.Itoa(i),
		})
	}

	SharedApp.eventsEmit(SelectVariant, msg)
	SharedApp.eventsOnce(OnVariantSelected, func(optionalData ...interface{}) {
		res := optionalData[0].(string)
		i, _ := strconv.Atoi(res)
		ch <- i
	})

	idx := <-ch
	return a.TaskAdd(tasks[idx])
}

func (a *App) handleM3U8Task(task *ParserTask, result *ParseResult) (err error) {
	mpl := result.Data.(*m3u8.MediaPlaylist)
	queue := &M3U8DownloadQueue{}
	a.queues[task.Url] = queue
	// 任务列表添加任务
	item := a.addTaskNotifyItem(task)
	queue.NotifyItem = item

	info, cnt := "", mpl.Count()
	if mpl.Closed {
		d := time.Unix(int64(queue.TotalDuration), 0).Format("15:07:51")
		info = fmt.Sprintf("点播资源解析成功，有%v个片段, 时长：%v, 即将开始缓存...", cnt, d)
	} else {
		info = "直播资源解析成功，即将开始缓存..."
	}
	// 通知前端任务即将开始
	item.Status = info
	a.eventsEmit(TaskStatusUpdate, item)

	go func() {
		defer delete(a.queues, task.Url)

		item.State = DownloadTaskProcessing
		a.eventsEmit(TaskStatusUpdate, item)
		// 开始下载
		a.concurrentLock <- struct{}{}
		err = queue.StartDownload(task, mpl)
		// 下载完毕，释放资源
		<-a.concurrentLock
		if err != nil {
			item.Status = "下载失败，请检查链接有效性"
			item.State = DownloadTaskError
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
			item.State = DownloadTaskError
			a.eventsEmit(TaskFinish, item)
			return
		}

		a.logInfo("切片合并完成")
		if task.DelOnComplete {
			err = os.RemoveAll(queue.DownloadDir)
			if err != nil {
				a.logErrorf("切片删除失败: %v", err)
			} else {
				a.logInfo("切片删除完成")
			}
		}

		item.Status = "已完成"
		item.IsDone = true
		item.State = DownloadTaskDone
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
		item.State = DownloadTaskProcessing
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

func (a *App) isCliMode() bool {
	cli := a.getCli()
	return cli != nil
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
		cli.eventBus.Publish(eventName, optionalData...)
	} else {
		runtime.EventsEmit(a.ctx, eventName, optionalData...)
	}
}

func (a *App) eventsOnce(eventName string, callback func(optionalData ...interface{})) {
	cli := a.getCli()
	if cli != nil {
		_ = cli.eventBus.Subscribe(eventName, callback)
	} else {
		runtime.EventsOnce(a.ctx, eventName, callback)
	}
}

func (a *App) messageDialog(dialogOptions runtime.MessageDialogOptions) (string, error) {
	if a.isCliMode() {
		return "", nil
	}
	return runtime.MessageDialog(a.ctx, dialogOptions)
}

func (a *App) eventsOn(eventName string, callback func(optionalData ...interface{})) {
	cli := a.getCli()
	if cli != nil {
		_ = cli.eventBus.Subscribe(eventName, callback)
	} else {
		runtime.EventsOn(a.ctx, eventName, callback)
	}
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
	if a.isCliMode() {
		fmt.Println("ERR | " + message)
	} else {
		runtime.LogError(a.ctx, message)
	}
}

func (a *App) logErrorf(format string, args ...interface{}) {
	if a.isCliMode() {
		fmt.Printf("ERR | "+format+"\n", args...)
	} else {
		runtime.LogErrorf(a.ctx, format, args...)
	}
}
