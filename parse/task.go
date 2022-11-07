package parse

import (
	"GiveMeAnOffer/eventbus"
	"GiveMeAnOffer/logger"
	"bufio"
	"context"
	"crypto/tls"
	"crypto/x509"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/grafov/m3u8"
	"github.com/jamesnetherton/m3u"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type ParserTask struct {
	Url           string            `gorm:"primaryKey" json:"url"`
	TaskName      string            `json:"taskName"`
	Prefix        string            `json:"prefix"`
	DelOnComplete bool              `json:"delOnComplete"`
	KeyIV         string            `json:"keyIV"`
	Headers       string            `json:"headersMap"`
	HeadersMap    map[string]string `gorm:"-"`

	Ctx     context.Context         `gorm:"-" json:"-"`
	Client  *http.Client            `gorm:"-" json:"-"`
	Handler eventbus.RuntimeHandler `gorm:"-" json:"-"`
	Logger  logger.AppLogger        `gorm:"-" json:"-"`
	DstPath string                  `gorm:"-" json:"-"`
}

type PlayListInfo struct {
	Uri  string `json:"uri"`
	Desc string `json:"desc"`
}

type EventMessage struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Info    []*PlayListInfo `json:"info"`
	Title   string          `json:"title"`
}

type TaskType int

const (
	TaskTypeM3U TaskType = iota + 1
	TaskTypeM3U8
	TaskTypeChinaAACC
	TaskTypeBilibili
	TaskTypeCommon
)

type Result struct {
	Type TaskType
	Data interface{}
}

func (t *ParserTask) SetupRuntime(handler eventbus.RuntimeHandler) {
	t.Handler = handler
}

func (t *ParserTask) SetupContext(c context.Context) {
	t.Ctx = c
}

func (t *ParserTask) SetupClient(c *http.Client) {
	t.Client = c
}

func (t *ParserTask) SetupLogger(l logger.AppLogger) {
	t.Logger = l
}

func (t *ParserTask) SetupDownloadPath(p string) {
	t.DstPath = p
}

func (t *ParserTask) Parse() (*Result, error) {
	if t.HeadersMap == nil {
		t.HeadersMap = make(map[string]string)
	}

	if len(t.Headers) > 0 {
		headers := strings.Split(t.Headers, "\n")
		for _, header := range headers {
			arr := strings.Split(header, ":")
			if len(arr) != 2 {
				continue
			}
			t.HeadersMap[strings.TrimSpace(arr[0])] = strings.TrimSpace(arr[1])
		}
	}

	u, err := url.Parse(t.Url)
	if err != nil {
		return nil, err
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
		return &Result{Type: TaskTypeCommon, Data: []string{t.Url}}, nil
	}
}

func (t *ParserTask) handleM3UList() (*Result, error) {
	playlist, err := m3u.Parse(t.Url)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
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

	return &Result{Type: TaskTypeM3U, Data: tasks}, nil
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

func (t *ParserTask) parseLocalFile(path string) (*Result, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	p, listType, err := m3u8.DecodeFrom(bufio.NewReader(f), true)
	if err != nil {
		panic(err)
	}
	switch listType {
	case m3u8.MEDIA:
		return &Result{Type: TaskTypeM3U8, Data: p}, nil
	case m3u8.MASTER:
		return t.selectVariant(p.(*m3u8.MasterPlaylist))
	}
	return nil, nil
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

	origin := t.HeadersMap["Origin"]
	if len(origin) == 0 {
		origin = t.HeadersMap["origin"]
	}
	if len(origin) == 0 {
		origin = playlistUrl.Host
	}
	req.Header.Set("Origin", origin)

	refer := t.HeadersMap["Referer"]
	if len(refer) == 0 {
		refer = t.HeadersMap["referer"]
	}
	if len(refer) == 0 {
		refer = playlistUrl.Host
	}
	req.Header.Set("Referer", refer)

	req = req.WithContext(t.Ctx)
	return req, nil
}

func (t *ParserTask) selectVariant(l *m3u8.MasterPlaylist) (*Result, error) {
	// 等待前端选择
	msg := &EventMessage{
		Code:    1,
		Message: "请选择一种画质",
		Title:   "* 画质",
	}

	for i, variant := range l.Variants {
		msg.Info = append(msg.Info, &PlayListInfo{
			Desc: variant.Resolution,
			Uri:  strconv.Itoa(i),
		})
	}

	ch := make(chan int)
	t.Handler.EventsEmit(eventbus.SelectVariant, msg)
	t.Handler.EventsOnce(eventbus.OnVariantSelected, func(optionalData ...interface{}) {
		res := optionalData[0].(string)
		i, _ := strconv.Atoi(res)
		ch <- i
	})

	idx := <-ch
	return t.handleVariant(l.Variants[idx])
}

func (t *ParserTask) handleVariant(v *m3u8.Variant) (*Result, error) {
	if v.Chunklist != nil {
		return &Result{Type: TaskTypeM3U8, Data: v.Chunklist}, nil
	}
	req, err := t.BuildReq(v.URI)
	if err != nil {
		return nil, err
	}

	t.Url = req.URL.String()
	return t.Parse()
}

func (t *ParserTask) RetrieveM3U8List() (m3u8.Playlist, m3u8.ListType, error) {
	req, err := t.BuildReq(t.Url)
	if err != nil {
		return nil, 0, err
	}

	resp, err := t.Client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			t.Logger.LogError(err.Error())
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
	result, err := t.Handler.MessageDialog(runtime.MessageDialogOptions{
		Type:          runtime.QuestionDialog,
		Title:         "遇到证书错误",
		Message:       "是否忽略?",
		DefaultButton: "No",
		Buttons:       []string{"Yes", "No"},
	})
	if err != nil || result == "No" {
		return nil, 0, err
	}
	tr, ok := t.Client.Transport.(*http.Transport)
	if ok {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true} // 忽略证书错误
	}
	return t.RetrieveM3U8List()
}

func (t *ParserTask) getPlayerList() (*Result, error) {
	playlist, listType, err := t.RetrieveM3U8List()
	if err != nil {
		playlist, listType, err = t.handleFor509Error(err)
		if err != nil {
			return nil, err
		}
	}
	if listType == m3u8.MASTER {
		return t.selectVariant(playlist.(*m3u8.MasterPlaylist))
	} else {
		return &Result{Type: TaskTypeM3U8, Data: playlist}, nil
	}
}
