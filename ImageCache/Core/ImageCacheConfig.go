package Core

import (
	"gopkg.in/djherbis/times.v1"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sort"
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
	ShouldCacheImageInMemory             bool
	ShouldUseWeakMemoryCache             bool
	ShouldRemoveExpiredDataWhenAvailable bool
	DiskCacheReadingOption               DataReadingOption
	DiskCacheWritingOption               DataWritingOption
	MaxDiskAge                           float64
	MaxMemoryCount                       uint64
	DiskCacheExpireType                  ImageCacheConfigExpireType
	MemoryCache                          reflect.Type // 实现 MemoryCache 协议的类型
	DiskCache                            reflect.Type
}

const kDefaultCacheMaxDiskAge = 60 * 60 * 24 * 7 // 1 week

var defaultImageCache ImageCacheConfig

func init() {
	defaultImageCache = NewImageCacheConfig()
}

func DefaultImageCacheConfig() ImageCacheConfig {
	return defaultImageCache
}

func NewImageCacheConfig() ImageCacheConfig {
	return ImageCacheConfig{
		ShouldDisableRemoteBackup:            true,
		ShouldCacheImageInMemory:             true,
		ShouldUseWeakMemoryCache:             true,
		ShouldRemoveExpiredDataWhenAvailable: false,
		DiskCacheReadingOption:               0,
		DiskCacheWritingOption:               WritingAtomic,
		MaxDiskAge:                           kDefaultCacheMaxDiskAge,
		MaxMemoryCount:                       0, // 表明没有限制
		DiskCacheExpireType:                  ImageCacheConfigExpireTypeModificationDate,
		MemoryCache:                          reflect.TypeOf((*memoryCache)(nil)).Elem(),
		DiskCache:                            reflect.TypeOf((*diskCache)(nil)).Elem(),
	}
}

func (config ImageCacheConfig) Copy() ImageCacheConfig {
	return ImageCacheConfig{
		ShouldDisableRemoteBackup:            config.ShouldDisableRemoteBackup,
		ShouldCacheImageInMemory:             config.ShouldCacheImageInMemory,
		ShouldUseWeakMemoryCache:             config.ShouldUseWeakMemoryCache,
		ShouldRemoveExpiredDataWhenAvailable: config.ShouldRemoveExpiredDataWhenAvailable,
		DiskCacheReadingOption:               config.DiskCacheReadingOption,
		DiskCacheWritingOption:               config.DiskCacheWritingOption,
		MaxDiskAge:                           config.MaxDiskAge,
		DiskCacheExpireType:                  config.DiskCacheExpireType,
		MemoryCache:                          config.MemoryCache,
		DiskCache:                            config.DiskCache,
	}
}
