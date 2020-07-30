package bilibili_api

import (
	"archive/zip"
	"compress/flate"
	"compress/gzip"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/reactivex/rxgo/v2"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

type HistoryDanmukuRequest struct {
	baseRequest
	Oid        string
	Date       string
	progressCb func(progress float64)
	fileSize   int64
}

func (request HistoryDanmukuRequest) Request() (*http.Request, error) {
	if !request.IsParamsValid() {
		return nil, kInvalidParamError
	}
	base := *kBaseURL
	base.Path = "x/v2/dm/history"
	base.RawQuery = request.queryItems(base.Query()).Encode()
	return http.NewRequest("GET", base.String(), nil)
}

func (request *HistoryDanmukuRequest) SetProgressFunc(cb func(progress float64)) {
	request.progressCb = cb
}

func (request HistoryDanmukuRequest) Download(to string, client *http.Client, opts ...rxgo.Option) rxgo.OptionalSingle {
	return request.download(to, client, opts...).First()
}

func (request HistoryDanmukuRequest) download(to string, client *http.Client, opts ...rxgo.Option) rxgo.Observable {
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
		receivedSize := 0
		_, e = io.Copy(file, readerFunc(func(p []byte) (n int, err error) {
			defer func() {
				receivedSize += n
				if request.progressCb != nil {
					request.progressCb(float64(receivedSize) / float64(request.fileSize))
				}
			}()
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
	})
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
		request.fileSize = r.ContentLength
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

type HistoryDanmukuIndexRequest struct {
	baseRequest
	Oid        string
	Month      string
	progressCB func(progress float64)
}

func (request HistoryDanmukuIndexRequest) queryItems(query url.Values) url.Values {
	query.Set("type", "1")
	query.Set("oid", request.Oid)
	query.Set("month", request.Month)
	return query
}

func (request HistoryDanmukuIndexRequest) IsParamsValid() bool {
	return len(request.Month) > 0 && len(request.Oid) > 0
}

func (request HistoryDanmukuIndexRequest) Request() (*http.Request, error) {
	if !request.IsParamsValid() {
		return nil, kInvalidParamError
	}
	base := *kBaseURL
	base.Path = "x/v2/dm/history/index"
	base.RawQuery = request.queryItems(base.Query()).Encode()
	return http.NewRequest("GET", base.String(), nil)
}

func (request *HistoryDanmukuIndexRequest) Fetch(client *http.Client, opts ...rxgo.Option) rxgo.Observable {
	req, err := request.Request()
	if err != nil {
		return rxgo.Thrown(err)
	}
	return request.baseRequest.fetch(client, req, opts...).Map(func(ctx context.Context, i interface{}) (interface{}, error) {
		var info HistoryDanmukuIndex
		data := i.([]byte)
		err := json.Unmarshal(data, &info)
		info.md5Sign = md5.Sum(data)
		if err != nil {
			return nil, err
		}
		return info, nil
	})
}

func (request *HistoryDanmukuIndexRequest) SetProgressFunc(cb func(progress float64)) {
	request.progressCB = cb
}

func (request *HistoryDanmukuIndexRequest) Download(to string, client *http.Client, opts ...rxgo.Option) rxgo.OptionalSingle {
	tmpDir := ""
	var zipFile *os.File
	totalProgress := float64(0)
	count := 0
	ob := request.Fetch(client, opts...).FlatMap(func(item rxgo.Item) rxgo.Observable {
		if item.E != nil {
			return rxgo.Thrown(item.E)
		}
		tmpDir = filepath.Join(os.TempDir(), fmt.Sprintf("%x", item.V.(HistoryDanmukuIndex).md5Sign))
		info, err := os.Stat(tmpDir)
		if os.IsNotExist(err) {
			err = os.MkdirAll(tmpDir, os.ModePerm)
		}
		if err != nil {
			return rxgo.Thrown(err)
		}

		info, err = os.Stat(to)
		if info != nil && info.IsDir() { // 如果是文件夹，则创建一个压缩包文件
			to = filepath.Join(to, fmt.Sprintf("%v.zip", request.Month))
		}
		if filepath.Ext(to) != ".zip" { // 如果后缀不是zip，则拼接一个zip
			to += ".zip"
		}
		zipFile, err = os.Create(to)
		if err != nil {
			return rxgo.Thrown(err)
		}
		dates := item.V.(HistoryDanmukuIndex).Data
		if len(dates) == 0 {
			return rxgo.Empty()
		}
		count = len(dates)
		return rxgo.Just(dates)()
	}).Map(func(ctx context.Context, i interface{}) (interface{}, error) {
		req := HistoryDanmukuRequest{}
		req.Session = request.Session
		req.Oid = request.Oid
		req.Date = i.(string)
		req.SetProgressFunc(func(progress float64) {
			totalProgress += progress / float64(count)
			if req.progressCb != nil {
				req.progressCb(totalProgress)
			}
		})
		return req, nil
	}).FlatMap(func(item rxgo.Item) rxgo.Observable {
		req := item.V.(HistoryDanmukuRequest)
		path := filepath.Join(tmpDir, req.Date+".xml")
		return req.download(path, client, opts...)
	})

	return ob.Reduce(func(ctx context.Context, acc interface{}, elem interface{}) (interface{}, error) {
		var archive *zip.Writer
		if acc == nil {
			archive = zip.NewWriter(zipFile)
		} else {
			archive = acc.(*zip.Writer)
		}
		info, err := os.Stat(elem.(string))
		if err != nil {
			return nil, err
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return nil, err
		}
		writer, err := archive.CreateHeader(header)
		if err != nil {
			return nil, err
		}
		file, err := os.Open(elem.(string))
		if err != nil {
			return nil, err
		}
		defer func() {
			_ = file.Close()
		}()
		_, err = io.Copy(writer, file)
		return archive, nil
	}).Map(func(ctx context.Context, i interface{}) (interface{}, error) {
		archive := i.(*zip.Writer)
		err := archive.Close()
		_ = os.RemoveAll(tmpDir)
		_ = zipFile.Close()
		return to, err
	})
}

type HistoryDanmukuIndex struct {
	BaseResponse
	Data    []string
	md5Sign [16]byte
}
