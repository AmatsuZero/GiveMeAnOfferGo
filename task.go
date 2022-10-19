package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/grafov/m3u8"
	"github.com/jamesnetherton/m3u"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

const (
	TaskAddEvent           = "task-add-reply"
	SelectVariant          = "select-variant"
	OnVariantSelected      = "variant-selected"
	TaskNotifyCreate       = "task-notify-create"
	StopLiveStreamDownload = "stop-live-stream-download"
)

type ParserTask struct {
	Url           string            `json:"url"`
	TaskName      string            `json:"taskName"`
	Prefix        string            `json:"prefix"`
	DelOnComplete bool              `json:"delOnComplete"`
	KeyIV         string            `json:"keyIV"`
	Headers       map[string]string `json:"headers"`
}

type playListInfo struct {
	Uri  string `json:"uri"`
	Desc string `json:"desc"`
}

type EventMessage struct {
	Code    int
	Message string
	Info    []*playListInfo
}

func (t *ParserTask) Parse() error {
	if t.Headers == nil {
		t.Headers = make(map[string]string)
	}

	u, err := url.Parse(t.Url)
	if err != nil {
		return err
	}

	if u.Host == "www.bilibili.com" {
		b := t.NewBilibiliTask(u)
		return b.Parse()
	} else if u.Host == "www.chinaacc.com" {
		c := t.NewChinaAACCTask()
		return c.Parse()
	}

	ext := path.Ext(u.Path)
	switch {
	case ext == ".m3u": // m3u 文件需要单独处理：https://baike.baidu.com/item/m3u%E6%96%87%E4%BB%B6/365977
		return t.handleM3UList()
	case strings.Contains(u.Path, "m3u8"):
		isLocal := u.Scheme == "http" || u.Scheme == "https"
		if !isLocal {
			return t.parseLocalFile(u.String())
		} else {
			return t.getPlayerList()
		}
	default:
		q := &CommonDownloader{}
		return q.StartDownload(t, []string{t.Url})
	}
}

func (t *ParserTask) handleM3UList() error {
	playlist, err := m3u.Parse(t.Url)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	tasks := make([]*ParserTask, 0, len(playlist.Tracks))
	for _, track := range playlist.Tracks {
		tasks = append(tasks, &ParserTask{
			Url:           track.URI,
			TaskName:      track.Name,
			Headers:       t.Headers,
			DelOnComplete: t.DelOnComplete,
			Prefix:        t.Prefix,
			KeyIV:         t.KeyIV,
		})
	}

	return t.handleMultiTasks(tasks)
}

func (t *ParserTask) handleMultiTasks(tasks []*ParserTask) error {
	if len(tasks) == 0 {
		return nil
	}

	for _, task := range tasks {
		err := task.Parse()
		if err != nil {
			SharedApp.LogInfof("❌下载任务失败：%v", t.Url)
		}
	}
	return nil
}

func (t *ParserTask) NewChinaAACCTask() *ChinaAACCParserTask {
	return &ChinaAACCParserTask{
		ParserTask: t,
	}
}

func (t *ParserTask) NewBilibiliTask(u *url.URL) *BilibiliParserTask {
	ret := &BilibiliParserTask{
		ParserTask: t,
	}

	segments := strings.Split(u.Path, "/")
	n := 0
	for _, val := range segments {
		if len(val) > 0 {
			segments[n] = val
			n++
		}
	}
	segments = segments[:n]

	if strings.Contains(u.Path, "/video/") {
		ret.vid = segments[1]
		ret.taskType = bilibiliTaskType(segments[0])
	}
	return ret
}

func (t *ParserTask) parseLocalFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	p, listType, err := m3u8.DecodeFrom(bufio.NewReader(f), true)
	if err != nil {
		panic(err)
	}
	switch listType {
	case m3u8.MEDIA:
		return t.handleMediaPlayList(p.(*m3u8.MediaPlaylist))
	case m3u8.MASTER:
		return t.selectVariant(p.(*m3u8.MasterPlaylist))
	}
	return nil
}

func (t *ParserTask) buildSegmentsURL(u string) (string, error) {
	if strings.HasPrefix(u, "http") {
		return url.QueryUnescape(u)
	} else {
		playlistUrl, err := url.Parse(t.Url)
		msUrl, err := playlistUrl.Parse(u)
		if err != nil {
			return "", err
		}
		return url.QueryUnescape(msUrl.String())
	}
}

func (t *ParserTask) BuildReq(u string) (*http.Request, error) {
	if u != t.Url {
		var err error
		u, err = t.buildSegmentsURL(u)
		if err != nil {
			return nil, err
		}
	}
	playlistUrl, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", playlistUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	origin := t.Headers["Origin"]
	if len(origin) == 0 {
		origin = t.Headers["origin"]
	}
	if len(origin) == 0 {
		origin = playlistUrl.Host
	}
	req.Header.Set("Origin", origin)

	refer := t.Headers["Referer"]
	if len(refer) == 0 {
		refer = t.Headers["referer"]
	}
	if len(refer) == 0 {
		refer = playlistUrl.Host
	}
	req.Header.Set("Referer", refer)

	req = req.WithContext(SharedApp.ctx)
	return req, nil
}

func (t *ParserTask) selectVariant(l *m3u8.MasterPlaylist) error {
	// 等待前端选择
	msg := EventMessage{
		Code:    1,
		Message: "",
	}

	for i, variant := range l.Variants {
		msg.Info = append(msg.Info, &playListInfo{
			Desc: variant.Resolution,
			Uri:  strconv.Itoa(i),
		})
	}

	ch := make(chan int)
	SharedApp.EventsEmit(SelectVariant, msg)
	SharedApp.EventsOnce(OnVariantSelected, func(optionalData ...interface{}) {
		res := optionalData[0].(string)
		i, _ := strconv.Atoi(res)
		ch <- i
	})

	idx := <-ch
	return t.handleVariant(l.Variants[idx])
}

func (t *ParserTask) handleVariant(v *m3u8.Variant) error {
	if v.Chunklist != nil {
		return t.handleMediaPlayList(v.Chunklist)
	}
	req, err := t.BuildReq(v.URI)
	if err != nil {
		return err
	}

	t.Url = req.URL.String()
	err = t.Parse()
	return err
}

func (t *ParserTask) handleMediaPlayList(mpl *m3u8.MediaPlaylist) error {
	cnt := 0
	info := ""

	queue := &M3U8DownloadQueue{}

	ch := make(chan bool)
	if mpl.Closed {
		d := time.Unix(int64(queue.TotalDuration), 0).Format("15:07:51")
		info = fmt.Sprintf("点播资源解析成功，有%v个片段，时长：%v，，即将开始缓存...", cnt, d)
	} else {
		info = "直播资源解析成功，即将开始缓存..."
	}

	var err error
	go func(c chan bool) {
		err = queue.StartDownload(t, mpl)
		c <- true
	}(ch)

	SharedApp.EventsEmit(TaskAddEvent, EventMessage{
		Code:    0,
		Message: info,
	})

	<-ch
	if err != nil {
		return err
	}

	SharedApp.LogInfof("切片下载完成，一共%v个", len(queue.tasks))

	merger := NewMergeConfigFromDownloadQueue(queue, t.TaskName)
	err = merger.Merge()
	if err != nil {
		return err
	}
	SharedApp.LogInfo("切片合并完成")

	if t.DelOnComplete {
		err = os.RemoveAll(queue.DownloadDir)
		SharedApp.LogInfo("切片删除完成")
	}
	return err
}

func (t *ParserTask) retrieveM3U8List() (m3u8.Playlist, m3u8.ListType, error) {
	req, err := t.BuildReq(t.Url)
	if err != nil {
		return nil, 0, err
	}

	resp, err := SharedApp.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			SharedApp.LogError(err.Error())
		}
	}(resp.Body)

	return m3u8.DecodeFrom(resp.Body, true)
}

func (t *ParserTask) handleFor509Error(err error) (m3u8.Playlist, m3u8.ListType, error) {
	e, ok := err.(*url.Error)
	if !ok {
		return nil, 0, err
	}
	if _, ok = e.Err.(x509.CertificateInvalidError); !ok {
		return nil, 0, err
	}
	result, err := SharedApp.MessageDialog(runtime.MessageDialogOptions{
		Type:          runtime.QuestionDialog,
		Title:         "遇到证书错误",
		Message:       "是否忽略?",
		DefaultButton: "No",
	})
	if err != nil || result == "No" {
		return nil, 0, err
	}
	tr := SharedApp.client.Transport.(*http.Transport)
	tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true} // 忽略证书错误
	return t.retrieveM3U8List()
}

func (t *ParserTask) getPlayerList() error {
	playlist, listType, err := t.retrieveM3U8List()
	if err != nil {
		playlist, listType, err = t.handleFor509Error(err)
		if err != nil {
			return err
		}
	}
	if listType == m3u8.MASTER {
		return t.selectVariant(playlist.(*m3u8.MasterPlaylist))
	} else {
		return t.handleMediaPlayList(playlist.(*m3u8.MediaPlaylist))
	}
}
