package bilibili_api

import (
	"context"
	"encoding/json"
	"github.com/reactivex/rxgo/v2"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	kPassportBaseURL, _ = url.Parse("http://passport.bilibili.com")
	kDefaultSessionPath = ""
	kDefaultSession     *Session
)

func init() {
	config, err := os.UserConfigDir()
	if err != nil {
		return
	}
	kDefaultSessionPath = filepath.Join(config, "user_session")
	kDefaultSession = NewSession(kDefaultSessionPath)
}

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

func (request LoginRequest) Login(client *http.Client, opts ...rxgo.Option) (s *Session) {
	s = NewSession("")
	for item := range request.Fetch(client).Observe(opts...) {
		if item.E == nil && item.V.(*Session).IsSuccess() {
			s = item.V.(*Session)
			break
		}
	}
	if s.IsSuccess() {
		_ = s.Serialize(kDefaultSessionPath)
		kDefaultSession = s // 更新 Session
	}
	return
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
		err = exec.CommandContext(ctx, "python", "-m", "webbrowser", "-t", param.QRLink()).Run()
		if err != nil {
			return nil, err
		}
		return param, nil
	}).FlatMap(func(item rxgo.Item) rxgo.Observable {
		if item.E != nil {
			return rxgo.Thrown(item.E)
		}
		start := time.Now() // 180秒计时
		return item.V.(sessionParam).Fetch(client, opts...).TakeUntil(func(i interface{}) bool {
			return i.(*Session).IsSuccess() || time.Since(start).Seconds() >= 180
		})
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

func (param sessionParam) QRLink() string {
	base, _ := url.Parse("https://cli.im/api/qrcode/code")
	q := base.Query()
	q.Set("text", param.Data.Url)
	base.RawQuery = q.Encode()
	return base.String()
}

func (param sessionParam) Request() (*http.Request, error) {
	if !param.IsParamsValid() {
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
	return rxgo.Interval(rxgo.WithDuration(3*time.Second)).Map(func(ctx context.Context, i interface{}) (interface{}, error) {
		req = req.WithContext(ctx)
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer func() {
			_ = resp.Body.Close()
		}()
		session := NewSession("")
		err = json.NewDecoder(resp.Body).Decode(session)
		if err != nil {
			return nil, err
		}
		if session.IsSuccess() {
			session.SetCookies(resp.Request.URL, resp.Cookies()) // 持久化 Cookie
		}
		return session, nil
	}, opts...)
}

type Session struct {
	BaseResponse
	Status     bool
	Data       interface{}
	cookieJar  map[*url.URL][]*http.Cookie
	cookieLock sync.Mutex
}

func NewSession(path string) (s *Session) {
	s = &Session{}
	if len(path) == 0 {
		return
	}
	info, err := os.Stat(path)
	if os.IsNotExist(err) || info.IsDir() {
		return
	}
	file, err := os.Open(path)
	defer func() {
		_ = file.Close()
	}()
	if err != nil {
		return
	}
	err = json.NewDecoder(file).Decode(s)
	return
}

func (s *Session) IsSuccess() bool {
	return s.Status
}

func (s *Session) RedirectURL() string {
	if !s.IsSuccess() {
		return ""
	}
	r, ok := s.Data.(map[string]interface{})
	if !ok {
		return ""
	}
	ret := r["url"]
	return ret.(string)
}

func (s *Session) SetCookies(u *url.URL, cookies []*http.Cookie) {
	s.cookieLock.Lock()
	if s.cookieJar == nil {
		s.cookieJar = map[*url.URL][]*http.Cookie{}
	}
	defer s.cookieLock.Unlock()
	jar, ok := s.cookieJar[u]
	if ok {
		jar = append(jar, cookies...)
		set := map[*http.Cookie]bool{}
		for _, c := range jar {
			set[c] = true
		}
		jar = jar[:0]
		for c := range set {
			jar = append(jar, c)
		}
		jar = append(jar, cookies...)
	} else {
		jar = cookies
	}
	s.cookieJar[u] = jar
}

func (s *Session) Cookies(u *url.URL) []*http.Cookie {
	s.cookieLock.Lock()
	defer s.cookieLock.Unlock()
	if s.cookieJar == nil {
		s.cookieJar = map[*url.URL][]*http.Cookie{}
	}
	return s.cookieJar[u]
}

func (s *Session) Serialize(path string) error {
	byteArr, err := json.Marshal(s)
	if err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()
	_, err = file.Write(byteArr)
	return err
}
