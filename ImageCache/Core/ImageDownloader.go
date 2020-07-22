package Core

import (
	"net/http"
	"sync"
)

type ImageDownloaderOptions BitsType

const (
	ImageDownloaderLowPriority ImageDownloaderOptions = 1 << iota
	ImageDownloaderProgressiveLoad
	ImageDownloaderUseNSURLCache
	ImageDownloaderIgnoreCachedResponse
	ImageDownloaderContinueInBackground
	ImageDownloaderHandleCookies
	ImageDownloaderHighPriority
	ImageDownloaderScaleDownLargeImages
	ImageDownloaderAvoidDecodeImage
	ImageDownloaderDecodeFirstFrameOnly
	ImageDownloaderPreloadAllFrames
	ImageDownloaderMatchAnimatedImageClass
)

type ImageDownloadToken struct {
	downloadOperationToken interface{}
	downloadOperation      ImageDownloadOperationProtocol
	lock                   sync.Mutex
}

func newImageDownloadTokenWithObserver(op ImageDownloadOperationProtocol) *ImageDownloadToken {
	return &ImageDownloadToken{downloadOperation: op}
}

func (token *ImageDownloadToken) GetRequest() *http.Request {
	if token == nil {
		return nil
	}
	return nil
}

func (token *ImageDownloadToken) GetResponse() *http.Response {
	if token == nil {
		return nil
	}
	return token.downloadOperation.GetResponse()
}

func (token *ImageDownloadToken) Cancel() {
	if token == nil || token.downloadOperation == nil || token.downloadOperation.GetIsCanceled() {
		return
	}
	token.downloadOperation.CancelWithToken(token.downloadOperationToken)
	token.downloadOperationToken = nil
}

type ImageDownloader struct {
	config               *ImageDownloaderConfig
	ResponseModifier     ImageDownloaderResponseModifier
	RequestModifier      ImageDownloaderRequestModifier
	Decryptor            ImageDownloaderDecryptor
	CurrentDownloadCount uint
	operationLock        sync.Mutex
	httpHeaderLock       sync.Mutex
	client               *http.Client
	httpHeaders          map[string]string
	downloadQueue        *ImageDownloadQueue
	urlOperations        map[string]ImageDownloadOperationProtocol
}

func (downloader *ImageDownloader) CreateDownloaderOperation(url string, ctx ImageContext, options ImageDownloaderOptions) (ImageDownloadOperationProtocol, error) {
	if downloader == nil {
		return nil, ImageOperationInvalidError
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	downloader.httpHeaderLock.Lock()
	for k, v := range downloader.httpHeaders {
		req.Header.Set(k, v)
	}
	downloader.httpHeaderLock.Unlock()
	contextDict := ctx
	if contextDict == nil {
		contextDict = ImageContext{}
	} else {
		for k, v := range ctx {
			contextDict[k] = v
		}
	}
	// req modifier
	var reqModifier ImageDownloaderRequestModifier
	if modifier, ok := contextDict[kImageContextDownloadResponseModifier]; ok {
		if m, ok := modifier.(ImageDownloaderRequestModifier); ok {
			reqModifier = m
		}
	}
	if reqModifier == nil {
		reqModifier = downloader.RequestModifier
	}
	if reqModifier != nil {
		contextDict[kImageContextDownloadResponseModifier] = reqModifier
	}
	// resp modifier
	var respModifier ImageDownloaderResponseModifier
	if modifier, ok := contextDict[kImageContextDownloadResponseModifier]; ok {
		if m, ok := modifier.(ImageDownloaderResponseModifier); ok {
			respModifier = m
		}
	}
	if respModifier == nil {
		respModifier = downloader.ResponseModifier
	}
	if respModifier != nil {
		contextDict[kImageContextDownloadResponseModifier] = respModifier
	}
	var op ImageDownloadOperationProtocol
	if downloader.GetConfig() != nil && downloader.GetConfig().GetDownloadOperation != nil {
		op = downloader.GetConfig().GetDownloadOperation(req, downloader.client, options, contextDict)
	} else {
		op = NewImageDownloadOperation(req, downloader.client, options, contextDict)
	}
	// Decryptor
	var decryptor ImageDownloaderDecryptor
	if modifier, ok := contextDict[kImageContextDownloadDecryptor]; ok {
		if m, ok := modifier.(ImageDownloaderDecryptor); ok {
			decryptor = m
		}
	}
	if decryptor == nil {
		decryptor = downloader.Decryptor
	}
	if decryptor != nil {
		contextDict[kImageContextDownloadResponseModifier] = decryptor
	}
	if BitsHas(BitsType(options), BitsType(ImageDownloaderHighPriority)) {
		op.SetOperationPriority(WebImageOperationPriorityHigh)
	} else if BitsHas(BitsType(options), BitsType(ImageDownloaderLowPriority)) {
		op.SetOperationPriority(WebImageOperationPriorityLow)
	}
	if downloader.GetConfig() != nil && downloader.GetConfig().ExecutionOrder == SDWebImageDownloaderLIFOExecutionOrder {
		for _, pendingOp := range downloader.downloadQueue.GetOperations() {
			pendingOp.AddDependency(op)
		}
	}
	return op, nil
}

func (downloader *ImageDownloader) SetValueForHTTPField(value, field string) {
	if downloader == nil {
		return
	}
	downloader.httpHeaderLock.Lock()
	downloader.httpHeaders[field] = value
	downloader.httpHeaderLock.Unlock()
}

func (downloader *ImageDownloader) GetValueForField(field string) (val string) {
	if downloader == nil || len(downloader.httpHeaders) == 0 {
		return
	}
	downloader.httpHeaderLock.Lock()
	val = downloader.httpHeaders[field]
	downloader.httpHeaderLock.Unlock()
	return
}

func (downloader *ImageDownloader) DownloadImage(
	url string,
	options ImageDownloaderOptions,
	ctx ImageContext,
	progressCb ImageDownloaderProgressBlock,
	completionCb ImageDownloaderCompletedBlock) *ImageDownloadToken {
	if downloader == nil {
		return nil
	}
	if len(url) == 0 {
		if completionCb != nil {
			completionCb(nil, ImageOperationInvalidURLError, true)
		}
		return nil
	}
	downloader.operationLock.Lock()
	operation, ok := downloader.urlOperations[url]
	var err error
	var token interface{}
	if !ok || operation.GetIsFinished() || operation.GetIsCanceled() { // 确保已经完成或者取消的任务移除
		operation, err = downloader.CreateDownloaderOperation(url, ctx, options)
		if err != nil {
			downloader.operationLock.Unlock()
			if completionCb != nil {
				completionCb(nil, err, true)
			}
			return nil
		}
		operation.SetCompletionFunc(func(op WebImageOperationProtocol) {
			downloader.operationLock.Lock()
			delete(downloader.urlOperations, url)
			downloader.operationLock.Unlock()
		})
		downloader.urlOperations[url] = operation
		token = operation.AddHandlersForProgressAndCompletion(progressCb, completionCb)
		downloader.downloadQueue.AddOperation(operation)
	} else {
		token = operation.AddHandlersForProgressAndCompletion(progressCb, completionCb)
		if !operation.GetIsRunning() {
			if BitsHas(BitsType(options), BitsType(ImageDownloaderHighPriority)) {
				operation.SetOperationPriority(WebImageOperationPriorityHigh)
			} else if BitsHas(BitsType(options), BitsType(ImageDownloaderLowPriority)) {
				operation.SetOperationPriority(WebImageOperationPriorityLow)
			} else {
				operation.SetOperationPriority(WebImageOperationPriorityNormal)
			}
		}
	}
	downloader.operationLock.Unlock()
	downloadToken := newImageDownloadTokenWithObserver(operation)
	downloadToken.downloadOperationToken = token
	return downloadToken
}

func (downloader *ImageDownloader) CancelAllDownloads() {
	if downloader == nil || downloader.downloadQueue == nil {
		return
	}
	downloader.downloadQueue.CancelAllOperations()
}

func (downloader *ImageDownloader) GetConfig() *ImageDownloaderConfig {
	if downloader == nil || downloader.config == nil {
		return nil
	}
	return downloader.config
}

var kSharedImageDownloader *ImageDownloader
var kSharedImageDownloaderToken sync.Once

func GetSharedImageDownloader() *ImageDownloader {
	if kSharedImageDownloader == nil {
		kSharedImageDownloaderToken.Do(func() {
			kSharedImageDownloader = NewImageDownloader()
		})
	}
	return kSharedImageDownloader
}

func NewImageDownloader() *ImageDownloader {
	return NewImageDownloaderWithConfig(GetDefaultDownloaderConfig())
}

func NewImageDownloaderWithConfig(config *ImageDownloaderConfig) *ImageDownloader {
	if config == nil {
		config = GetDefaultDownloaderConfig()
	}
	headers := map[string]string{
		"Accept":     "image/*,*/*;q=0.8",
		"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.89 Safari/537.36",
	}
	client := config.Client
	if client == nil {
		client = http.DefaultClient
	}
	return &ImageDownloader{
		config:        config.Copy(),
		httpHeaders:   headers,
		client:        client,
		downloadQueue: NewImageDownloadQueue(config.MaxConcurrentDownloads),
		urlOperations: map[string]ImageDownloadOperationProtocol{},
	}
}
