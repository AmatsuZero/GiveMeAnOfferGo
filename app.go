package main

import (
	"context"
	"errors"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"os"
	"path/filepath"
)

var configFilePath string

func init() {
	configFilePath, _ = os.UserConfigDir()
	if len(configFilePath) == 0 {
		configFilePath = os.Getenv("APPDATA")
	}
	configFilePath = filepath.Join(configFilePath, "M3U8-Downloader")

	if _, err := os.Stat(configFilePath); errors.Is(err, os.ErrNotExist) {
		os.Mkdir(configFilePath, os.ModePerm)
	}

	configFilePath = filepath.Join(configFilePath, "config.json")
}

// App struct
type App struct {
	config *UserConfig
	ctx    context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.config, _ = NewConfig(configFilePath)
}

func (a *App) shutdown(ctx context.Context) {
	a.config.Save()
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

type MergeFilesConfig struct {
	Files     []string `json:"files"`
	MergeType string   `json:"age"` // copy: 快速合并 / transcoding：修复合并(慢|转码)
	TsName    string   `json:"taskName"`
}

func (a *App) StartMergeTs(config MergeFilesConfig) {

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
