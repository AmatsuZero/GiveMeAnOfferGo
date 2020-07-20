package Core

import (
	"fmt"
	lru "github.com/hashicorp/golang-lru"
	"io/ioutil"
	"net/http"
	"sync"
)

type BitsType int

func BitsSet(options, flag BitsType) BitsType {
	return options | flag
}

func BitsClear(options, flag BitsType) BitsType {
	return options &^ flag
}

func BitsToggle(options, flag BitsType) BitsType {
	return options ^ flag
}

func BitsHas(options, flag BitsType) bool {
	return options&flag != 0
}

type ImageContext map[string]interface{}

var InvalidParamError error

func init() {
	InvalidParamError = fmt.Errorf("check input")
}

type URLCache struct {
	*lru.Cache
}

type CachedResponse struct {
	*http.Response
	BodyData []byte
}

func NewCachedResponse(resp *http.Response) *CachedResponse {
	if resp == nil {
		return nil
	}
	data, _ := ioutil.ReadAll(resp.Body)
	return &CachedResponse{
		Response: resp,
		BodyData: data,
	}
}

func (cache *URLCache) AddCache(req *http.Request, response *http.Response) bool {
	if req == nil || cache == nil {
		return false
	}
	return cache.Add(req, NewCachedResponse(response))
}

func (cache *URLCache) GetCachedResponseForRequest(req *http.Request) (resp *CachedResponse, ok bool) {
	if req == nil || cache == nil {
		return nil, false
	}
	r, ok := cache.Get(req)
	if !ok {
		return nil, false
	}
	resp, ok = r.(*CachedResponse)
	return
}

var kDefaultURLCache *URLCache
var kDefaultURLCacheToken sync.Once

func getURLCache() *URLCache {
	if kDefaultURLCache == nil {
		kDefaultURLCacheToken.Do(func() {
			cache, _ := lru.New(2000)
			kDefaultURLCache = &URLCache{cache}
		})
	}
	return kDefaultURLCache
}
