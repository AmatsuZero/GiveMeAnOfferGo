package main

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

type Cli struct {
	rootCmd       *cobra.Command
	downloadDir   string
	parserTask    *ParserTask
	delOnComplete *bool
	ctx           context.Context
}

func NewCli() *Cli {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "mode", "headless")
	cli := &Cli{ctx: ctx}
	SharedApp.headlessMode = true
	SharedApp.startup(ctx)

	rootCmd := &cobra.Command{
		Use:   "m3u8-download",
		Short: "m3u8 下载器",
		Long:  `m3u8 下载器，基于 wails 打造，支持 GUI 与命令行`,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}

	base, _ := os.UserHomeDir()
	rootCmd.PersistentFlags().StringVar(&SharedApp.config.PathDownloader, "downloadDir", filepath.Join(base, "Downloads"), "设置下载文件夹")
	rootCmd.PersistentFlags().BoolP("headless", "", true, "无 UI 启动")

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "输出 m3u8 下载器版本",
		Run:   cli.printVersion,
	}
	rootCmd.AddCommand(versionCmd)

	parseCmd := &cobra.Command{
		Use:   "parse",
		Short: "解析并下载 m3u8 文件",
		RunE:  cli.parse,
	}

	cli.parserTask = new(ParserTask)
	parseCmd.PersistentFlags().StringVar(&cli.parserTask.Url, "url", "", "设置 m3u8 地址, 多个地址用分号分割")
	parseCmd.MarkFlagRequired("url")

	cli.delOnComplete = parseCmd.Flags().BoolP("delOnComplete", "d", true, "合并完成后是否删除 ts 文件")
	parseCmd.PersistentFlags().StringVar(&cli.parserTask.KeyIV, "keyIV", "", "设置自定义密钥")
	parseCmd.PersistentFlags().StringVar(&cli.parserTask.Prefix, "prefix", "", "设置前缀")
	parseCmd.PersistentFlags().StringVar(&cli.parserTask.TaskName, "name", "", "输入文件名")

	rootCmd.AddCommand(parseCmd)

	cli.rootCmd = rootCmd

	return cli
}

func (c *Cli) parse(cmd *cobra.Command, args []string) error {
	adders := strings.Split(c.parserTask.Url, ",")
	if len(adders) == 0 {
		return fmt.Errorf("输入 m3u8 地址")
	}

	if len(adders) > 1 {
		c.parserTask.TaskName = ""
	}

	var tasks []*ParserTask
	for _, s := range adders {
		tasks = append(tasks, &ParserTask{
			Url:           s,
			TaskName:      c.parserTask.TaskName,
			Prefix:        c.parserTask.Prefix,
			DelOnComplete: *c.delOnComplete,
			KeyIV:         c.parserTask.KeyIV,
			Headers:       nil,
		})
	}

	return SharedApp.TaskAddMuti(tasks)
}

func (c *Cli) printVersion(cmd *cobra.Command, args []string) {
	fmt.Println("v1.0.0")
}

func (c *Cli) Execute() error {
	defer SharedApp.shutdown(c.ctx)
	return c.rootCmd.Execute()
}
