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

type ImageDownloaderConfig struct {
	MaxConcurrentDownloads int                           // 最大并发下载数，默认6
	DownloadTimeout        float64                       // 下载超时数，默认15s
	Client                 *http.Client                  // 自定义下载设置
	ExecutionOrder         ImageDownloaderExecutionOrder // 下载顺序，默认是FIFO
	GetDownloadOperation   func(req *http.Request, client *http.Client,
		options ImageDownloaderOptions, ctx ImageContext) ImageDownloadOperationProtocol
}

func NewImageDownloaderConfig() *ImageDownloaderConfig {
	return &ImageDownloaderConfig{
		MaxConcurrentDownloads: 6,
		DownloadTimeout:        15,
		Client:                 http.DefaultClient,
		ExecutionOrder:         ImageDownloaderFIFOExecutionOrder,
		GetDownloadOperation: func(req *http.Request, client *http.Client,
			options ImageDownloaderOptions, ctx ImageContext) ImageDownloadOperationProtocol {
			return NewImageDownloadOperation(req, client, options, ctx)
		},
	}
}

func (config *ImageDownloaderConfig) Copy() *ImageDownloaderConfig {
	return &ImageDownloaderConfig{
		MaxConcurrentDownloads: config.MaxConcurrentDownloads,
		DownloadTimeout:        config.DownloadTimeout,
		Client:                 config.Client,
		ExecutionOrder:         config.ExecutionOrder,
		GetDownloadOperation:   config.GetDownloadOperation,
	}
}

var defaultDownloaderConfig *ImageDownloaderConfig
var defaultDownloaderConfigToken sync.Once

func GetDefaultDownloaderConfig() *ImageDownloaderConfig {
	if defaultDownloaderConfig == nil {
		defaultDownloaderConfigToken.Do(func() {
			defaultDownloaderConfig = NewImageDownloaderConfig()
		})
	}
	return defaultDownloaderConfig
}
