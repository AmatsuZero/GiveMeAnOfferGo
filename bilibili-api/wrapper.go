package bilibili_api

import (
	"context"
	"fmt"
	"github.com/reactivex/rxgo/v2"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func (request VideoStreamRequest) DownloadWithLatestDanmuku(to string, config ASSConfig, client *http.Client, opts ...rxgo.Option) rxgo.OptionalSingle {
	dir := filepath.Dir(to)
	danmukuTmpPath := filepath.Join(dir, "danmuku.ass")
	t := time.Now()
	danReq := HistoryDanmukuIndexRequest{}
	danReq.Month = t.Format("2006-01")
	danReq.Oid = request.Cid
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
		req.Oid = request.Cid
		req.Date = i.(string)
		item, err := req.DownloadAss(danmukuTmpPath, config, client).Get(rxgo.WithContext(ctx))
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
		err := exec.CommandContext(ctx, "ffmpeg", "-i", videoPath, "-vf", fmt.Sprintf("ass=%v", danmukuTmpPath), to).Run()
		_ = os.Remove(danmukuTmpPath)
		_ = os.Remove(videoPath)
		return to, err
	})
}
