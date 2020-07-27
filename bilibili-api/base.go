package bilibili_api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/reactivex/rxgo/v2"
	"io"
	"net/http"
	"net/url"
)

var (
	kBaseURL, _        = url.Parse("http://api.bilibili.com")
	kInvalidParamError = fmt.Errorf("invalid param")
)

type BadResponseError struct {
	Code    int
	Message string
}

func (be BadResponseError) Error() string {
	return be.Message
}

type BaseResponse struct {
	Code    int
	Message string
	Ttl     int
}

func (resp BaseResponse) IsSuccess() bool {
	return resp.Code == 0
}

func (resp BaseResponse) GetError() error {
	if resp.IsSuccess() {
		return nil
	}
	return BadResponseError{
		Code:    resp.Code,
		Message: resp.Message,
	}
}

type BaseRequest interface {
	Request() (*http.Request, error)
	IsParamsValid() bool
	Fetch(client *http.Client, opts ...rxgo.Option) rxgo.Observable
}

type baseRequest struct{}

func (bq baseRequest) Request() (*http.Request, error) {
	return nil, kInvalidParamError
}

func (bq baseRequest) IsParamsValid() bool {
	return false
}

func (bq baseRequest) fetch(client *http.Client, req *http.Request, opts ...rxgo.Option) rxgo.Observable {
	return rxgo.Defer([]rxgo.Producer{func(ctx context.Context, next chan<- rxgo.Item) {
		if client == nil {
			client = http.DefaultClient
		}
		resp, err := client.Do(req)
		if err != nil {
			next <- rxgo.Error(err)
			return
		}
		next <- rxgo.Of(resp)
	}}, opts...).Map(func(ctx context.Context, i interface{}) (interface{}, error) {
		var info BaseResponse
		resp := i.(*http.Response)
		defer func() {
			_ = resp.Body.Close()
		}()
		var dst bytes.Buffer
		src := io.TeeReader(resp.Body, bq)
		_, err := io.Copy(&dst, src)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(dst.Bytes(), &info)
		if err != nil {
			return nil, err
		}
		if !info.IsSuccess() {
			return nil, info.GetError()
		}
		return dst.Bytes(), nil
	})
}

func (bq baseRequest) Write(p []byte) (n int, err error) {
	n = len(p)
	return
}

func (bq baseRequest) Fetch(client *http.Client, opts ...rxgo.Option) rxgo.Observable {
	req, err := bq.Request()
	if err != nil {
		return rxgo.Thrown(err)
	}
	return bq.fetch(client, req, opts...)
}

func (bq baseRequest) queryItems(query url.Values) url.Values {
	return query
}
