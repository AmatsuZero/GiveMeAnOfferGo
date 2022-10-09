package main

import (
	"encoding/json"
	"errors"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"net/url"
	"strconv"
	"strings"
)

const (
	baseUri     = "http://api.bilibili.com"
	videoStream = "/x/player/playurl"
	videoInfo   = "/x/web-interface/view"
)

var baseURl *url.URL

func init() {
	baseURl, _ = url.Parse(baseUri)
}

type baseResp struct {
	Code    int
	Message string
	Ttl     int
}

type dimension struct {
	Width, Height int64
	Rotate        bool // 是否宽高互换
}

func (d *dimension) UnmarshalJSON(data []byte) error {
	tmp := struct {
		Rotate int
	}{}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	d.Rotate = tmp.Rotate == 1

	return err
}

type videoInfoResp struct {
	baseResp
	Data struct {
		Bvid      string
		aid       int64
		Videos    int64  // 分区数量
		Tid       int64  // 分区tid
		Pic       string // 稿件封面
		Title     string // 稿件标题
		Cid       int64  // 稿件1P 分辨率
		Dimension dimension
		Pages     []*videoPageData // 视频分P列表
	}
}

func (v *videoInfoResp) download() error {
	u := baseURl.JoinPath(videoStream)
	values := u.Query()
	if len(v.Data.Bvid) > 0 {
		values.Add("bvid", v.Data.Bvid)
	} else {
		values.Add("aid", strconv.Itoa(int(v.Data.aid)))
	}
	u.RawQuery = values.Encode()

	if len(v.Data.Pages) == 0 {
		return errors.New("no page data")
	}

	res, err := v.Data.Pages[0].selectResolution(u)
	if err != nil {
		return err
	}
	println(res)
	return err
}

type playUrlResp struct {
	baseResp
	Data struct {
		From, Result, Message string
		AcceptDescription     []string         `json:"accept_description"`
		AcceptQuality         []int            `json:"accept_quality"`
		SupportFormats        []*supportFormat `json:"support_formats"`
		Durl                  []struct {
			Order        int64 // 分段序号
			Length       int64 // 视频长度，单位为毫秒
			Size         int64 // 视频大小，单位为 byte
			AHead, VHead string
			Url          string   // 有效时间为 120min
			BackupUrl    []string `json:"backup_url"`
		}
	}
}

type supportFormat struct {
	Quality             int
	Format, Superscript string
	newDescription      string `json:"new_description"`
	displayDescription  string `json:"display_description"`
	Codecs              []string
}

type videoPageData struct {
	Cid, Page int64
	From      string
	Part      string // 分P标题
	Duration  int64  // 分P时长（单位为秒）
	Dimension dimension
}

func (d *videoPageData) selectResolution(u *url.URL) (string, error) {
	values := u.Query()
	values.Add("cid", strconv.Itoa(int(d.Cid)))
	u.RawQuery = values.Encode()

	// 先获取信息
	resp, err := SharedApp.client.Get(u.String())
	if err != nil {
		return "", err
	}

	var info playUrlResp
	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		return "", err
	}

	msg := EventMessage{
		Code:    1,
		Message: "",
	}

	for _, format := range info.Data.SupportFormats {
		i := &playListInfo{
			Desc: format.Format,
		}
		values.Set("qn", strconv.Itoa(format.Quality))
		u.RawQuery = values.Encode()
		i.Uri = u.String()
		msg.Info = append(msg.Info, i)
	}

	values.Del("qn")
	u.RawQuery = values.Encode()

	ch := make(chan int)
	runtime.EventsEmit(SharedApp.ctx, SelectVariant, msg)
	runtime.EventsOnce(SharedApp.ctx, OnVariantSelected, func(optionalData ...interface{}) {
		idx := optionalData[0].(int)
		ch <- info.Data.AcceptQuality[idx]
	})

	ans := <-ch

	return strconv.Itoa(ans), nil
}

func (d *videoPageData) download(u *url.URL) error {
	return nil
}

type bilibiliTaskType string

const (
	video bilibiliTaskType = "video"
)

type BilibiliParserTask struct {
	vid      string
	taskType bilibiliTaskType
}

func NewBilibiliTask(u *url.URL) *BilibiliParserTask {
	t := &BilibiliParserTask{}
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
		t.vid = segments[1]
		t.taskType = bilibiliTaskType(segments[0])
	}
	return t
}

func (t *BilibiliParserTask) Parse() error {
	// 获取视频信息
	infoURL := baseURl.JoinPath(videoInfo)
	values := infoURL.Query()
	if strings.HasPrefix(t.vid, "BV") {
		values.Add("bvid", t.vid)
	} else {
		values.Add("avid", t.vid)
	}
	infoURL.RawQuery = values.Encode()

	resp, err := SharedApp.client.Get(infoURL.String())
	if err != nil {
		return err
	}

	var info videoInfoResp
	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		return err
	}

	err = info.download()

	if err != nil {
		return err
	}

	return nil
}
