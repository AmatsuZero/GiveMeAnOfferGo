package main

import (
	"bytes"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type MergeType string

const (
	MergeTypeSpeed       MergeType = "speed"
	MergeTypeTransCoding MergeType = "transcoding"
)

type MergeFilesConfig struct {
	Files     []string  `json:"files"`
	MergeType MergeType `json:"mergeType"` // speed: 快速合并 / transcoding：修复合并(慢|转码)
	TsName    string    `json:"taskName"`
}

func NewMergeConfigFromDownloadQueue(q *M3U8DownloadQueue, fileName string) *MergeFilesConfig {
	config := &MergeFilesConfig{
		MergeType: MergeTypeSpeed,
		TsName:    fileName,
	}

	sort.Slice(q.tasks, func(i, j int) bool {
		return q.tasks[i].Idx < q.tasks[j].Idx
	})

	for _, task := range q.tasks {
		if task.Done {
			config.Files = append(config.Files, task.Dst)
		}
	}

	return config
}

func Cmd(commandName string, params []string) (string, error) {
	cmd := exec.Command(commandName, params...)
	//fmt.Println("Cmd", cmd.Args)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return "", err
	}
	err = cmd.Wait()
	return out.String(), err
}

func (c *MergeFilesConfig) Merge() error {
	name := c.TsName
	if len(name) == 0 {
		name = "*"
	}
	name += ".txt"

	f, err := os.CreateTemp("", name)
	if err != nil {
		return err
	}

	for _, file := range c.Files {
		_, err = f.WriteString(fmt.Sprintf("file '%v'\n", file))
		if err != nil {
			return err
		}
	}

	err = f.Close()
	if err != nil {
		return err
	}

	audioCodec, videoCodec := "", ""
	switch c.MergeType {
	case MergeTypeSpeed:
		audioCodec = "copy"
		videoCodec = "copy"
	case MergeTypeTransCoding:
		audioCodec = "aac"
		videoCodec = "libx264"
	default:
		break
	}

	output := c.TsName
	if len(output) == 0 {
		output = fmt.Sprintf("%v", time.Now().Unix())
	}

	output = filepath.Join(SharedApp.config.PathDownloader, output+".mp4")
	cmdStr := fmt.Sprintf("ffmpeg -loglevel quiet -f concat -safe 0 -i %v -vcodec %v -acodec %v %v", f.Name(), videoCodec, audioCodec, output)
	args := strings.Split(cmdStr, " ")
	msg, err := Cmd(args[0], args[1:])
	if err != nil {
		runtime.LogError(SharedApp.ctx, fmt.Sprintf("videoConvert failed, %v, output: %v\n", err, msg))
	}
	defer func(name string) {
		err = os.Remove(name)
		if err != nil {
			runtime.LogInfof(SharedApp.ctx, "删除合并文件临时列表失败：%v", err.Error())
		}
	}(f.Name())
	return err
}
