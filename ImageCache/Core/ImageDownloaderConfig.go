package Core

import (
	"net/http"
	"sync"
)

type ImageDownloaderExecutionOrder int

const (
	ImageDownloaderFIFOExecutionOrder ImageDownloaderExecutionOrder = iota
	SDWebImageDownloaderLIFOExecutionOrder
)

type imageDownloaderConfig struct {
	MaxConcurrentDownloads int                           // 最大并发下载数，默认6
	DownloadTimeout        float64                       // 下载超时数，默认15s
	Client                 *http.Client                  // 自定义下载设置
	ExecutionOrder         ImageDownloaderExecutionOrder // 下载顺序，默认是FIFO
	GetDownloadOperation   func(req *http.Request, client *http.Client,
		options ImageDownloaderOptions, ctx ImageContext) ImageDownloadOperation
}

func newImageDownloaderConfig() *imageDownloaderConfig {
	return &imageDownloaderConfig{
		MaxConcurrentDownloads: 6,
		DownloadTimeout:        15,
		Client:                 http.DefaultClient,
		ExecutionOrder:         ImageDownloaderFIFOExecutionOrder,
		GetDownloadOperation: func(req *http.Request, client *http.Client,
			options ImageDownloaderOptions, ctx ImageContext) ImageDownloadOperation {
			return newImageDownloadOperation(req, client, options, ctx)
		},
	}
}

func (config *imageDownloaderConfig) Copy() *imageDownloaderConfig {
	return &imageDownloaderConfig{
		MaxConcurrentDownloads: config.MaxConcurrentDownloads,
		DownloadTimeout:        config.DownloadTimeout,
		Client:                 config.Client,
		ExecutionOrder:         config.ExecutionOrder,
		GetDownloadOperation:   config.GetDownloadOperation,
	}
}

var defaultDownloaderConfig *imageDownloaderConfig
var defaultDownloaderConfigToken sync.Once

func getDefaultDownloaderConfig() *imageDownloaderConfig {
	if defaultDownloaderConfig == nil {
		defaultDownloaderConfigToken.Do(func() {
			defaultDownloaderConfig = newImageDownloaderConfig()
		})
	}
	return defaultDownloaderConfig
}
