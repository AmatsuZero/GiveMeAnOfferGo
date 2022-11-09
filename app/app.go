package app

import (
	"GiveMeAnOffer/downloader"
	"GiveMeAnOffer/eventbus"
	"GiveMeAnOffer/merge"
	"GiveMeAnOffer/parse"
	"GiveMeAnOffer/sniffer"
	"GiveMeAnOffer/utils"
	"context"
	"errors"
	"fmt"
	"github.com/flytam/filenamify"
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

type CtxKey string

var SharedApp *App

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

	configFilePath = filepath.Join(appFolder, "Config.json")

	// Create an instance of the app structure
	SharedApp = NewApp()
}

// App struct
type App struct {
	Config         *UserConfig
	ctx            context.Context
	client         *http.Client
	stopTasks      context.CancelFunc
	sniffer        *sniffer.Sniffer
	concurrentLock chan struct{}

	db     *gorm.DB
	tasks  []*downloader.DownloadTaskUIItem
	queues map[string]downloader.StoppableTask
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

func (a *App) Startup(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	a.ctx, a.stopTasks = ctx, cancel
	a.initDB()
	config, err := NewConfig(configFilePath)
	if err != nil {
		a.LogError(err.Error())
	} else if config.ConfigProxy != nil {
		// 写入代理配置
		err = os.Setenv("HTTP_PROXY", config.ConfigProxy.http)
		if err != nil {
			a.LogError(err.Error())
		}
		err = os.Setenv("HTTPS_PROXY", config.ConfigProxy.https)
		if err != nil {
			a.LogError(err.Error())
		}
	}

	a.Config = config
	a.concurrentLock = make(chan struct{}, config.ConCurrentCnt)
	a.tasks = make([]*downloader.DownloadTaskUIItem, 0)
}

func (a *App) Shutdown(ctx context.Context) {
	err := a.Config.Save()
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

		return files, err
	}

	dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
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

	if len(dir) == 0 {
		return nil, err
	}

	return a.OpenSelectTsDir(dir)
}

func (a *App) StartMergeTs(config *merge.FilesConfig) (string, error) {
	config.SetupLogger(a)
	fn := ""
	if len(config.Output) == 0 {
		fn = config.TsName
		if len(fn) > 0 {
			n, err := filenamify.Filenamify(fn, filenamify.Options{})
			if err != nil {
				n = fmt.Sprintf("%v", time.Now().Unix())
			}
			fn = n
		} else {
			fn = fmt.Sprintf("%v", time.Now().Unix())
		}
		config.Output = filepath.Join(a.Config.PathDownloader, fn+".mp4")
	}
	return config.Merge()
}

func (a *App) OpenConfigDir() (string, error) {
	defaultDir := a.Config.PathDownloader
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

	a.Config.PathDownloader = dir
	return dir, err
}

func (a *App) TaskAdd(task *parse.ParserTask) error {
	task.SetupRuntime(a)
	task.SetupClient(a.client)
	task.SetupContext(a.ctx)
	task.SetupLogger(a)
	task.SetupDownloadPath(a.Config.PathDownloader)
	ret, err := task.Parse()

	if err != nil {
		return err
	}
	switch ret.Type {
	case parse.TaskTypeCommon:
		err = a.handleCommonTask(task, ret)
	case parse.TaskTypeChinaAACC:
		err = a.handleAACCTask(ret)
	case parse.TaskTypeM3U8:
		err = a.handleM3U8Task(task, ret)
	case parse.TaskTypeM3U:
		err = a.handleM3UTask(ret)
	case parse.TaskTypeBilibili:
		err = a.handleBilibiliTask(ret)
	}
	return err
}

func (a *App) handleBilibiliTask(result *parse.Result) error {
	tasks := result.Data.([]*parse.BilibiliParserTask)
	for _, task := range tasks {
		go func(t *parse.BilibiliParserTask) {
			item := &downloader.DownloadTaskUIItem{
				ParserTask: t.ParserTask,
				Status:     "初始化...",
				State:      downloader.DownloadTaskProcessing,
			}
			a.EventsEmit(eventbus.TaskNotifyCreate, item)

			d := &downloader.CommonDownloader{}
			d.NotifyItem = item
			err := d.StartDownload(t.ParserTask, t.Urls)

			if err != nil {
				item.Status = "下载失败，请检查链接有效性"
				item.State = downloader.DownloadTaskError
				a.EventsEmit(eventbus.TaskFinish, item)
				return
			}

			// 遍历下载文件夹，调整顺序
			files, err := os.ReadDir(d.DownloadDir)
			var fileList []string
			if err != nil {
				item.Status = "读取文件夹失败"
				item.State = downloader.DownloadTaskError
				a.EventsEmit(eventbus.TaskFinish, item)
				return
			}

			for _, f := range files {
				fileList = append(fileList, filepath.Join(d.DownloadDir, f.Name()))
			}

			sort.Slice(fileList, func(i, j int) bool {
				lhs, rhs := path.Base(fileList[i]), path.Base(fileList[j])
				return t.OrderDict[lhs] < t.OrderDict[rhs]
			})

			merger := &merge.FilesConfig{
				Files:     fileList,
				TsName:    t.TaskName,
				MergeType: merge.Speed,
			}

			output, err := a.StartMergeTs(merger)
			if err != nil {
				item.Status = "合并出错，请尝试手动合并"
				item.State = downloader.DownloadTaskError
				a.EventsEmit(eventbus.TaskFinish, item)
				return
			}

			if t.DelOnComplete {
				err = os.RemoveAll(d.DownloadDir)
				if err != nil {
					a.LogErrorf("临时文件删除失败：%v", err)
				} else {
					a.LogInfo("临时文件删除完成")
				}
			}

			item.Status = "已完成"
			item.IsDone = true
			item.VideoPath = output
			item.State = downloader.DownloadTaskDone
			a.EventsEmit(eventbus.TaskFinish, item)
		}(task)
	}
	return nil
}

func (a *App) handleM3UTask(result *parse.Result) (err error) {
	tasks := result.Data.([]*parse.ParserTask)
	ch := make(chan int)

	msg := &parse.EventMessage{
		Code:    1,
		Message: "请选择要下载的链接",
		Title:   "* 片源",
	}

	for i, task := range tasks {
		msg.Info = append(msg.Info, &parse.PlayListInfo{
			Desc: task.TaskName,
			Uri:  strconv.Itoa(i),
		})
	}

	a.EventsEmit(eventbus.SelectVariant, msg)
	a.EventsOnce(eventbus.OnVariantSelected, func(optionalData ...interface{}) {
		res := optionalData[0].(string)
		i, _ := strconv.Atoi(res)
		ch <- i
	})

	idx := <-ch
	return a.TaskAdd(tasks[idx])
}

func (a *App) handleM3U8Task(task *parse.ParserTask, result *parse.Result) (err error) {
	mpl := result.Data.(*m3u8.MediaPlaylist)
	queue := &downloader.M3U8DownloadQueue{}
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
	a.EventsEmit(eventbus.TaskStatusUpdate, item)

	shouldStop := false
	a.EventsOn(eventbus.TaskStop, func(optionalData ...interface{}) {
		u := optionalData[0].(string)
		if u != task.Url { // 只停止自己的任务
			return
		}
		queue.Stop()
		shouldStop = true
		item.Status = "任务停止中..."
		a.EventsEmit(eventbus.TaskFinish, item)
	})

	go func() {
		defer delete(a.queues, task.Url)

		item.State = downloader.DownloadTaskProcessing
		a.EventsEmit(eventbus.TaskStatusUpdate, item)
		// 开始下载
		a.concurrentLock <- struct{}{}
		err = queue.StartDownload(task, mpl)
		// 下载完毕，释放资源
		<-a.concurrentLock
		if shouldStop {
			item.Status = "任务已经停止"
			item.State = downloader.DownloadTaskError
			a.EventsEmit(eventbus.TaskFinish, item)
			return
		} else if err != nil {
			item.Status = "下载失败，请检查链接有效性"
			item.State = downloader.DownloadTaskError
			a.EventsEmit(eventbus.TaskFinish, item)
			return
		}

		item.Status = "已完成，合并中..."
		a.EventsEmit(eventbus.TaskFinish, item)
		a.LogInfof("切片下载完成，一共%v个", len(queue.Tasks))

		merger := merge.NewMergeConfigFromDownloadQueue(queue, task.TaskName)
		output, err := a.StartMergeTs(merger)
		if err != nil {
			item.Status = "合并出错，请尝试手动合并"
			item.State = downloader.DownloadTaskError
			a.EventsEmit(eventbus.TaskFinish, item)
			return
		}

		a.LogInfo("切片合并完成")
		if task.DelOnComplete {
			err = os.RemoveAll(queue.DownloadDir)
			if err != nil {
				a.LogErrorf("切片删除失败: %v", err)
			} else {
				a.LogInfo("切片删除完成")
			}
		}

		item.Status = "已完成"
		item.IsDone = true
		item.State = downloader.DownloadTaskDone
		item.VideoPath = output
		a.EventsEmit(eventbus.TaskFinish, item)
	}()
	return
}

func (a *App) handleAACCTask(result *parse.Result) error {
	ch := result.Data.(chan *parse.ParserTask)
	for parserTask := range ch {
		err := a.TaskAdd(parserTask)
		if err != nil {
			a.LogErrorf("正保网校课程下载失败，任务名：%v，链接：%v", parserTask.TaskName, parserTask.Url)
		}
	}
	return nil
}

func (a *App) handleCommonTask(task *parse.ParserTask, result *parse.Result) (err error) {
	item := a.addTaskNotifyItem(task)

	go func() {
		item.State = downloader.DownloadTaskProcessing
		a.concurrentLock <- struct{}{}
		q := &downloader.CommonDownloader{}
		a.queues[task.Url] = q
		q.NotifyItem = item
		err = q.StartDownload(task, result.Data.([]string))
		<-a.concurrentLock

		item.Status = "已完成"
		item.IsDone = true
		item.State = downloader.DownloadTaskDone
		a.EventsEmit(eventbus.TaskFinish, item)
	}()

	return
}

func (a *App) TaskAddMuti(tasks []*parse.ParserTask) error {
	var wg sync.WaitGroup
	for _, task := range tasks {
		wg.Add(1)
		go func(t *parse.ParserTask) {
			defer wg.Done()
			e := a.TaskAdd(t)
			if e != nil {
				a.LogErrorf("下载任务失败:%v, 原因：%v", t.Url, e.Error())
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
	a.sniffer = sniffer.NewSniffer(u)
	a.sniffer.SetupLogger(a)
	a.sniffer.SetupLogger(a)
	return a.sniffer.GetLinks()
}

func (a *App) Open(link string) error {
	if len(link) == 0 {
		link = a.Config.PathDownloader
	}
	return open.Run(link)
}

func (a *App) Play(file string) error {
	msg, err := utils.Cmd("ffplay", []string{file})
	if err == nil {
		a.LogInfof("播放文件 %v \n %v", file, msg)
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

func (a *App) EventsEmit(eventName string, optionalData ...interface{}) {
	cli := a.getCli()
	if cli != nil {
		cli.eventBus.Publish(eventName, optionalData...)
	} else {
		runtime.EventsEmit(a.ctx, eventName, optionalData...)
	}
}

func (a *App) EventsOnce(eventName string, callback func(optionalData ...interface{})) {
	cli := a.getCli()
	if cli != nil {
		_ = cli.eventBus.Subscribe(eventName, callback)
	} else {
		runtime.EventsOnce(a.ctx, eventName, callback)
	}
}

func (a *App) MessageDialog(dialogOptions runtime.MessageDialogOptions) (string, error) {
	if a.isCliMode() {
		return "", nil
	}
	return runtime.MessageDialog(a.ctx, dialogOptions)
}

func (a *App) EventsOn(eventName string, callback func(optionalData ...interface{})) {
	cli := a.getCli()
	if cli != nil {
		_ = cli.eventBus.Subscribe(eventName, callback)
	} else {
		runtime.EventsOn(a.ctx, eventName, callback)
	}
}

func (a *App) LogInfof(format string, args ...interface{}) {
	cli := a.getCli()
	if cli != nil {
		if *cli.verbose {
			_, _ = fmt.Fprintf(os.Stdout, "INFO | "+format+"\n", args...)
		}
	} else {
		runtime.LogInfof(a.ctx, format, args...)
	}
}

func (a *App) LogInfo(message string) {
	cli := a.getCli()
	if cli != nil {
		if *cli.verbose {
			_, _ = fmt.Fprintln(os.Stdout, "INFO | "+message)
		}
	} else {
		runtime.LogInfo(a.ctx, message)
	}
}

func (a *App) LogError(message string) {
	if a.isCliMode() {
		_, _ = fmt.Fprintln(os.Stderr, "ERR | "+message)
	} else {
		runtime.LogError(a.ctx, message)
	}
}

func (a *App) LogErrorf(format string, args ...interface{}) {
	if a.isCliMode() {
		_, _ = fmt.Fprintf(os.Stderr, "ERR | "+format+"\n", args...)
	} else {
		runtime.LogErrorf(a.ctx, format, args...)
	}
}
