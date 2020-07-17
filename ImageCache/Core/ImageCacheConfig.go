package Core

import (
	"github.com/hashicorp/golang-lru/simplelru"
	"gopkg.in/djherbis/times.v1"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

type ImageCacheConfigExpireType int

const (
	// When the image cache is accessed it will update this value
	ImageCacheConfigExpireTypeAccessDate ImageCacheConfigExpireType = iota
	// When the image cache is created or modified it will update this value (Default)
	ImageCacheConfigExpireTypeModificationDate
	// When the image cache is created it will update this value
	ImageCacheConfigExpireTypeCreationDate
	// When the image cache is created, modified, renamed, file attribute updated (like permission, xattr)  it will update this value
	ImageCacheConfigExpireTypeChangeDate
)

func (expireType ImageCacheConfigExpireType) AccessFilesInDir(dir string) (files []os.FileInfo, err error) {
	files, err = ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	// 按照类型，排序
	sort.Slice(files, func(i, j int) bool {
		lhs := filepath.Join(dir, files[i].Name())
		rhs := filepath.Join(dir, files[j].Name())
		return expireType.getReferenceDate(lhs).Second() > expireType.getReferenceDate(rhs).Second()
	})
	return
}

func (expireType ImageCacheConfigExpireType) getReferenceDate(path string) time.Time {
	switch expireType {
	case ImageCacheConfigExpireTypeAccessDate:
		t, _ := times.Stat(path)
		return t.AccessTime()
	case ImageCacheConfigExpireTypeCreationDate:
		t, _ := times.Stat(path)
		return t.BirthTime()
	case ImageCacheConfigExpireTypeChangeDate:
		t, _ := times.Stat(path)
		return t.ChangeTime()
	default:
		t, _ := times.Stat(path)
		return t.ChangeTime()
	}
}

type ImageCacheConfig struct {
	ShouldDisableRemoteBackup            bool // 远端同步，对应iCloud
	ShouldCacheImageInMemory             bool // 是否启用内存缓存
	ShouldUseWeakMemoryCache             bool
	ShouldRemoveExpiredDataWhenAvailable bool
	DiskCacheReadingOption               DataReadingOption
	DiskCacheWritingOption               DataWritingOption
	MaxDiskAge                           float64 // 磁盘缓存时间，默认一周
	MaxDiskSize                          int64   // 磁盘缓存大小设置，默认为0，表示不限制
	MaxMemoryCount                       uint64
	DiskCacheExpireType                  ImageCacheConfigExpireType
	MemoryCacheCreationFn                func(config *ImageCacheConfig, cb simplelru.EvictCallback) (MemoryCache, error)
	DiskCacheCreationFn                  func(config *ImageCacheConfig, path string) DishCache
}

const kDefaultCacheMaxDiskAge = 60 * 60 * 24 * 7 // 1 week

var defaultImageCacheConfig *ImageCacheConfig
var imageCacheConfigOnce sync.Once

func GetDefaultImageCacheConfig() *ImageCacheConfig {
	imageCacheConfigOnce.Do(func() {
		if defaultImageCacheConfig != nil {
			defaultImageCacheConfig = NewImageCacheConfig()
		}
	})
	return defaultImageCacheConfig
}

func NewImageCacheConfig() *ImageCacheConfig {
	config := &ImageCacheConfig{
		ShouldDisableRemoteBackup:            true,
		ShouldCacheImageInMemory:             true,
		ShouldUseWeakMemoryCache:             true,
		ShouldRemoveExpiredDataWhenAvailable: false,
		DiskCacheReadingOption:               0,
		DiskCacheWritingOption:               WritingAtomic,
		MaxDiskAge:                           kDefaultCacheMaxDiskAge,
		MaxMemoryCount:                       0, // 表明没有限制
		MaxDiskSize:                          0, // 表示没有限制
		DiskCacheExpireType:                  ImageCacheConfigExpireTypeModificationDate,
	}
	config.MemoryCacheCreationFn = func(config *ImageCacheConfig, cb simplelru.EvictCallback) (MemoryCache, error) {
		return newMemoryCacheWithConfig(config, cb)
	}
	config.DiskCacheCreationFn = func(config *ImageCacheConfig, path string) DishCache {
		return newDiskCache(config, path)
	}
	return config
}

func (config *ImageCacheConfig) Copy() *ImageCacheConfig {
	if config == nil {
		return nil
	}
	return &ImageCacheConfig{
		ShouldDisableRemoteBackup:            config.ShouldDisableRemoteBackup,
		ShouldCacheImageInMemory:             config.ShouldCacheImageInMemory,
		ShouldUseWeakMemoryCache:             config.ShouldUseWeakMemoryCache,
		ShouldRemoveExpiredDataWhenAvailable: config.ShouldRemoveExpiredDataWhenAvailable,
		DiskCacheReadingOption:               config.DiskCacheReadingOption,
		DiskCacheWritingOption:               config.DiskCacheWritingOption,
		MaxDiskAge:                           config.MaxDiskAge,
		MaxDiskSize:                          config.MaxDiskSize,
		DiskCacheExpireType:                  config.DiskCacheExpireType,
		MemoryCacheCreationFn:                config.MemoryCacheCreationFn,
		DiskCacheCreationFn:                  config.DiskCacheCreationFn,
	}
}
