package Core

import (
	"context"
	"github.com/reactivex/rxgo/v2"
	"os"
	"path/filepath"
	"sync"
)

func BitsHas(options, flag BitsType) bool {
	return options&flag != 0
}

type ImageCacheType int

const (
	ImageCacheTypeNone ImageCacheType = iota
	ImageCacheTypeDisk
	ImageCacheTypeMemory
	ImageCacheTypeAll
)

type ImageCacheCheckCompletionBlock func(isInCache bool)
type ImageCacheQueryDataCompletionBlock func(data []byte)
type ImageCacheCalculateSizeBlock func(fileCount, totalSize uint64)
type ImageCacheAdditionalCachePathBlock func(key string) string
type ImageCacheQueryCompletionBlock func(data []byte, cacheType ImageCacheType)
type ImageCacheContainsCompletionBlock func(containsCacheType ImageCacheType)
type ImageNoParamsBlock func()

func ImageCacheDecodeImageData(imageData []byte, cacheKey string, options ImageCacheOptions, context WebImageContext) {

}

type ImageCacheOptions BitsType

const (
	ImageCacheQueryMemoryData ImageCacheOptions = 1 << iota
	ImageCacheQueryMemoryDataSync
	ImageCacheQueryDiskDataSync
	ImageCacheScaleDownLargeImages
	ImageCacheAvoidDecodeImage
	ImageCacheDecodeFirstFrameOnly
	ImageCachePreloadAllFrames
	ImageCacheMatchAnimatedImageClass
)

type ImageCacheProtocol interface {
	QueryImage(key string, ops ImageCacheOptions, ctx WebImageContext, cb ImageCacheQueryCompletionBlock) WebImageOperation
	QueryImageWithCacheType(key string, ops ImageCacheOptions,
		ctx WebImageContext, ct ImageCacheType, cb ImageCacheQueryCompletionBlock) WebImageOperation
	StoreImage(data []byte, key string, ct ImageCacheType, cb ImageNoParamsBlock)
	RemoveImageForKey(key string, cacheType ImageCacheType, cb ImageNoParamsBlock)
	ContainsImageForKey(key string, cacheType ImageCacheType, cb ImageCacheContainsCompletionBlock)
	ClearWithCacheType(cacheType ImageCacheType, cb ImageNoParamsBlock)
}

type ImageCache struct {
	config                   *ImageCacheConfig
	memoryCache              MemoryCache
	diskCache                DishCache
	diskCachePath            string
	AdditionalCachePathBlock ImageCacheAdditionalCachePathBlock
}

func (cache *ImageCache) QueryImage(key string, ops ImageCacheOptions, ctx WebImageContext, cb ImageCacheQueryCompletionBlock) WebImageOperation {
	panic("implement me")
}

func (cache *ImageCache) QueryImageWithCacheType(key string, ops ImageCacheOptions,
	ctx WebImageContext, ct ImageCacheType, cb ImageCacheQueryCompletionBlock) WebImageOperation {
	panic("implement me")
}

func (cache *ImageCache) ImageDataFromMemoryForKey(key string) []byte {
	if cache == nil || cache.memoryCache == nil {
		return nil
	}
	return cache.memoryCache.ObjectForKey(key).([]byte)
}

func (cache *ImageCache) StoreImage(data []byte, key string, ct ImageCacheType, cb ImageNoParamsBlock) {
	defer func() {
		if cb != nil {
			cb()
		}
	}()
	if len(data) == 0 || ct == ImageCacheTypeNone || cache == nil {
		return
	}
	if ct == ImageCacheTypeMemory || ct == ImageCacheTypeAll {
		if cache.config.ShouldCacheImageInMemory {
			cache.StoreImageToMemory(data, key)
		}
	}
	if ct == ImageCacheTypeDisk || ct == ImageCacheTypeAll {
		cache.StoreImageToDisk(data, key)
	}
}

func (cache *ImageCache) StoreImageToDisk(data []byte, key string) {
	disposed := cache.storeImageToDisk(data, key).Run()
	<-disposed
}

func (cache *ImageCache) StoreImageToMemory(data []byte, key string) {
	if cache.memoryCache == nil {
		return
	}
	cache.memoryCache.SetObjectWithKey(key, data)
}

func (cache *ImageCache) storeImageToDisk(data []byte, key string) rxgo.Observable {
	if cache == nil || cache.diskCache == nil {
		return rxgo.Empty()
	}
	return rxgo.Create([]rxgo.Producer{func(ctx context.Context, next chan<- rxgo.Item) {
		next <- rxgo.Error(cache.diskCache.SetDataForKey(key, data))
	}})
}

func (cache *ImageCache) RemoveImageForKey(key string, cacheType ImageCacheType, cb ImageNoParamsBlock) {
	panic("implement me")
}

func (cache *ImageCache) ContainsImageForKey(key string, cacheType ImageCacheType, cb ImageCacheContainsCompletionBlock) {
	panic("implement me")
}

func (cache *ImageCache) DiskImageExistsWithKey(key string) bool {
	item, err := rxgo.Create([]rxgo.Producer{func(ctx context.Context, next chan<- rxgo.Item) {
		next <- rxgo.Of(cache.diskImageExistsWithKey(key))
	}}).First().Get()
	if err != nil {
		return false
	}
	return item.V.(bool)
}

func (cache *ImageCache) diskImageExistsWithKey(key string) bool {
	if cache == nil || cache.diskCache == nil {
		return false
	}
	return cache.diskCache.ContainDataForKey(key)
}

func (cache *ImageCache) DiskImageExistsWithKeyAsync(key string, cb ImageCacheCheckCompletionBlock) {
	rxgo.Create([]rxgo.Producer{func(ctx context.Context, next chan<- rxgo.Item) {
		next <- rxgo.Of(cache.diskImageExistsWithKey(key))
	}}).DoOnNext(func(i interface{}) {
		if cb != nil {
			cb(i.(bool))
		}
	})
}

func (cache *ImageCache) ClearWithCacheType(cacheType ImageCacheType, cb ImageNoParamsBlock) {
	panic("implement me")
}

func (cache *ImageCache) GetImageCacheConfig() *ImageCacheConfig {
	return cache.config
}

func (cache *ImageCache) GetMemoryCache() MemoryCache {
	return cache.memoryCache
}

func (cache *ImageCache) GetDiskCache() DishCache {
	return cache.diskCache
}

func (cache *ImageCache) GetDiskCachePath() string {
	return cache.diskCachePath
}

func (cache *ImageCache) onMemoryCacheEjection(key interface{}, value interface{}) {

}

func (cache *ImageCache) GetCachePathForKey(key string) string {
	if cache == nil || cache.diskCache == nil {
		return ""
	}
	return cache.diskCache.GetCachePathForKey(key)
}

var sharedImageCache *ImageCache
var defaultDiskUserCachePath string

func init() {
	path, err := os.UserCacheDir()
	if err != nil {
		path = os.TempDir()
	}
	defaultDiskUserCachePath = filepath.Join(path, "com.daubert.imageCache")
}

func NewImageCache() *ImageCache {
	return NewImageCacheWithNameSpace("com.daubert.imageCache")
}

func NewImageCacheWithNameSpace(ns string) *ImageCache {
	return NewImageCacheWithNameSpaceInDir(ns, defaultDiskUserCachePath)
}

func NewImageCacheWithNameSpaceInDir(ns, dir string) *ImageCache {
	return NewImageCacheWithNameSpaceInDirAndConfig(ns, dir, GetDefaultImageCacheConfig())
}

func NewImageCacheWithNameSpaceInDirAndConfig(ns, dir string, config *ImageCacheConfig) *ImageCache {
	cache := &ImageCache{}
	if config == nil {
		config = GetDefaultImageCacheConfig()
	}
	cache.config = config.Copy()
	if len(ns) == 0 {
		ns = "default"
	}
	if len(dir) > 0 {
		dir = filepath.Join(dir, ns)
	} else {
		dir = filepath.Join(defaultDiskUserCachePath, ns)
	}
	cache.diskCachePath = dir
	if config.DiskCacheCreationFn != nil {
		cache.diskCache = config.DiskCacheCreationFn(config, dir)
	}
	if config.MemoryCacheCreationFn != nil {
		mc, _ := config.MemoryCacheCreationFn(config, cache.onMemoryCacheEjection)
		cache.memoryCache = mc
	}
	return cache
}

var imageCacheOnce sync.Once

func GetSharedImageCache() *ImageCache {
	imageCacheOnce.Do(func() {
		if defaultImageCacheConfig == nil {
			sharedImageCache = NewImageCache()
		}
	})
	return sharedImageCache
}
