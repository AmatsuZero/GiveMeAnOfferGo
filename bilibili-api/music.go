package bilibili_api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/reactivex/rxgo/v2"
	"net/http"
	"net/url"
)

type MusicQuality int

const (
	MusicQuality128K MusicQuality = iota
	MusicQuality192K
	MusicQuality320K
	MusicQualityFLAC
	MusicQualityTrial = -1
)

type MusicStreamRequest struct {
	baseRequest
	SongId   string
	Quality  MusicQuality
	Mid      string
	Platform string
}

func (request MusicStreamRequest) Request() (*http.Request, error) {
	if !request.IsParamsValid() {
		return nil, kInvalidParamError
	}
	base := *kBaseURL
	base.Path = "audio/music-service-c/url"
	query := request.queryItems(base.Query())
	if query.Get("mid") == "" {
		if request.Session == nil {
			request.Mid = "293793435"
		} else {
			for _, c := range request.Session.Cookies(&base) {
				if c.Name == "DedeUserID" {
					request.Mid = c.Value
					break
				}
			}
		}
		query.Set("mid", request.Mid)
	}
	base.RawQuery = query.Encode()
	return http.NewRequest("GET", base.String(), nil)
}

func (request MusicStreamRequest) Fetch(client *http.Client, opts ...rxgo.Option) rxgo.Observable {
	req, err := request.Request()
	if err != nil {
		return rxgo.Thrown(err)
	}
	return request.fetch(client, req, opts...).Map(func(ctx context.Context, i interface{}) (interface{}, error) {
		var info MusicStreamInfo
		data := i.([]byte)
		err := json.Unmarshal(data, &info)
		if err != nil {
			return nil, err
		}
		info.Session = request.Session
		return info, nil
	})
}

func (request MusicStreamRequest) queryItems(query url.Values) url.Values {
	query.Set("songid", request.SongId)
	query.Set("quality", fmt.Sprintf("%d", request.Quality))
	query.Set("privilege", "2")
	query.Set("platform", "ios")
	query.Set("mid", request.Mid)
	return query
}

func (request MusicStreamRequest) IsParamsValid() bool {
	return len(request.SongId) > 0
}

type MusicQualityInfo struct {
	Type        MusicQuality
	Desc        string
	Size        int
	Bps         string
	Tag         string
	Require     string
	RequireDesc string `json:"requiredesc"`
}

type MusicStreamInfo struct {
	BaseResponse
	baseRequest
	Data struct {
		Sid       int
		Type      MusicQuality
		Info      string `json:"info, omitempty"`
		Timeout   int
		Size      int
		Cdns      []string
		Qualities []MusicQualityInfo
		Title     string
		Cover     string
	}
	progressCB func(progress float64)
	retryCount int
}

func (info *MusicStreamInfo) SetProgressFunc(cb func(progress float64)) {
	info.progressCB = cb
}

func (info *MusicStreamInfo) Download(to string, client *http.Client, opts ...rxgo.Option) rxgo.OptionalSingle {
	return info.download(to, client, opts...).First()
}

func (info *MusicStreamInfo) download(to string, client *http.Client, opts ...rxgo.Option) rxgo.Observable {
	return info.Fetch(client, opts...).Retry(len(info.Data.Cdns), func(err error) bool {
		if e, ok := err.(BadResponseError); ok && e.Code == http.StatusForbidden { // 禁止请求，不必重试了
			return false
		}
		info.retryCount++
		return true
	}).Map(func(ctx context.Context, i interface{}) (interface{}, error) {
		return nil, nil
	})
}

func (info *MusicStreamInfo) Fetch(client *http.Client, opts ...rxgo.Option) rxgo.Observable {
	req, err := info.Request()
	if err != nil {
		return rxgo.Thrown(err)
	}
	return info.fetch(client, req, opts...)
}

func (info *MusicStreamInfo) fetch(client *http.Client, req *http.Request, opts ...rxgo.Option) rxgo.Observable {
	if client == nil {
		client = http.DefaultClient
	}
	if info.Session == nil {
		info.Session = kDefaultSession
	}
	if client.Jar == nil {
		for _, c := range info.Session.Cookies(req.URL) {
			req.AddCookie(c)
		}
	}
	return rxgo.Defer([]rxgo.Producer{func(ctx context.Context, next chan<- rxgo.Item) {
		req = req.WithContext(ctx)
		resp, err := client.Do(req)
		if err != nil {
			next <- rxgo.Error(err)
			return
		}
		defer func() {
			if client.Jar != nil {
				return
			}
			kDefaultSession.SetCookies(resp.Request.URL, resp.Cookies())
			_ = kDefaultSession.Serialize(kDefaultSessionPath)
		}()
		if resp.StatusCode == http.StatusForbidden {
			next <- rxgo.Of(BadResponseError{
				Code:    resp.StatusCode,
				Message: resp.Status,
			})
			_ = resp.Body.Close()
			return
		}
		next <- rxgo.Of(resp)
	}}, opts...)
}

func (info *MusicStreamInfo) Request() (*http.Request, error) {
	if !info.IsParamsValid() {
		return nil, kInvalidParamError
	}
	if info.retryCount > len(info.Data.Cdns) {
		return nil, kOutOfTimesError
	}
	u, err := url.Parse(info.Data.Cdns[info.retryCount])
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", kFakeWebUA)
	return req, nil
}

func (info *MusicStreamInfo) IsParamsValid() bool {
	return !info.IsSuccess() || len(info.Data.Cdns) > 0
}
