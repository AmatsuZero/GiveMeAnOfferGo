package bilibili_api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/reactivex/rxgo/v2"
	"net/http"
	"net/url"
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
		Durl              []struct {
			Order     int // 分段序号
			Length    int
			Size      int // 单位为Byte
			Ahead     string
			Vhead     string
			Url       string   // 视频流url
			BackupUrl []string `json:"backup_url"`
		}
	}
}

func (info VideoStreamInfo) SupportFormats() []string {
	return strings.Split(info.Data.AcceptFormat, ",")
}

func (info VideoStreamInfo) Download(client *http.Client, opts ...rxgo.Option) rxgo.Observable {
	return rxgo.Empty()
}
