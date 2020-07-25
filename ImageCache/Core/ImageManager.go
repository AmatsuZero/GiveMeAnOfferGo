package Core

import (
	"net/url"
	"sync"
)

type ImageManagerDelegate interface {
	ShouldDownloadForImageManager(mgr *ImageManager, url *url.URL) bool
	ShouldBlockFailedURLForImageManager(mgr *ImageManager, url *url.URL, err error) bool
}

type InternalCompletionBlock func(data []byte, err error, cacheType ImageCacheType, finished bool, imageURL *url.URL)

type ImageCombineOperation struct {
	*WebImageOperation
	CacheOperation WebImageOperationProtocol
	mgr            *ImageManager
}

type ImageManager struct {
	imageCache            ImageCacheProtocol
	Delegate              ImageManagerDelegate
	failedURLs            map[*url.URL]bool
	failedURLsLock        sync.Mutex
	CacheKeyFilter        ImageCacheKeyFilterProtocol
	runningOperations     map[*ImageCombineOperation]bool
	runningOperationsLock sync.Mutex
}

func (mgr *ImageManager) GetImageCache() ImageCacheProtocol {
	if mgr == nil {
		return nil
	}
	return mgr.imageCache
}

func (mgr *ImageManager) CancelAll() {
	if mgr == nil {
		return
	}
	mgr.runningOperationsLock.Lock()
	defer mgr.runningOperationsLock.Unlock()
	for op := range mgr.runningOperations {
		op.Cancel()
	}
}

func (mgr *ImageManager) IsRunning() bool {
	if mgr == nil {
		return false
	}
	mgr.runningOperationsLock.Lock()
	defer mgr.runningOperationsLock.Unlock()
	return len(mgr.runningOperations) > 0
}

func (mgr *ImageManager) RemoveFailedRL(url *url.URL) {
	if mgr == nil || url == nil {
		return
	}
	mgr.failedURLsLock.Lock()
	defer mgr.failedURLsLock.Unlock()
	delete(mgr.failedURLs, url)
}

func (mgr *ImageManager) RemoveFailedURLs() {
	if mgr == nil {
		return
	}
	mgr.failedURLsLock.Lock()
	defer mgr.failedURLsLock.Unlock()
	mgr.failedURLs = map[*url.URL]bool{}
}

func (mgr *ImageManager) CacheKeyForURL(url *url.URL) string {
	if url == nil {
		return ""
	}
	if mgr != nil && mgr.CacheKeyFilter != nil {
		return mgr.CacheKeyFilter.CacheKeyForURL(url)
	}
	return url.String()
}

func (mgr *ImageManager) CacheKeyForURLAndContext(url *url.URL, ctx ImageContext) string {
	if mgr != nil && mgr.CacheKeyFilter != nil {
		return mgr.CacheKeyForURL(url)
	}
	filter, ok := ctx[kImageContextImageCache].(ImageCacheKeyFilterProtocol)
	if !ok {
		return ""
	}
	return filter.CacheKeyForURL(url)
}

func (mgr *ImageManager) callCacheProcessForOperation(
	op *ImageCombineOperation,
	url *url.URL,
	options ImageOptions,
	ctx ImageContext,
	completedBlock InternalCompletionBlock) {
	if mgr == nil {
		return
	}
	var imageCache ImageCacheProtocol
	if ic, ok := ctx[kImageContextImageCache]; ok {
		c, ok := ic.(ImageCacheProtocol)
		if ok {
			imageCache = c
		}
	}
	if imageCache == nil {
		imageCache = mgr.imageCache
	}
	originalQueryCacheType := ImageCacheTypeNone
	if t, ok := ctx[kbImageContextOriginalQueryCacheType]; ok {
		originalQueryCacheType = t.(ImageCacheType)
	}
	if originalQueryCacheType != ImageCacheTypeNone {
		key := mgr.CacheKeyForURLAndContext(url, ctx)
		op.CacheOperation = imageCache.QueryImageCacheWithCacheType(key,
			options, ctx, originalQueryCacheType,
			func(data []byte, cacheType ImageCacheType) {
				if op.GetIsCanceled() {

					return
				}
			})
	} else {
		mgr.callDownloadProcess(op, url, options, ctx, nil, originalQueryCacheType, completedBlock)
	}
}

func (mgr *ImageManager) callDownloadProcess(
	op *ImageCombineOperation,
	url *url.URL,
	options ImageOptions,
	ctx ImageContext,
	cachedData []byte,
	cacheType ImageCacheType,
	completedBlock InternalCompletionBlock) {
	if mgr == nil {
		return
	}
	shouldDownload := BitsHas(BitsType(options), BitsType(ImageFromCacheOnly))
	shouldDownload = shouldDownload && BitsHas(BitsType(options), BitsType(ImageRefreshCached))

}

func (mgr *ImageManager) LoadImage(url *url.URL, options ImageOptions, completedBlock InternalCompletionBlock) *ImageCombineOperation {
	if mgr == nil || completedBlock == nil {
		return nil
	}
	op := &ImageCombineOperation{
		WebImageOperation: NewWebImageOperation(nil),
		CacheOperation:    nil,
		mgr:               mgr,
	}
	isFailedURL := false
	if url != nil {
		mgr.failedURLsLock.Lock()
		_, isFailedURL = mgr.failedURLs[url]
		mgr.failedURLsLock.Unlock()
	}
	if (url != nil && len(url.String()) == 0) || BitsHas(BitsType(options), BitsType(ImageRetryFailed)) {
		err := ImageInvalidURLError
		if isFailedURL {
			err = ImageBlackListedError
		}
		mgr.callCompletionBlock(completedBlock, nil, err, ImageCacheTypeNone, true, url)
		return op
	}
	mgr.runningOperationsLock.Lock()
	mgr.runningOperations[op] = true
	mgr.runningOperationsLock.Unlock()

	return op
}

func (mgr *ImageManager) callCompletionBlock(
	completionBlock InternalCompletionBlock,
	data []byte,
	err error,
	cacheType ImageCacheType,
	finished bool,
	url *url.URL) {
	if completionBlock == nil {
		return
	}
	go completionBlock(data, err, cacheType, finished, url)
}

func NewImageManagerWithCache(cache ImageCacheProtocol) *ImageManager {
	return &ImageManager{
		imageCache: cache,
		failedURLs: map[*url.URL]bool{},
	}
}

var sharedImageManager *ImageManager
var sharedImageManagerToken sync.Once

func GetSharedImageManager() *ImageManager {
	if sharedImageManager == nil {
		sharedImageManagerToken.Do(func() {
			sharedImageManager = NewImageManagerWithCache(GetSharedImageCache())
		})
	}
	return sharedImageManager
}
