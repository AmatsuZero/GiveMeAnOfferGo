package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
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
	a.config = config
	a.concurrentLock = make(chan struct{}, config.ConCurrentCnt)
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
	return config.Merge()
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
	err := task.Parse()
	return err
}

func (a *App) TaskAddMuti(tasks []*ParserTask) error {
	var wg sync.WaitGroup
	for _, task := range tasks {
		a.concurrentLock <- struct{}{}
		wg.Add(1)

		go func(t *ParserTask) {
			defer wg.Done()
			e := t.Parse()
			if e != nil {
				a.LogErrorf("下载任务失败:%v, 原因：%v", t.Url, e.Error())
			}
			<-a.concurrentLock
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
		a.LogInfof("播放文件 %v \n %v", file, msg)
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

func (a *App) EventsEmit(eventName string, optionalData ...interface{}) {
	cli := a.getCli()
	if cli != nil {
		return
	}
	runtime.EventsEmit(a.ctx, eventName, optionalData...)
}

func (a *App) EventsOnce(eventName string, callback func(optionalData ...interface{})) {
	cli := a.getCli()
	if cli != nil {
		return
	}
	runtime.EventsOnce(a.ctx, eventName, callback)
}

func (a *App) MessageDialog(dialogOptions runtime.MessageDialogOptions) (string, error) {
	cli := a.ctx.Value(CliKey).(*Cli)
	if cli != nil {
		return "", nil
	}
	return runtime.MessageDialog(a.ctx, dialogOptions)
}

func (a *App) EventsOn(eventName string, callback func(optionalData ...interface{})) {
	cli := a.getCli()
	if cli != nil {
		return
	}
	runtime.EventsOn(a.ctx, eventName, callback)
}

func (a *App) LogInfof(format string, args ...interface{}) {
	cli := a.getCli()
	if cli != nil {
		if *cli.verbose {
			fmt.Printf("INFO | "+format+"\n", args...)
		}
	} else {
		runtime.LogInfof(a.ctx, format, args...)
	}
}

func (a *App) LogInfo(message string) {
	cli := a.getCli()
	if cli != nil {
		if *cli.verbose {
			fmt.Println("INFO | " + message)
		}
	} else {
		runtime.LogInfo(a.ctx, message)
	}
}

func (a *App) LogError(message string) {
	cli := a.getCli()
	if cli != nil {
		fmt.Println("ERR | " + message)
	} else {
		runtime.LogError(a.ctx, message)
	}
}

func (a *App) LogErrorf(format string, args ...interface{}) {
	cli := a.getCli()
	if cli != nil {
		fmt.Printf("ERR | "+format+"\n", args...)
	} else {
		runtime.LogErrorf(a.ctx, format, args...)
	}
}
