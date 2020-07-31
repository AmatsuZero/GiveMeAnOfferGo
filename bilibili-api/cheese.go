package bilibili_api

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"github.com/reactivex/rxgo/v2"
	"net/http"
	"net/url"
)

type CheeseVideoStreamRequest struct {
	videoStreamBaseRequest
	Avid        string
	EpId        string
	Cid         string
	Resolution  VideoStreamResolution
	ShouldUse4K bool
}

func (request CheeseVideoStreamRequest) IsParamsValid() bool {
	return len(request.Avid) > 0 && len(request.EpId) > 0 && len(request.Cid) > 0
}

func (request CheeseVideoStreamRequest) Request() (*http.Request, error) {
	if !request.IsParamsValid() {
		return nil, kInvalidParamError
	}
	base := *kBaseURL
	base.Path = "pugv/player/web/playurl"
	base.RawQuery = request.queryItems(base.Query()).Encode()
	return http.NewRequest("GET", base.String(), nil)
}

func (request CheeseVideoStreamRequest) queryItems(query url.Values) url.Values {
	query.Set("avid", request.Avid)
	query.Set("cid", request.Cid)
	query.Set("ep_id", request.EpId)
	if request.Resolution > 0 {
		query.Set("qn", request.Resolution.Encode())
	}
	if request.ShouldUse4K {
		query.Set("fourk", "1")
	}
	return query
}

func (request CheeseVideoStreamRequest) Download(to string, client *http.Client, opts ...rxgo.Option) rxgo.OptionalSingle {
	defaultClient := http.DefaultClient // 这里改用没有超时的默认 Client，避免任务被 Cancel
	ob := request.download(request.Fetch(client), defaultClient, opts...)
	return request.export(ob, to)
}

func (request CheeseVideoStreamRequest) Fetch(client *http.Client, opts ...rxgo.Option) rxgo.Observable {
	req, err := request.Request()
	if err != nil {
		return rxgo.Thrown(err)
	}
	return request.fetch(client, req, opts...).Map(func(ctx context.Context, i interface{}) (interface{}, error) {
		info := &CheeseVideoStreamInfo{} // interface 要使用，必须 cast 为指针类型，参考：https://forum.golangbridge.org/t/how-to-cast-interface-to-a-given-interface/13997
		data := i.([]byte)
		err := json.Unmarshal(data, info)
		if err != nil {
			return nil, err
		}
		info.md5Sign = md5.Sum(data)
		return info, nil
	})
}

type CheeseVideoStreamInfo struct {
	BaseResponse
	baseVideoStreamInfo
	Data struct {
		VideoStreamInfoData
		Code                int
		NoRexcode           int
		Fnval               int
		VideoProject        bool
		Fnver               int
		Type                string
		HasPaid             bool
		SupportFormatsInfos []struct {
			Formt       string
			Description string
			Quality     VideoStreamResolution
		} `json:"supportFormats"`
		Status int
	}
}

func (info CheeseVideoStreamInfo) TotalSize() (size int64) {
	for _, seg := range info.GetSegments() {
		size += int64(seg.Size)
	}
	return
}

func (info CheeseVideoStreamInfo) GetSegments() []VideoSegments {
	return info.Data.Durl
}

func (info CheeseVideoStreamInfo) download(tmpDir string, client *http.Client, opts ...rxgo.Option) rxgo.Observable {
	return info.operation(info.GetSegments(), info.TotalSize(), tmpDir, client, opts...)
}

func (info CheeseVideoStreamInfo) Download(to string, client *http.Client, opts ...rxgo.Option) rxgo.OptionalSingle {
	return info.download(to, client, opts...).First()
}
