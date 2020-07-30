package bilibili_api

import (
	"compress/flate"
	"compress/gzip"
	"context"
	"encoding/json"
	"github.com/reactivex/rxgo/v2"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

type HistoryDanmukuRequest struct {
	baseRequest
	Oid        string
	Date       string
	progressCb func(progress float64)
}

func (request HistoryDanmukuRequest) Request() (*http.Request, error) {
	base := *kBaseURL
	base.Path = "x/v2/dm/history"
	base.RawQuery = request.queryItems(base.Query()).Encode()
	return http.NewRequest("GET", base.String(), nil)
}

func (request *HistoryDanmukuRequest) SetProgressFunc(cb func(progress float64)) {
	request.progressCb = cb
}

func (request HistoryDanmukuRequest) Download(to string, client *http.Client, opts ...rxgo.Option) rxgo.OptionalSingle {
	return request.Fetch(client, opts...).Map(func(ctx context.Context, i interface{}) (interface{}, error) {
		r := i.(*http.Response)
		var e error
		// Check that the server actually sent compressed data
		var reader io.ReadCloser
		switch r.Header.Get("Content-Encoding") {
		case "gzip":
			reader, e = gzip.NewReader(r.Body)
			if e != nil {
				return nil, e
			}
			defer func() {
				_ = reader.Close()
			}()
		case "deflate":
			reader = flate.NewReader(r.Body)
			defer func() {
				_ = reader.Close()
			}()
		default:
			reader = r.Body
		}
		file, e := os.Create(to)
		if e != nil {
			return nil, e
		}
		defer func() {
			_ = file.Close()
		}()
		_, e = io.Copy(file, readerFunc(func(p []byte) (int, error) {
			select {
			case <-ctx.Done():
				return 0, ctx.Err()
			default:
				return reader.Read(p)
			}
		}))
		if e != nil {
			return nil, e
		}
		return to, e
	}).First()
}

func (request *HistoryDanmukuRequest) Fetch(client *http.Client, opts ...rxgo.Option) rxgo.Observable {
	req, err := request.Request()
	if err != nil {
		return rxgo.Thrown(err)
	}
	if client == nil {
		client = http.DefaultClient
	}
	if request.Session == nil {
		request.Session = kDefaultSession
	}
	if client.Jar == nil {
		for _, c := range request.Session.Cookies(req.URL) {
			req.AddCookie(c)
		}
	}
	return rxgo.Defer([]rxgo.Producer{func(ctx context.Context, next chan<- rxgo.Item) {
		req = req.WithContext(ctx)
		r, err := client.Do(req)
		if err != nil {
			next <- rxgo.Error(err)
			return
		}
		defer func() {
			if client.Jar != nil {
				return
			}
			kDefaultSession.SetCookies(r.Request.URL, r.Cookies())
			_ = kDefaultSession.Serialize(kDefaultSessionPath)
		}()
		if r.Header.Get("Content-Type") != "text/xml" {
			data, err := ioutil.ReadAll(r.Body)
			defer func() {
				next <- rxgo.Error(err)
			}()
			if err != nil {
				return
			}
			var info BaseResponse
			err = json.Unmarshal(data, &info)
			if err != nil {
				return
			}
			err = info.GetError()
			return
		}
		next <- rxgo.Of(r)
	}}, opts...)
}

func (request HistoryDanmukuRequest) IsParamsValid() bool {
	return len(request.Date) > 0 && len(request.Oid) > 0
}

func (request HistoryDanmukuRequest) queryItems(query url.Values) url.Values {
	query.Set("type", "1")
	query.Set("oid", request.Oid)
	query.Set("date", request.Date)
	return query
}
