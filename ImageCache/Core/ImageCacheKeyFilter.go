package Core

import "net/url"

type ImageCacheKeyFilterBlock func(url *url.URL) string

type ImageCacheKeyFilterProtocol interface {
	CacheKeyForURL(url *url.URL) string
}

type ImageCacheKeyFilter struct {
	cb ImageCacheKeyFilterBlock
}

func (filter *ImageCacheKeyFilter) CacheKeyForURL(url *url.URL) string {
	if filter == nil || filter.cb == nil {
		return ""
	}
	return filter.cb(url)
}

func NewCacheKeyFilter(cb ImageCacheKeyFilterBlock) *ImageCacheKeyFilter {
	return &ImageCacheKeyFilter{cb: cb}
}
