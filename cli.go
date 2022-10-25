package main

import (
	"context"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"GiveMeAnOffer/eventbus"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

type Cli struct {
	rootCmd       *cobra.Command
	parserTask    *ParserTask
	delOnComplete *bool
	verbose       *bool
	concurrentCnt *int
	ctx           context.Context
	eventBus      *eventbus.AsyncEventBus
}

const CliKey AppCtxKey = "cli"

func NewCli() *Cli {
	ctx := context.Background()
	cli := &Cli{
		ctx:      ctx,
		eventBus: eventbus.NewAsyncEventBus(),
	}

	ctx = context.WithValue(ctx, CliKey, cli)
	SharedApp.startup(ctx)

	rootCmd := &cobra.Command{
		Use:   "m3u8-download",
		Short: "m3u8 下载器",
		Long:  `m3u8 下载器，基于 wails 打造，支持 GUI 与命令行`,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}

	downloadDir := SharedApp.config.PathDownloader
	if len(downloadDir) == 0 {
		base, _ := os.UserHomeDir()
		downloadDir = filepath.Join(base, "Downloads")
	}

	rootCmd.PersistentFlags().StringVar(&SharedApp.config.PathDownloader, "downloadDir", downloadDir, "设置下载文件夹")
	rootCmd.PersistentFlags().Bool("headless", true, "无 UI 启动")
	cli.concurrentCnt = rootCmd.PersistentFlags().IntP("concurrent", "n", 3, "并发任务下载数量")
	cli.verbose = rootCmd.PersistentFlags().BoolP("verbose", "v", false, "是否打印日志信息")

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "输出 m3u8 下载器版本",
		Run:   cli.printVersion,
	}
	rootCmd.AddCommand(versionCmd)

	parseCmd := &cobra.Command{
		Use:   "parse",
		Short: "解析并下载 m3u8 文件，按 q 终止",
		RunE:  cli.parse,
	}

	cli.parserTask = new(ParserTask)
	parseCmd.PersistentFlags().StringVar(&cli.parserTask.Url, "url", "", "设置 m3u8 地址, 多个地址用分号分割")
	_ = parseCmd.MarkFlagRequired("url")

	cli.delOnComplete = parseCmd.Flags().BoolP("delOnComplete", "d", true, "合并完成后是否删除 ts 文件")
	parseCmd.PersistentFlags().StringVar(&cli.parserTask.KeyIV, "keyIV", "", "设置自定义密钥")
	parseCmd.PersistentFlags().StringVar(&cli.parserTask.Prefix, "prefix", "", "设置前缀")
	parseCmd.PersistentFlags().StringVar(&cli.parserTask.TaskName, "name", "", "输入文件名")

	rootCmd.AddCommand(parseCmd)

	cli.rootCmd = rootCmd

	return cli
}

func (c *Cli) parse(cmd *cobra.Command, args []string) (err error) {
	SharedApp.logInfof("🐱 下载地址: %v", SharedApp.config.PathDownloader)
	_ = c.eventBus.Subscribe(TaskStatusUpdate, func(item *DownloadTaskUIItem) {
		SharedApp.logInfof(item.Status)
	})

	_ = c.eventBus.Subscribe(TaskNotifyCreate, func(item *DownloadTaskUIItem) {
		SharedApp.logInfof(item.Status)
	})

	_ = c.eventBus.Subscribe(TaskAddEvent, func(item *DownloadTaskUIItem) {
		SharedApp.logInfof(item.Status)
	})

	_ = c.eventBus.Subscribe(SelectVariant, c.selectVariant)

	SharedApp.concurrentLock = make(chan struct{}, *c.concurrentCnt)
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
		})
	}
	defer fmt.Println("解析结束")
	if len(tasks) == 1 {
		err = SharedApp.TaskAdd(tasks[0])
	} else {
		err = SharedApp.TaskAddMuti(tasks)
	}

	if err != nil {
		return err
	}

	go c.quitKeyListening(tasks)
	done := make(chan bool)
	_ = c.eventBus.Subscribe(TaskFinish, func(item *DownloadTaskUIItem) {
		if item.State == DownloadTaskError {
			SharedApp.logError(item.Status)
		}
		if item.IsDone || item.State == DownloadTaskError {
			done <- true
		}
	})

	<-done
	return
}

func (c *Cli) printVersion(cmd *cobra.Command, args []string) {
	fmt.Println("v1.0.0")
}

func (c *Cli) Execute() error {
	defer SharedApp.shutdown(c.ctx)
	return c.rootCmd.Execute()
}

func (c *Cli) MessageDialog(ops runtime.MessageDialogOptions) (string, error) {
	prompt := promptui.Select{
		Label: ops.Title,
		Items: ops.Buttons,
	}

	_, result, err := prompt.Run()
	return result, err
}

func (c *Cli) selectVariant(msg EventMessage) {
	var labels []string

	for _, info := range msg.Info {
		labels = append(labels, info.Desc)
	}

	prompt := promptui.Select{
		Label: "请选择",
		Items: labels,
	}

	i, _, err := prompt.Run()
	if err != nil {
		SharedApp.logError(err.Error())
		return
	}

	c.eventBus.Publish(OnVariantSelected, msg.Info[i].Uri)
}

func (c *Cli) quitKeyListening(tasks []*ParserTask) {
	ch := make(chan string)
	go func(ch chan string) {
		// disable input buffering
		exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
		// do not display entered characters on the screen
		exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
		var b []byte = make([]byte, 1)
		for {
			os.Stdin.Read(b)
			ch <- string(b)
		}
	}(ch)

	for {
		stdin := <-ch
		if stdin == "q" {
			for _, task := range tasks {
				c.eventBus.Publish(TaskStop, task.Url)
			}
			break
		}
		time.Sleep(time.Millisecond * 100)
	}
}
