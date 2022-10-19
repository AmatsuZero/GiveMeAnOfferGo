package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/skratchdot/open-golang/open"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var configFilePath string

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
	config       *UserConfig
	ctx          context.Context
	client       *http.Client
	stopTasks    context.CancelFunc
	headlessMode bool
	sniffer      *Sniffer
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
				"视频片段", "*.ts",
			},
			{
				"所有文件", "*",
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
	return task.Parse()
}

func (a *App) TaskAddMuti(tasks []*ParserTask) error {
	for _, task := range tasks {
		e := task.Parse()
		if e != nil {
			a.LogErrorf("下载任务失败:%v", e.Error())
		}
	}
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

func (a *App) EventsEmit(eventName string, optionalData ...interface{}) {
	if a.headlessMode {
		return
	}
	runtime.EventsEmit(a.ctx, eventName, optionalData...)
}

func (a *App) EventsOnce(eventName string, callback func(optionalData ...interface{})) {
	if a.headlessMode {
		return
	}
	runtime.EventsOnce(a.ctx, eventName, callback)
}

func (a *App) MessageDialog(dialogOptions runtime.MessageDialogOptions) (string, error) {
	if a.headlessMode {
		return "", nil
	}
	return runtime.MessageDialog(a.ctx, dialogOptions)
}

func (a *App) EventsOn(eventName string, callback func(optionalData ...interface{})) {
	if a.headlessMode {
		return
	}
	runtime.EventsOn(a.ctx, eventName, callback)
}

func (a *App) LogInfof(format string, args ...interface{}) {
	if a.headlessMode {
		fmt.Printf("INFO | "+format+"\n", args...)
	} else {
		runtime.LogInfof(a.ctx, format, args...)
	}
}

func (a *App) LogInfo(message string) {
	if a.headlessMode {
		fmt.Println("INFO | " + message)
	} else {
		runtime.LogInfo(a.ctx, message)
	}
}

func (a *App) LogError(message string) {
	if a.headlessMode {
		fmt.Println("ERR | " + message)
	} else {
		runtime.LogError(a.ctx, message)
	}
}

func (a *App) LogErrorf(format string, args ...interface{}) {
	if a.headlessMode {
		fmt.Printf("ERR | "+format+"\n", args...)
	} else {
		runtime.LogErrorf(a.ctx, format, args...)
	}
}
