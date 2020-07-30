package bilibili_api

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/reactivex/rxgo/v2"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
)

type MusicQuality int

const (
	MusicQuality128K MusicQuality = iota
	MusicQuality192K
	MusicQuality320K
	MusicQualityFLAC
	MusicQualityTrial MusicQuality = -1
)

func (quality MusicQuality) String() string {
	switch quality {
	case MusicQuality128K:
		return "流畅 128K"
	case MusicQuality192K:
		return "标准 192K"
	case MusicQualityFLAC:
		return "高品质 320K"
	case MusicQualityTrial:
		return "无损 FLAC （大会员）"
	}
	return "未知"
}

type MusicStreamRequest struct {
	baseRequest
	SongId     string
	Quality    MusicQuality
	Mid        string
	Platform   string
	progressCB func(progress float64)
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

func (request *MusicStreamRequest) SetProgressFunc(cb func(progress float64)) {
	request.progressCB = cb
}

func (request *MusicStreamRequest) Download(to string, client *http.Client, opts ...rxgo.Option) rxgo.OptionalSingle {
	return request.Fetch(client, opts...).Map(func(ctx context.Context, i interface{}) (interface{}, error) {
		info := i.(MusicStreamInfo)
		info.Session = request.Session
		info.progressCB = request.progressCB
		return info.Download(to, client, opts...).Get()
	}).First()
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
		info.md5Sign = md5.Sum(data)
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
	Require     int
	RequireDesc string `json:"requiredesc"`
}

func (info MusicQualityInfo) NeedVIP() bool {
	return info.Require == 1
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
	md5Sign    [16]byte
}

func (info *MusicStreamInfo) SetProgressFunc(cb func(progress float64)) {
	info.progressCB = cb
}

func (info *MusicStreamInfo) Download(to string, client *http.Client, opts ...rxgo.Option) rxgo.OptionalSingle {
	tmpDir := filepath.Join(os.TempDir(), fmt.Sprintf("%x", info.md5Sign))
	_, err := os.Stat(tmpDir)
	if os.IsNotExist(err) {
		err = os.Mkdir(tmpDir, os.ModePerm)
	}
	if err != nil {
		return rxgo.Thrown(err).First()
	}
	tmpCoverPath := filepath.Join(tmpDir, info.Data.Title+filepath.Ext(info.Data.Cover))
	tmpMusicPath := filepath.Join(tmpDir, "music.mp3")
	if info.Data.Type == MusicQualityFLAC {
		tmpMusicPath = filepath.Join(tmpDir, "music.aac")
	}
	return rxgo.Merge([]rxgo.Observable{
		info.downloadCover(tmpCoverPath, client, opts...),
		info.download(tmpMusicPath, client, opts...),
	}).Reduce(func(ctx context.Context, i interface{}, i2 interface{}) (interface{}, error) {
		if i == nil {
			return i2, nil
		}
		return i.(string) + i2.(string), nil
	}).Map(func(ctx context.Context, i interface{}) (interface{}, error) {
		cmd := exec.CommandContext(ctx, "ffmpeg", "-i", tmpMusicPath, "-i", tmpCoverPath,
			"-map_metadata", "0", "-map", "0", "-map", "1", to)
		err = cmd.Run()
		fmt.Println(cmd.String())
		_ = os.RemoveAll(tmpDir)
		return to, err
	})
}

func (info *MusicStreamInfo) downloadCover(to string, client *http.Client, opts ...rxgo.Option) rxgo.Observable {
	if client == nil {
		client = http.DefaultClient
	}
	req, err := http.NewRequest("GET", info.Data.Cover, nil)
	if err != nil {
		return rxgo.Thrown(err)
	}
	return rxgo.Defer([]rxgo.Producer{func(ctx context.Context, next chan<- rxgo.Item) {
		req = req.WithContext(ctx)
		r, err := client.Do(req)
		if err != nil {
			next <- rxgo.Error(err)
			return
		}
		defer func() {
			_ = r.Body.Close()
		}()
		file, err := os.Create(to)
		if err != nil {
			next <- rxgo.Error(err)
			return
		}
		defer func() {
			_ = file.Close()
		}()
		_, err = io.Copy(file, readerFunc(func(p []byte) (n int, err error) {
			select {
			case <-ctx.Done():
				return 0, ctx.Err()
			default:
				return r.Body.Read(p)
			}
		}))
		if err != nil {
			next <- rxgo.Error(err)
			return
		}
		next <- rxgo.Of(to)
	}}, opts...)
}

func (info *MusicStreamInfo) download(to string, client *http.Client, opts ...rxgo.Option) rxgo.Observable {
	return info.Fetch(client, opts...).Retry(len(info.Data.Cdns), func(err error) bool {
		if e, ok := err.(BadResponseError); ok && e.Code == http.StatusForbidden { // 禁止请求，不必重试了
			return false
		}
		info.retryCount++
		return true
	}).Map(func(ctx context.Context, i interface{}) (interface{}, error) {
		r := i.(*http.Response)
		defer func() {
			_ = r.Body.Close()
		}()
		file, err := os.Create(to)
		if err != nil {
			return nil, err
		}
		defer func() {
			_ = file.Close()
		}()
		receivedSize := 0
		if info.Data.Size == 0 { // 试听类型，大小为 0
			info.Data.Size = int(r.ContentLength)
		}
		_, err = io.Copy(file, readerFunc(func(p []byte) (n int, err error) {
			select {
			case <-ctx.Done():
				return 0, ctx.Err()
			default:
				n, err = r.Body.Read(p)
				receivedSize += n
				if info.progressCB != nil {
					info.progressCB(float64(receivedSize) / float64(info.Data.Size))
				}
				return n, err
			}
		}))
		return to, err
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
		for _, c := range info.Session.jar.AllCookies() {
			if c.Name == "SESSDATA" {
				req.AddCookie(c)
				break
			}
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
