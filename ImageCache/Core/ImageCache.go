package Core

import (
	"context"
	"github.com/reactivex/rxgo/v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

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
	QueryImageCache(key string, ops ImageCacheOptions, ctx ImageContext, cb ImageCacheQueryCompletionBlock) WebImageOperationProtocol
	QueryImageCacheWithCacheType(key string, ops ImageCacheOptions,
		ctx ImageContext, ct ImageCacheType, cb ImageCacheQueryCompletionBlock) WebImageOperationProtocol
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

func (cache *ImageCache) QueryImageCache(key string, ops ImageCacheOptions, ctx ImageContext, cb ImageCacheQueryCompletionBlock) WebImageOperationProtocol {
	panic("implement me")
}

func (cache *ImageCache) QueryImageCacheWithCacheType(key string, ops ImageCacheOptions,
	ctx ImageContext, ct ImageCacheType, cb ImageCacheQueryCompletionBlock) WebImageOperationProtocol {
	if cache == nil || len(key) == 0 || ct == ImageCacheTypeNone {
		if cb != nil {
			cb(nil, ImageCacheTypeNone)
		}
		return nil
	}
	memoryData := cache.memoryCache.GetObjectForKey(key).([]byte)
	if ct == ImageCacheTypeMemory && BitsHas(BitsType(ops), BitsType(ImageCacheQueryMemoryData)) {
		if cb != nil {
			cb(memoryData, ImageCacheTypeMemory)
		}
		return nil
	}
	type query struct {
		t ImageCacheType
		d []byte
	}
	shouldQuerySync := (memoryData != nil && BitsHas(BitsType(ops), BitsType(ImageCacheQueryMemoryDataSync))) ||
		(memoryData == nil && BitsHas(BitsType(ops), BitsType(ImageCacheQueryDiskDataSync)))
	ob := rxgo.Create([]rxgo.Producer{func(ctx context.Context, next chan<- rxgo.Item) {
		data := query{}
		if memoryData != nil {
			data.t = ImageCacheTypeMemory
			data.d = memoryData
		} else {
			diskData, err := cache.getDiskImageDataBySearchingAllPathsForKey(key)
			if err == nil && cache.memoryCache != nil && cache.config.ShouldCacheImageInMemory {
				cache.StoreImageDataToMemory(diskData, key)
			}
			data.t = ImageCacheTypeDisk
			data.d = diskData
		}
		next <- rxgo.Item{V: data}
	}})
	ob.DoOnNext(func(i interface{}) {
		if cb == nil {
			return
		}
		q := i.(query)
		if shouldQuerySync {
			cb(q.d, q.t)
		} else {
			go cb(q.d, q.t)
		}
	})
	return newWebImageOperation(ob)
}

func (cache *ImageCache) ImageDataFromMemoryForKey(key string) []byte {
	if cache == nil || cache.memoryCache == nil {
		return nil
	}
	return cache.memoryCache.GetObjectForKey(key).([]byte)
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
			cache.StoreImageDataToMemory(data, key)
		}
	}
	if ct == ImageCacheTypeDisk || ct == ImageCacheTypeAll {
		cache.StoreImageDataToDisk(data, key)
	}
}

func (cache *ImageCache) StoreImageDataToDisk(data []byte, key string) {
	<-cache.storeImageToDisk(data, key).Run()
}

func (cache *ImageCache) StoreImageDataToMemory(data []byte, key string) {
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
	if cache == nil {
		if cb != nil {
			cb()
		}
		return
	}
	switch cacheType {
	case ImageCacheTypeMemory:
		cache.removeCacheDataFromMemory(key).DoOnCompleted(rxgo.CompletedFunc(cb))
	case ImageCacheTypeDisk:
		cache.removeCacheDataFromDisk(key).DoOnCompleted(rxgo.CompletedFunc(cb))
	case ImageCacheTypeAll:
		rxgo.Concat([]rxgo.Observable{
			cache.removeCacheDataFromMemory(key),
			cache.removeCacheDataFromDisk(key),
		}).DoOnCompleted(rxgo.CompletedFunc(cb))
	case ImageCacheTypeNone:
		if cb != nil {
			cb()
		}
	}
}

func (cache *ImageCache) RemoveCacheDataFromMemory(key string) {
	<-cache.removeCacheDataFromMemory(key).Run()
}

func (cache *ImageCache) removeCacheDataFromMemory(key string) rxgo.Observable {
	if cache == nil || cache.memoryCache == nil || len(key) == 0 {
		return rxgo.Empty()
	}
	return rxgo.Create([]rxgo.Producer{func(ctx context.Context, next chan<- rxgo.Item) {
		next <- rxgo.Of(cache.memoryCache.RemoveObjectForKey(key))
	}})
}

func (cache *ImageCache) RemoveCacheDataFromDisk(key string) {
	<-cache.removeCacheDataFromDisk(key).Run()
}

func (cache *ImageCache) removeCacheDataFromDisk(key string) rxgo.Observable {
	if cache == nil || cache.diskCache == nil || len(key) == 0 {
		return rxgo.Empty()
	}
	return rxgo.Create([]rxgo.Producer{func(ctx context.Context, next chan<- rxgo.Item) {
		next <- rxgo.Of(cache.diskCache.RemoveDataForKey(key) == nil)
	}})
}

func (cache *ImageCache) ContainsImageForKey(key string, cacheType ImageCacheType, cb ImageCacheContainsCompletionBlock) {
	if cb == nil {
		return
	}
	if cache == nil {
		cb(ImageCacheTypeNone)
		return
	}
	switch cacheType {
	case ImageCacheTypeNone:
		cb(ImageCacheTypeNone)
	case ImageCacheTypeMemory:
		isInMemory := cache.ImageDataFromMemoryForKey(key) != nil
		if isInMemory {
			cb(ImageCacheTypeMemory)
		} else {
			cb(ImageCacheTypeNone)
		}
	case ImageCacheTypeDisk:
		// 检查磁盘缓存是否有
		cache.DiskImageExistsWithKeyAsync(key, func(isInCache bool) {
			if isInCache {
				cb(ImageCacheTypeDisk)
			} else {
				cb(ImageCacheTypeNone)
			}
		})
	case ImageCacheTypeAll:
		cache.ContainsImageForKey(key, ImageCacheTypeMemory, func(containsCacheType ImageCacheType) {
			if containsCacheType == ImageCacheTypeMemory {
				cb(containsCacheType)
			} else {
				cache.ContainsImageForKey(key, ImageCacheTypeDisk, cb)
			}
		})
	default:
		cb(ImageCacheTypeNone)
	}
}

func (cache *ImageCache) DiskImageExistsWithKey(key string) bool {
	item, err := rxgo.JustItem(cache.diskImageExistsWithKey(key)).Get()
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

func (cache *ImageCache) CalculateSizeWithCompletionBlock(cb ImageCacheCalculateSizeBlock) {
	count := uint64(0)
	size := uint64(0)
	if cache == nil || cache.diskCache == nil {
		if cb != nil {
			cb(count, size)
		}
		return
	}
	rxgo.Create([]rxgo.Producer{func(ctx context.Context, next chan<- rxgo.Item) {
		count = cache.diskCache.GetTotalCount()
		next <- rxgo.Of(count)
		size = uint64(cache.diskCache.GetTotalSize())
		next <- rxgo.Of(size)
	}}).DoOnCompleted(func() {
		if cb != nil {
			cb(count, size)
		}
	})
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
	if cache != nil {
		if cb != nil {
			cb()
		}
		return
	}
	switch cacheType {
	case ImageCacheTypeMemory:
		cache.clearMemoryCache().DoOnCompleted(rxgo.CompletedFunc(cb))
	case ImageCacheTypeDisk:
		cache.clearDiskCache().DoOnCompleted(rxgo.CompletedFunc(cb))
	case ImageCacheTypeAll:
		rxgo.Concat([]rxgo.Observable{
			cache.clearMemoryCache(),
			cache.clearDiskCache(),
		}).DoOnCompleted(rxgo.CompletedFunc(cb))
	default:
		if cb != nil {
			cb()
		}
	}
}

func (cache *ImageCache) ClearMemoryCache() {
	<-cache.clearMemoryCache().Run()
}

func (cache *ImageCache) clearMemoryCache() rxgo.Observable {
	if cache == nil || cache.memoryCache == nil {
		return rxgo.Empty()
	}
	return rxgo.Create([]rxgo.Producer{func(ctx context.Context, next chan<- rxgo.Item) {
		cache.memoryCache.RemoveAllObjects()
		next <- rxgo.Of(true)
	}})
}

func (cache *ImageCache) ClearDiskCache(cb ImageNoParamsBlock) {
	ob := cache.clearDiskCache()
	ob.DoOnCompleted(rxgo.CompletedFunc(cb))
}

func (cache *ImageCache) clearDiskCache() rxgo.Observable {
	if cache == nil || cache.diskCache == nil {
		return rxgo.Empty()
	}
	return rxgo.Create([]rxgo.Producer{func(ctx context.Context, next chan<- rxgo.Item) {
		cache.diskCache.RemoveAllData()
		next <- rxgo.Of(true)
	}})
}

func (cache *ImageCache) DeleteOldFilesWithCompletionBlock(cb ImageNoParamsBlock) {
	rxgo.Create([]rxgo.Producer{func(ctx context.Context, next chan<- rxgo.Item) {
		if cache == nil || cache.diskCache == nil {
			next <- rxgo.Of(true)
		} else {
			cache.diskCache.RemoveExpiredData()
			next <- rxgo.Of(true)
		}
	}}).DoOnCompleted(rxgo.CompletedFunc(cb))
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

func (cache *ImageCache) GetDiskImageData(key string) ([]byte, error) {
	item, err := rxgo.Create([]rxgo.Producer{func(ctx context.Context, next chan<- rxgo.Item) {
		data, err := cache.getDiskImageDataBySearchingAllPathsForKey(key)
		next <- rxgo.Item{
			V: data,
			E: err,
		}
	}}).First().Get()
	if err != nil {
		return nil, err
	}
	return item.V.([]byte), nil
}

func (cache *ImageCache) GetDiskCachePath() string {
	return cache.diskCachePath
}

func (cache *ImageCache) getDiskImageDataBySearchingAllPathsForKey(key string) ([]byte, error) {
	if cache == nil || cache.diskCache == nil || len(key) == 0 {
		return nil, InvalidParamError
	}
	data, err := cache.diskCache.GetDataForKey(key)
	if err != nil {
		return nil, err
	}
	if cache.AdditionalCachePathBlock != nil {
		path := cache.AdditionalCachePathBlock(key)
		data, err = ioutil.ReadFile(path)
	}
	return data, err
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
