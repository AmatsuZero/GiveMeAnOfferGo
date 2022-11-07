package merge

import (
	"GiveMeAnOffer/downloader"
	"GiveMeAnOffer/logger"
	"GiveMeAnOffer/utils"
	"fmt"
	"github.com/flytam/filenamify"
	"os"
	"runtime"
	"sort"
	"strings"
)

type Type string

const (
	Speed       Type = "speed"
	TransCoding Type = "transcoding"
)

type FilesConfig struct {
	Files     []string `json:"files"`
	MergeType Type     `json:"mergeType"` // speed: 快速合并 / transcoding：修复合并(慢|转码)
	TsName    string   `json:"taskName"`
	Output    string   `json:"output"`

	logger logger.AppLogger
}

func (c *FilesConfig) SetupLogger(l logger.AppLogger) {
	c.logger = l
}

func NewMergeConfigFromDownloadQueue(q *downloader.M3U8DownloadQueue, fileName string) *FilesConfig {
	config := &FilesConfig{
		MergeType: Speed,
		TsName:    fileName,
	}

	sort.Slice(q.Tasks, func(i, j int) bool {
		return q.Tasks[i].Idx < q.Tasks[j].Idx
	})

	for _, task := range q.Tasks {
		if task.Done {
			config.Files = append(config.Files, task.Dst)
		}
	}

	return config
}

func (c *FilesConfig) Merge() (string, error) {
	fileName := c.TsName // 处理非法文件名
	fileName, err := filenamify.Filenamify(fileName, filenamify.Options{})
	if err != nil {
		return "", err
	}
	fileName = strings.ReplaceAll(fileName, " ", "") // 移除空格

	name := fileName
	if len(name) == 0 {
		name = "*"
	}
	name += ".txt"

	f, err := os.CreateTemp("", name)
	if err != nil {
		return "", err
	}

	for _, file := range c.Files {
		_, err = f.WriteString(fmt.Sprintf("file '%v'\n", file))
		if err != nil {
			return "", err
		}
	}

	err = f.Close()
	if err != nil {
		return "", err
	}

	audioCodec, videoCodec := "", ""
	switch c.MergeType {
	case Speed:
		audioCodec = "copy"
		videoCodec = "copy"
	case TransCoding:
		audioCodec = "aac"
		videoCodec = "libx264"
	default:
		break
	}

	cmdStr := fmt.Sprintf("ffmpeg -loglevel 16 -f concat -safe 0 -i %v -vcodec %v -acodec %v", f.Name(), videoCodec, audioCodec)
	args := strings.Split(cmdStr, " ")
	if runtime.GOOS == "linux" { // FIX：linux 主机合并失败
		args = append(args, "-bsf:a", "aac_adtstoasc")
	}
	// 处理进度
	//	args = append(args, "-progress")
	args = append(args, c.Output)
	msg, err := utils.Cmd(args[0], args[1:])
	if err != nil {
		c.logger.LogErrorf("video merge failed, %v, output: %v\n", err, msg)
	}
	defer func(name string) {
		err = os.Remove(name)
		if err != nil {
			c.logger.LogErrorf("删除合并文件临时列表失败：%v", err.Error())
		}
	}(f.Name())
	return c.Output, err
}
