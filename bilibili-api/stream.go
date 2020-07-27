package bilibili_api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/reactivex/rxgo/v2"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

type VideoStreamResolution int

const (
	VideoStreamResolution360             VideoStreamResolution = 16
	VideoStreamResolution480                                   = VideoStreamResolution360 * 2
	VideoStreamResolution720                                   = VideoStreamResolution480 * 2
	VideoStreamResolutionVIP720                                = VideoStreamResolution720 + 10
	VideoStreamResolution1080                                  = VideoStreamResolutionVIP720 + 6
	VideoStreamResolutionVIP1080Plus                           = VideoStreamResolution1080 + 32
	VideoStreamResolutionVIP1080and60FPS                       = VideoStreamResolutionVIP1080Plus + 4
	VideoStreamResolutionVIP4K                                 = VideoStreamResolutionVIP1080and60FPS + 4
)

func (resolution VideoStreamResolution) String() string {
	switch resolution {
	case VideoStreamResolutionVIP4K:
		return "4K 超清（大会员）"
	case VideoStreamResolutionVIP1080and60FPS:
		return "1080P60 高清（大会员）"
	case VideoStreamResolutionVIP1080Plus:
		return "1080P+ 高清（大会员）"
	case VideoStreamResolution1080:
		return "1080P 高清（登录）"
	case VideoStreamResolutionVIP720:
		return "720P60 高清（大会员）"
	case VideoStreamResolution720:
		return "720P 高清（登录）"
	case VideoStreamResolution480:
		return "480P 清晰"
	default:
		return "360P 流畅"
	}
}

func (resolution VideoStreamResolution) Encode() string {
	return fmt.Sprintf("%d", resolution)
}

type VideoStreamRequest struct {
	baseRequest
	Avid        string
	Bvid        string
	Cid         string
	Resolution  VideoStreamResolution
	ShouldUse4K bool
}

func (request VideoStreamRequest) IsParamsValid() bool {
	return (len(request.Avid) > 0 || len(request.Bvid) > 0) && len(request.Cid) > 0
}

func (request VideoStreamRequest) Request() (*http.Request, error) {
	if !request.IsParamsValid() {
		return nil, kInvalidParamError
	}
	base := *kBaseURL
	base.Path = "x/player/playurl"
	base.RawQuery = request.queryItems(base.Query()).Encode()
	return http.NewRequest("GET", base.String(), nil)
}

func (request VideoStreamRequest) Fetch(client *http.Client, opts ...rxgo.Option) rxgo.Observable {
	req, err := request.Request()
	if err != nil {
		return rxgo.Thrown(err)
	}
	return request.fetch(client, req, opts...).Map(func(ctx context.Context, i interface{}) (interface{}, error) {
		var info VideoStreamInfo
		data := i.([]byte)
		err := json.Unmarshal(data, &info)
		if err != nil {
			return nil, err
		}
		return info, nil
	})
}

func (request VideoStreamRequest) queryItems(query url.Values) url.Values {
	if len(request.Avid) > 0 {
		query.Add("avid", request.Avid)
	} else {
		query.Add("bvid", request.Bvid)
	}
	query.Set("cid", request.Cid)
	if request.Resolution > 0 {
		query.Set("qn", request.Resolution.Encode())
	}
	if request.ShouldUse4K {
		query.Set("fourk", "1")
	}
	return query
}

type VideoSegments struct {
	Order      int // 分段序号
	Length     int
	Size       int // 单位为Byte
	Ahead      string
	Vhead      string
	Url        string   // 视频流url
	BackupUrl  []string `json:"backup_url"`
	retryCount int      // 标志当前重试次数
	targetPath string   // 要写入的路径
}

func (segments VideoSegments) Request() (req *http.Request, err error) {
	if segments.retryCount > len(segments.BackupUrl) {
		return nil, fmt.Errorf("can not retry")
	}
	if segments.retryCount == 0 {
		req, err = http.NewRequest("GET", segments.Url, nil)
	} else {
		req, err = http.NewRequest("GET", segments.BackupUrl[segments.retryCount-1], nil)
	}
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", kFakeWebUA)
	req.Header.Set("Refer", kFakeRefer)
	return
}

func (segments VideoSegments) IsParamsValid() bool {
	if segments.retryCount == 0 {
		return len(segments.Url) > 0
	}
	return len(segments.BackupUrl[segments.retryCount-1]) > 0
}

func (segments *VideoSegments) Fetch(client *http.Client, opts ...rxgo.Option) rxgo.Observable {
	if client == nil {
		client = http.DefaultClient
	}
	return rxgo.Defer([]rxgo.Producer{func(ctx context.Context, next chan<- rxgo.Item) {
		req, err := segments.Request()
		if err != nil {
			next <- rxgo.Error(err)
			return
		}
		req = req.WithContext(ctx)
		resp, err := client.Do(req)
		next <- rxgo.Item{
			V: resp,
			E: err,
		}
	}}).Retry(len(segments.BackupUrl), func(err error) bool {
		segments.retryCount++
		return true
	})
}

type VideoStreamInfo struct {
	BaseResponse
	Data struct {
		From              string
		Result            string
		Message           string
		Quality           VideoStreamResolution
		Format            string
		TimeLength        int                     `json:"timelength"`
		AcceptFormat      string                  `json:"accept_format"`
		AcceptDescription []string                `json:"accept_description"`
		AcceptQuality     []VideoStreamResolution `json:"accept_quality"`
		VideoCodecid      int                     `json:"video_codecid"`
		SeekParam         string                  `json:"seek_param"`
		SeekType          string                  `json:"seek_type"`
		Durl              []VideoSegments
	}
}

func (info VideoStreamInfo) SupportFormats() []string {
	return strings.Split(info.Data.AcceptFormat, ",")
}

func (info VideoStreamInfo) Download(path string, client *http.Client, opts ...rxgo.Option) rxgo.Disposed {
	source := make([]VideoSegments, 0, len(info.Data.Durl))
	_ = copy(source, info.Data.Durl)
	sort.Slice(source, func(i, j int) bool {
		return source[i].Order < source[j].Order
	})
	ob := rxgo.Just(source)().Map(func(ctx context.Context, i interface{}) (interface{}, error) {
		segment := i.(VideoSegments)
		segment.targetPath = path
		return info.downloadSegment(&segment, client, ctx)
	})
	ob.DoOnNext(func(i interface{}) {

	})
	return ob.Run(opts...)
}

func (info VideoStreamInfo) downloadSegment(seg *VideoSegments, client *http.Client, ctx context.Context) (string, error) {
	ob := seg.Fetch(client, rxgo.WithContext(ctx))
	ob.DoOnNext(func(i interface{}) {

	})
	return "", nil
}
