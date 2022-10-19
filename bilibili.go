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

func (v *videoInfoResp) download(t *ParserTask) error {
	t.TaskName = v.Data.Title

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

	wg := &sync.WaitGroup{}
	wg.Add(int(v.Data.Videos))

	for _, page := range v.Data.Pages {
		tmp, e := u.Parse(u.String())
		if e != nil {
			return err
		}
		vars := tmp.Query()
		vars.Add("qn", res)
		tmp.RawQuery = vars.Encode()
		go func(p *videoPageData) {
			err = p.download(u, wg, t)
			if err != nil {
				SharedApp.LogInfof("B站任务下载失败：%v", err)
			}
		}(page)
	}

	wg.Wait()
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

	err = merger.Merge()
	if err != nil {
		return err
	}

	if t.DelOnComplete {
		err = os.RemoveAll(downloader.DownloadDir)
		SharedApp.LogInfo("临时文件删除完成")
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
		Message: "请选择一种画质",
	}

	for _, format := range info.Data.SupportFormats {
		msg.Info = append(msg.Info, &playListInfo{
			Desc: format.NewDescription,
			Uri:  strconv.Itoa(format.Quality),
		})
	}

	ch := make(chan string)
	SharedApp.EventsEmit(SelectVariant, msg)
	SharedApp.EventsOnce(OnVariantSelected, func(optionalData ...interface{}) {
		ch <- optionalData[0].(string)
	})

	res := <-ch

	return res, nil
}

func (d *videoPageData) download(u *url.URL, g *sync.WaitGroup, t *ParserTask) error {
	defer g.Done()

	values := u.Query()
	values.Add("cid", strconv.Itoa(int(d.Cid)))
	u.RawQuery = values.Encode()

	// 获取指定清晰度的 url 列表
	resp, err := SharedApp.client.Get(u.String())
	if err != nil {
		return err
	}

	var info playUrlResp
	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		return err
	}

	t.Headers["Referer"] = "https://www.bilibili.com"
	t.Headers["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.114 Safari/537.36"
	return info.download(t, d.Part)
}

type bilibiliTaskType string

const (
	video bilibiliTaskType = "video"
)

type BilibiliParserTask struct {
	*ParserTask
	vid      string
	taskType bilibiliTaskType
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

	err = info.download(t.ParserTask)

	if err != nil {
		return err
	}

	return nil
}
