package main

import (
	"github.com/AmatsuZero/mycli/commands"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"time"
)

func main() {
	now := time.Now().Local()
	log.Printf("🚀 启动任务：%v", now.Format(time.UnixDate))
	host := commands.FindAvailableHost()
	log.Printf("☁️ 使用域名为：%v", host)

	app := &cli.App{
		Commands: []*cli.Command{
			commands.CreateNewListCommand(host),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
