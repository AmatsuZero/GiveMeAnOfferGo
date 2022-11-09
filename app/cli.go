package app

import (
	"GiveMeAnOffer/downloader"
	"GiveMeAnOffer/downloader/aria"
	"GiveMeAnOffer/merge"
	"GiveMeAnOffer/parse"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"GiveMeAnOffer/eventbus"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

type Cli struct {
	rootCmd  *cobra.Command
	verbose  *bool
	ctx      context.Context
	eventBus *eventbus.AsyncEventBus
}

const CliKey CtxKey = "cli"

func NewCli() *Cli {
	ctx := context.Background()
	cli := &Cli{
		ctx:      ctx,
		eventBus: eventbus.NewAsyncEventBus(),
		verbose:  new(bool),
	}

	ctx = context.WithValue(ctx, CliKey, cli)
	SharedApp.Startup(ctx)

	rootCmd := &cobra.Command{
		Use:   "m3u8-download",
		Short: "m3u8 下载器",
		Long:  `m3u8 下载器，基于 wails 打造，支持 GUI 与命令行`,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}
	cli.rootCmd = rootCmd

	downloadDir := SharedApp.Config.PathDownloader
	if len(downloadDir) == 0 {
		base, _ := os.UserHomeDir()
		downloadDir = filepath.Join(base, "Downloads")
	}

	rootCmd.PersistentFlags().StringVar(&SharedApp.Config.PathDownloader, "downloadDir", downloadDir, "设置下载文件夹")
	rootCmd.PersistentFlags().Bool("headless", true, "无 UI 启动")
	cli.verbose = rootCmd.PersistentFlags().BoolP("verbose", "v", false, "是否打印日志信息")

	cli.addVersionCmd().addParseCmd().addMergeFileCmd().addAriaCmd()

	return cli
}

func (c *Cli) addVersionCmd() *Cli {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "输出 m3u8 下载器版本",
		Run:   c.printVersion,
	}
	c.rootCmd.AddCommand(versionCmd)

	return c
}

func (c *Cli) addParseCmd() *Cli {
	var delOnComplete *bool
	var concurrentCnt *int
	parserTask := new(parse.ParserTask)

	parseCmd := &cobra.Command{
		Use:   "parse",
		Short: "解析并下载 m3u8 文件，按 q 终止",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			SharedApp.DomReady(c.ctx)
		},
		Run: func(cmd *cobra.Command, args []string) {
			SharedApp.concurrentLock = make(chan struct{}, *concurrentCnt)
			parserTask.DelOnComplete = *delOnComplete
			e := c.parse(parserTask)
			code := 0
			if e != nil {
				SharedApp.LogErrorf("解析失败: %v", e)
				code = 1
			} else {
				_, _ = fmt.Fprintln(os.Stdout, "解析结束")
			}
			os.Exit(code)
		},
	}

	parseCmd.Flags().StringVarP(&parserTask.Url, "url", "u", "", "设置 m3u8 地址, 多个地址用分号分割")
	_ = parseCmd.MarkFlagRequired("url")

	delOnComplete = parseCmd.Flags().BoolP("delOnComplete", "d", true, "合并完成后是否删除 ts 文件")
	parseCmd.Flags().StringVarP(&parserTask.KeyIV, "keyIV", "k", "", "设置自定义密钥")
	parseCmd.Flags().StringVarP(&parserTask.Prefix, "prefix", "p", "", "设置前缀")
	parseCmd.Flags().StringVarP(&parserTask.TaskName, "name", "n", "", "输入文件名")
	concurrentCnt = parseCmd.Flags().IntP("concurrent", "c", 3, "并发任务下载数量")

	c.rootCmd.AddCommand(parseCmd)

	return c
}

func (c *Cli) addMergeFileCmd() *Cli {
	config := new(merge.FilesConfig)
	files, dir := "", ""

	mergeFileCmd := &cobra.Command{
		Use:   "merge",
		Short: "合并 ts 文件",
		Run: func(cmd *cobra.Command, args []string) {
			code := 0
			defer os.Exit(code)
			if len(files) > 0 {
				config.Files = strings.Split(files, ",")
			} else if len(dir) > 0 {
				tsFiles, err := os.ReadDir(dir)
				if err != nil {
					SharedApp.LogErrorf("读取文件夹失败: %v", err)
					code = 1
					return
				}
				fileList := make([]string, 0, len(tsFiles))
				// 文件名排序
				sort.Slice(tsFiles, func(i, j int) bool {
					return tsFiles[i].Name() < tsFiles[j].Name()
				})
				for _, f := range tsFiles {
					fileList = append(fileList, filepath.Join(dir, f.Name()))
				}
				config.Files = fileList
			} else {
				SharedApp.LogError("必须指定要合并的文件或文件夹")
				code = 1
				return
			}
			o, e := SharedApp.StartMergeTs(config)
			if e != nil {
				SharedApp.LogErrorf("合并失败：%v", e)
				code = 1
			} else {
				_, _ = fmt.Fprintf(os.Stdout, "合并结束: %v", o)
			}
		},
	}

	c.rootCmd.AddCommand(mergeFileCmd)
	mergeFileCmd.Flags().StringVarP(&config.TsName, "name", "n", "", "输入文件名")
	mergeFileCmd.Flags().StringVarP((*string)(&config.MergeType), "type", "t", "speed", "转换类型")
	mergeFileCmd.Flags().StringVarP(&files, "files", "f", "", "要合并的 ts 文件，用逗号隔开")
	mergeFileCmd.Flags().StringVarP(&dir, "dir", "d", "", "要合并的 ts 文件夹")
	mergeFileCmd.Flags().StringVarP(&config.Output, "output", "o", "", "输出路径")
	mergeFileCmd.MarkFlagsMutuallyExclusive("files", "dir")

	return c
}

func (c *Cli) parse(task *parse.ParserTask) (err error) {
	SharedApp.LogInfof("🐱 下载地址: %v", SharedApp.Config.PathDownloader)
	_ = c.eventBus.Subscribe(eventbus.TaskStatusUpdate, func(item *downloader.DownloadTaskUIItem) {
		SharedApp.LogInfof(item.Status)
	})

	_ = c.eventBus.Subscribe(eventbus.TaskNotifyCreate, func(item *downloader.DownloadTaskUIItem) {
		SharedApp.LogInfof(item.Status)
	})

	_ = c.eventBus.Subscribe(eventbus.TaskAddEvent, func(item *downloader.DownloadTaskUIItem) {
		SharedApp.LogInfof(item.Status)
	})

	_ = c.eventBus.Subscribe(eventbus.SelectVariant, func(msg *parse.EventMessage) {
		c.selectVariant(msg)
	})

	_ = c.eventBus.Subscribe(eventbus.TaskFinish, func(item *downloader.DownloadTaskUIItem) {
		SharedApp.LogInfof(item.Status)
	})

	adders := strings.Split(task.Url, ",")
	if len(adders) == 0 {
		return fmt.Errorf("输入 m3u8 地址")
	}

	if len(adders) > 1 {
		task.TaskName = ""
	}

	var tasks []*parse.ParserTask
	for _, s := range adders {
		tasks = append(tasks, &parse.ParserTask{
			Url:           s,
			TaskName:      task.TaskName,
			Prefix:        task.Prefix,
			DelOnComplete: task.DelOnComplete,
			KeyIV:         task.KeyIV,
		})
	}

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
	_ = c.eventBus.Subscribe(eventbus.TaskFinish, func(item *downloader.DownloadTaskUIItem) {
		if item.State == downloader.DownloadTaskError {
			SharedApp.LogError(item.Status)
		}
		if item.IsDone || item.State == downloader.DownloadTaskError {
			done <- true
		}
	})

	<-done
	return
}

func (c *Cli) printVersion(_ *cobra.Command, _ []string) {
	_, _ = fmt.Fprintln(os.Stdout, "v1.0.0")
}

func (c *Cli) addAriaCmd() *Cli {
	port, secret := new(int), ""
	localPort := new(int)

	cmd := &cobra.Command{
		Use:   "aria",
		Short: "启动 aria 服务",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			SharedApp.Config.AriaConfig.RPCListenPort = *port
			SharedApp.Config.AriaConfig.RPCSecret = secret
			SharedApp.DomReady(c.ctx)
		},
		Run: func(cmd *cobra.Command, args []string) {

			client := aria.Client{
				Config: SharedApp.Config.AriaConfig,
			}

			client.RunLocal(*localPort)
		},
	}
	port = cmd.Flags().IntP("port", "p", 6800, "启动 aria2 端口号")
	cmd.Flags().StringVarP(&secret, "secret", "s", "123456", "设置 TOKEN")
	localPort = cmd.Flags().Int("localport", 8080, "本地服务端口号")
	c.rootCmd.AddCommand(cmd)
	return c
}

func (c *Cli) Execute() error {
	defer SharedApp.Shutdown(c.ctx)
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

func (c *Cli) selectVariant(msg *parse.EventMessage) {
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
		SharedApp.LogError(err.Error())
		return
	}

	c.eventBus.Publish(eventbus.OnVariantSelected, msg.Info[i].Uri)
}

func (c *Cli) quitKeyListening(tasks []*parse.ParserTask) {
	ch := make(chan string)
	go func(ch chan string) {
		// disable input buffering
		_ = exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
		// do not display entered characters on the screen
		_ = exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
		b := make([]byte, 1)
		for {
			_, _ = os.Stdin.Read(b)
			ch <- string(b)
		}
	}(ch)

	for {
		stdin := <-ch
		if stdin == "q" {
			for _, task := range tasks {
				c.eventBus.Publish(eventbus.TaskStop, task.Url)
			}
			break
		}
		time.Sleep(time.Millisecond * 100)
	}
}
