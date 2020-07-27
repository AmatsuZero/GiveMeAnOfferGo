package bilibili_api

import (
	"fmt"
	"net/url"
)

var (
	kBaseURL, _        = url.Parse("http://api.bilibili.com")
	kInvalidParamError = fmt.Errorf("invalid param")
)

func GetVideoDetailURLByAID(aid string) (*url.URL, error) {
	return getVideoDetailURL(aid, "")
}

func GetVideoDetailURLByBvid(id string) (*url.URL, error) {
	return getVideoDetailURL("", id)
}

func getVideoDetailURL(aid, bvid string) (*url.URL, error) {
	if len(aid) == 0 && len(bvid) == 0 {
		return nil, kInvalidParamError
	}
	base := *kBaseURL
	base.Path = "x/web-interface/view"
	query := base.Query()
	if len(aid) > 0 {
		query.Add("aid", aid)
	} else {
		query.Add("bvid", bvid)
	}
	base.RawQuery = query.Encode()
	return &base, nil
}

type BaseResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Ttl     int    `json:"ttl"`
}
