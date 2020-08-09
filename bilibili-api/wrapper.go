package bilibili_api

import (
	"context"
	"fmt"
	"github.com/reactivex/rxgo/v2"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

// ass文件的头部
const templateHeader = `
[Script Info]
ScriptType: v4.00+
Collisions: Normal
PlayResX: {width}
PlayResY: {height}

[V4+ Styles]
Format: Name, Fontname, Fontsize, PrimaryColour, SecondaryColour, OutlineColour, BackColour, Bold, Italic, Underline, StrikeOut, ScaleX, ScaleY, Spacing, Angle, BorderStyle, Outline, Shadow, Alignment, MarginL, MarginR, MarginV, Encoding
Style: Default,{fontname},54,&H00FFFFFF,&H00FFFFFF,&H00000000,&H00000000,0,0,0,0,100,100,0.00,0.00,1,2.00,0.00,2,30,30,120,0
Style: Alternate,{fontname},36,&H00FFFFFF,&H00FFFFFF,&H00000000,&H00000000,0,0,0,0,100,100,0.00,0.00,1,2.00,0.00,2,30,30,84,0
Style: Danmaku,{fontname},{fontsize},&H00FFFFFF,&H00FFFFFF,&H00000000,&H00000000,0,0,0,0,100,100,0.00,0.00,1,1.00,0.00,2,30,30,30,0

[Events]
Format: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text
`

func NewDefaultAssConfig() ASSConfig {
	return ASSConfig{
		FontSize:       32,
		TuneDuration:   0,
		LineCount:      5,
		TemplateHeader: templateHeader,
		DropOffset:     2,
	}
}

func (request VideoStreamRequest) DownloadWithLatestDanmuku(to string, config ASSConfig, client *http.Client, opts ...rxgo.Option) rxgo.OptionalSingle {
	req := VideoInfoRequest{}
	req.Aid = request.Avid
	req.Bvid = request.Bvid
	cid := ""
	item, err := req.Fetch(client, opts...).Map(func(ctx context.Context, i interface{}) (interface{}, error) {
		dimension := i.(VideoInfo).Data.Dimension
		if config.ScreenWidth == 0 || config.ScreenHeight == 0 { // 如果宽高没指定，先获取视频信息，然后按照视频的宽高设置
			if dimension.IsPortrait() {
				config.ScreenWidth = dimension.Height
				config.ScreenHeight = dimension.Width
			} else {
				config.ScreenWidth = dimension.Width
				config.ScreenHeight = dimension.Height
			}
		}
		cid = strconv.Itoa(i.(VideoInfo).Data.Cid)
		return config, nil
	}).First().Get()
	if err != nil {
		return rxgo.Thrown(err).First()
	}
	if item.E != nil {
		return rxgo.Thrown(item.E).First()
	}
	config = item.V.(ASSConfig)
	dir := filepath.Dir(to)
	danmukuTmpPath := filepath.Join(dir, "danmuku.ass")
	t := time.Now()
	danReq := HistoryDanmukuIndexRequest{}
	danReq.Month = t.Format("2006-01")
	danReq.Oid = cid
	danmuku := danReq.Fetch(client, opts...).FlatMap(func(item rxgo.Item) rxgo.Observable {
		if item.E != nil {
			return rxgo.Thrown(item.E)
		}
		return rxgo.Just(item.V.(HistoryDanmukuIndex).Data)()
	}).Max(func(i interface{}, i2 interface{}) int {
		lhs, _ := time.Parse("2020-06-01", i.(string))
		rhs, _ := time.Parse("2006-06-01", i2.(string))
		return int(lhs.Sub(rhs).Seconds())
	}).Map(func(ctx context.Context, i interface{}) (interface{}, error) {
		req := HistoryDanmukuRequest{}
		req.Oid = cid
		req.Date = i.(string)
		item, err = req.DownloadAss(danmukuTmpPath, config, client).Get(rxgo.WithContext(ctx))
		if err != nil {
			return nil, err
		}
		return item.V, item.E
	})

	danmukuOB := rxgo.Create([]rxgo.Producer{func(ctx context.Context, next chan<- rxgo.Item) {
		item, err := danmuku.Get(rxgo.WithContext(ctx))
		if err != nil {
			next <- rxgo.Error(err)
		} else {
			next <- item
		}
	}})

	videoPath := filepath.Join(dir, "video.flv")
	videoOb := rxgo.Create([]rxgo.Producer{func(ctx context.Context, next chan<- rxgo.Item) {
		item, err := request.Download(videoPath, client, opts...).Get(rxgo.WithContext(ctx))
		if err != nil {
			next <- rxgo.Error(err)
		} else {
			next <- item
		}
	}})

	return rxgo.Concat([]rxgo.Observable{danmukuOB, videoOb}).Reduce(func(ctx context.Context, i interface{}, i2 interface{}) (interface{}, error) {
		return nil, nil
	}).Map(func(ctx context.Context, i interface{}) (interface{}, error) {
		cmd := exec.CommandContext(ctx, "ffmpeg", "-i", videoPath, "-vf", fmt.Sprintf("ass=%v", danmukuTmpPath), to)
		err := cmd.Run()
		_ = os.Remove(danmukuTmpPath)
		_ = os.Remove(videoPath)
		return to, err
	})
}
