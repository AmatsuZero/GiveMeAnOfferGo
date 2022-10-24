package main

import (
	"encoding/json"
	"errors"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
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

func (r *playUrlResp) download(t *ParserTask, title string) error {

	downloader := &CommonDownloader{}
	dict := map[string]int64{}
	var list []string

	for _, item := range r.Data.Durl {
		list = append(list, item.Url)
		dict[path.Base(item.Url)] = item.Order
	}

	err := downloader.StartDownload(t, list)
	if err != nil {
		return err
	}

	// 遍历下载文件夹，调整顺序
	files, err := os.ReadDir(downloader.DownloadDir)
	var fileList []string
	if err != nil {
		return err
	}

	for _, f := range files {
		fileList = append(fileList, filepath.Join(downloader.DownloadDir, f.Name()))
	}

	sort.Slice(fileList, func(i, j int) bool {
		lhs, rhs := path.Base(fileList[i]), path.Base(fileList[j])
		return dict[lhs] < dict[rhs]
	})

	merger := &MergeFilesConfig{
		Files:     fileList,
		TsName:    title,
		MergeType: MergeTypeSpeed,
	}

	_, err = merger.Merge()
	if err != nil {
		return err
	}

	if t.DelOnComplete {
		err = os.RemoveAll(downloader.DownloadDir)
		SharedApp.logInfo("临时文件删除完成")
	}

	return err
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
		Message: "请选择一种画质",
	}

	for _, format := range info.Data.SupportFormats {
		msg.Info = append(msg.Info, &playListInfo{
			Desc: format.NewDescription,
			Uri:  strconv.Itoa(format.Quality),
		})
	}

	ch := make(chan string)
	SharedApp.eventsEmit(SelectVariant, msg)
	SharedApp.eventsOnce(OnVariantSelected, func(optionalData ...interface{}) {
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

func (t *BilibiliParserTask) Parse() (*ParseResult, error) {
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
		return nil, err
	}

	info := &videoInfoResp{}
	err = json.NewDecoder(resp.Body).Decode(info)
	if err != nil {
		return nil, err
	}

	return info.parse(t)
}

func (v *videoInfoResp) parse(t *BilibiliParserTask) (*ParseResult, error) {
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

	resolution, err := v.Data.Pages[0].selectResolution(u)
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
				SharedApp.logErrorf("B站分P解析失败: %v", e)
				return
			}
			rs = append(rs, r)
		}(page, tmp)
	}
	wg.Wait()

	return &ParseResult{Type: TaskTypeBilibili, Data: rs}, nil
}

func (d *videoPageData) parse(u *url.URL, t *BilibiliParserTask) (*BilibiliParserTask, error) {
	values := u.Query()
	values.Add("cid", strconv.Itoa(int(d.Cid)))
	u.RawQuery = values.Encode()

	// 获取指定清晰度的 url 列表
	resp, err := SharedApp.client.Get(u.String())
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
	ret.headers["Referer"] = "https://www.bilibili.com"
	ret.headers["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.114 Safari/537.36"
	return ret, nil
}
