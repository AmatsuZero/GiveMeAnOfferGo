package main

import (
	"GiveMeAnOffer/app"
	"embed"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"os"
	"runtime"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
)

//go:embed frontend/dist
var assets embed.FS

func main() {

	for _, arg := range os.Args { // 检查是否以命令行模式启动
		if arg == "--headless" {
			cli := app.NewCli()
			err := cli.Execute()
			if err != nil {
				println("Error:", err.Error())
			}
			return
		}
	}

	// Create application with options
	err := wails.Run(&options.App{
		Title:              "m3u8-downloader",
		Width:              1024,
		Height:             768,
		Assets:             assets,
		BackgroundColour:   &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:          app.SharedApp.Startup,
		OnShutdown:         app.SharedApp.Shutdown,
		OnDomReady:         app.SharedApp.DomReady,
		LogLevelProduction: logger.ERROR,
		Frameless:          runtime.GOOS != "darwin",
		Bind: []interface{}{
			app.SharedApp,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
