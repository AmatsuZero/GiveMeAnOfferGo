package parse

import (
	"GiveMeAnOffer/eventbus"
	"encoding/json"
	"errors"
	"net/url"
	"path"
	"strconv"
	"strings"
	"sync"
)

const (
	baseUri     = "http://api.bilibili.com"
	videoStream = "/x/player/playurl"
	videoInfo   = "/x/web-interface/view"
	pageList    = "/x/player/pagelist"
)

var baseURl *url.URL

func init() {
	baseURl, _ = url.Parse(baseUri)
}

type baseResp struct {
	Code, Ttl int
	Message   string
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
		Aid       int64
		Videos    int64  // 分区数量
		Tid       int64  // 分区tid
		Pic       string // 稿件封面
		Title     string // 稿件标题
		Cid       int64  // 稿件1P cid
		Dimension *dimension
		Pages     []*videoPageData // 视频分P列表
	}
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
	NewDescription      string `json:"new_description"`
	DisplayDescription  string `json:"display_description"`
	Codecs              []string
}

type videoPageData struct {
	Cid, Page  int64
	From       string
	Part       string // 分P标题
	Duration   int64  // 分P时长（单位为秒）
	Dimension  dimension
	FirstFrame string `json:"first_frame"`
}

func (d *videoPageData) selectResolution(u *url.URL, t *BilibiliParserTask) (string, error) {
	values := u.Query()
	values.Add("cid", strconv.Itoa(int(d.Cid)))
	u.RawQuery = values.Encode()

	// 先获取信息
	resp, err := t.Client.Get(u.String())
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
		Message: "请选择一种画质",
	}

	for _, format := range info.Data.SupportFormats {
		msg.Info = append(msg.Info, &PlayListInfo{
			Desc: format.NewDescription,
			Uri:  strconv.Itoa(format.Quality),
		})
	}

	ch := make(chan string)
	t.Handler.EventsEmit(eventbus.SelectVariant, msg)
	t.Handler.EventsOnce(eventbus.OnVariantSelected, func(optionalData ...interface{}) {
		ch <- optionalData[0].(string)
	})

	res := <-ch

	return res, nil
}

type bilibiliTaskType string

const (
	video bilibiliTaskType = "video"
)

type BilibiliParserTask struct {
	*ParserTask
	vid       string
	taskType  bilibiliTaskType
	Urls      []string
	OrderDict map[string]int64
}

func (t *BilibiliParserTask) Parse() (*Result, error) {
	// 获取视频信息
	infoURL := baseURl.JoinPath(videoInfo)
	values := infoURL.Query()
	if strings.HasPrefix(t.vid, "BV") {
		values.Add("bvid", t.vid)
	} else {
		values.Add("avid", t.vid)
	}
	infoURL.RawQuery = values.Encode()

	resp, err := t.Client.Get(infoURL.String())
	if err != nil {
		return nil, err
	}

	info := &videoInfoResp{}
	err = json.NewDecoder(resp.Body).Decode(info)
	if err != nil {
		return nil, err
	}

	return info.parse(t)
}

func (v *videoInfoResp) parse(t *BilibiliParserTask) (*Result, error) {
	t.TaskName = v.Data.Title

	u := baseURl.JoinPath(videoStream)
	values := u.Query()
	if len(v.Data.Bvid) > 0 {
		values.Add("bvid", v.Data.Bvid)
	} else {
		values.Add("aid", strconv.Itoa(int(v.Data.Aid)))
	}
	values.Add("fnval", "0")
	u.RawQuery = values.Encode()

	if len(v.Data.Pages) == 0 {
		return nil, errors.New("no page data")
	}

	resolution, err := v.Data.Pages[0].selectResolution(u, t)
	if err != nil {
		return nil, err
	}

	var rs []*BilibiliParserTask
	wg := &sync.WaitGroup{}

	for _, page := range v.Data.Pages {
		tmp, _ := u.Parse(u.String()) // 拷贝一份原来的 URL
		vars := tmp.Query()
		vars.Add("qn", resolution)
		tmp.RawQuery = vars.Encode()
		wg.Add(1)
		go func(p *videoPageData, u *url.URL) {
			defer wg.Done()
			r, e := p.parse(u, t)
			if e != nil {
				t.Logger.LogErrorf("B站分P解析失败: %v", e)
				return
			}
			rs = append(rs, r)
		}(page, tmp)
	}
	wg.Wait()

	return &Result{Type: TaskTypeBilibili, Data: rs}, nil
}

func (d *videoPageData) parse(u *url.URL, t *BilibiliParserTask) (*BilibiliParserTask, error) {
	values := u.Query()
	values.Add("cid", strconv.Itoa(int(d.Cid)))
	u.RawQuery = values.Encode()

	// 获取指定清晰度的 url 列表
	resp, err := t.Client.Get(u.String())
	if err != nil {
		return nil, err
	}
	info := new(playUrlResp)
	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		return nil, err
	}

	t, err = info.parse(t, d.Part)
	if err != nil {
		return nil, err
	}
	t.Url = u.String()
	return t, nil
}

func (r *playUrlResp) parse(t *BilibiliParserTask, title string) (*BilibiliParserTask, error) {
	var list []string
	dict := map[string]int64{}
	for _, item := range r.Data.Durl {
		list = append(list, item.Url)
		dict[path.Base(item.Url)] = item.Order
	}
	ret := &BilibiliParserTask{
		ParserTask: t.ParserTask,
		vid:        t.vid,
		taskType:   t.taskType,
		Urls:       list,
		OrderDict:  dict,
	}
	ret.TaskName = title
	ret.HeadersMap["Referer"] = "https://www.bilibili.com"
	ret.HeadersMap["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.114 Safari/537.36"
	return ret, nil
}
