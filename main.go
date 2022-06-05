package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"sort"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "func",
				Usage: "解析方法",
			},
			&cli.StringFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "是否详细打印",
			},
		},

		Commands: []*cli.Command{
			{
				Name:        "unused",
				Description: "分析项目的Mach-O文件，检查Mach-O文件中无用的类和方法",
				Action:      Unused,
			},
			{
				Name:        "size",
				Description: "分析项目的LinkMap文件，得出每个类或者库所占用的空间大小（代码段+数据段）",
				Action:      ClassSize,
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
