package bilibili_api

import (
	"context"
	"encoding/json"
	"github.com/reactivex/rxgo/v2"
	"math"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
	"time"
)

var (
	kPassportBaseURL, _ = url.Parse("http://passport.bilibili.com")
)

type LoginRequest struct {
	baseRequest
}

func (request LoginRequest) IsParamsValid() bool {
	return true
}

func (request LoginRequest) Request() (*http.Request, error) {
	u := *kPassportBaseURL
	u.Path = "qrcode/getLoginUrl"
	return http.NewRequest("GET", u.String(), nil)
}

func (request LoginRequest) Fetch(client *http.Client, opts ...rxgo.Option) rxgo.Observable {
	req, err := request.Request()
	if err != nil {
		return rxgo.Thrown(err)
	}
	return request.fetch(client, req, opts...).Map(func(ctx context.Context, i interface{}) (interface{}, error) {
		var param sessionParam
		data := i.([]byte)
		err := json.Unmarshal(data, &param)
		if err != nil {
			return nil, err
		}
		// 打开链接
		err = exec.CommandContext(ctx, "python", "-m", "webbrowser", "-t", param.Data.Url).Run()
		if err != nil {
			return nil, err
		}
		return param, nil
	}).FlatMap(func(item rxgo.Item) rxgo.Observable {
		if item.E != nil {
			return rxgo.Thrown(item.E)
		}
		return item.V.(sessionParam).Fetch(client, opts...)
	})
}

type sessionParam struct {
	BaseResponse
	baseRequest
	Data struct {
		Url      string
		OauthKey string `json:"oauthKey"`
	}
}

func (param sessionParam) Request() (*http.Request, error) {
	if param.IsParamsValid() {
		return nil, kInvalidParamError
	}
	u := *kPassportBaseURL
	u.Path = "qrcode/getLoginInfo"
	data := url.Values{}
	data.Set("oauthKey", param.Data.OauthKey)
	data.Set("gourl", "http://www.bilibili.com")
	req, err := http.NewRequest("POST", u.String(), strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return req, nil
}

func (param sessionParam) IsParamsValid() bool {
	return len(param.Data.OauthKey) > 0
}

func (param sessionParam) Fetch(client *http.Client, opts ...rxgo.Option) rxgo.Observable {
	req, err := param.Request()
	if err != nil {
		return rxgo.Thrown(err)
	}
	if client == nil {
		client = http.DefaultClient
	}
	start := time.Now()
	return rxgo.Defer([]rxgo.Producer{func(ctx context.Context, next chan<- rxgo.Item) {
		var session Session
		var err error
		defer func() {
			next <- rxgo.Item{
				V: session,
				E: err,
			}
		}()
		resp, err := client.Do(req)
		if err != nil {
			return
		}
		defer func() {
			_ = resp.Body.Close()
		}()
		err = json.NewDecoder(resp.Body).Decode(&session)
		if err != nil {
			return
		}
		if !session.IsSuccess() {
			err = session.GetError()
			return
		}
		session.SetCookies(resp.Request.URL, resp.Cookies()) // 持久化 Cookie
	}}).Retry(math.MaxInt32, func(err error) bool {
		return time.Since(start).Seconds() <= 180
	})
}

type Session struct {
	BaseResponse
	Status bool
}

func (s Session) IsSuccess() bool {
	return s.Status
}

func (s Session) SetCookies(u *url.URL, cookies []*http.Cookie) {

}

func (s Session) Cookies(u *url.URL) []*http.Cookie {
	return nil
}
