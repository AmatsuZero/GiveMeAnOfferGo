package Core

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const DiskCacheExtendedAttributeName = "com.daubert.imageCache"

type DishCache interface {
	GetCachePath() string
	GetCacheConfig() *ImageCacheConfig
	ContainDataForKey(key string) bool
	GetDataForKey(key string) []byte
	SetDataForKey(key string, data []byte)
	GetExtendedDataForKey(key string)
	SetExtendedDataForKey(key string, data []byte)
	RemoveDataForKey(key string)
	RemoveAllData()
	RemoveExpiredData()
	GetCachePathForKey(key string)
	/// 缓存目录下的缓存数量
	GetTotalCount() uint64
	/// 缓存目录下缓存大小
	GetTotalSize() int64
}

type diskCache struct {
	Config        *ImageCacheConfig
	diskCachePath string
}

func NewDiskCache(config *ImageCacheConfig, path string) *diskCache {
	return &diskCache{
		Config:        config,
		diskCachePath: path,
	}
}

func (cache *diskCache) CachePathForKey(key string) string {
	return cache.cachePathForKeyInPath(key, cache.diskCachePath)
}

func (cache *diskCache) cachePathForKeyInPath(key, path string) string {
	if len(key) == 0 {
		return ""
	}
	ext := filepath.Ext(key)
	has := md5.Sum([]byte(key))
	key = fmt.Sprintf("%x", has) //将[]byte转成16进制
	return filepath.Join(path, key, fmt.Sprintf(".%v", ext))
}

func (cache *diskCache) MoveCacheDirectory(from, to string) {

}

func (cache *diskCache) GetTotalCount() uint64 {
	files, _ := ioutil.ReadDir(cache.diskCachePath)
	return uint64(len(files))
}

func (cache *diskCache) GetTotalSize() (size int64) {
	files, _ := ioutil.ReadDir(cache.diskCachePath)
	for _, file := range files {
		size += file.Size()
	}
	return
}

func (cache *diskCache) SetDataForKey(key string, data []byte) {
	if len(key) == 0 || data == nil || len(cache.diskCachePath) == 0 {
		return
	}
	// 检查目录是否存在, 不存在，则创建
	if _, err := os.Stat(cache.diskCachePath); os.IsNotExist(err) {
		err = os.MkdirAll(cache.diskCachePath, os.ModeDir)
		if err != nil {
			return
		}
	}
	path := cache.CachePathForKey(key)
	_ = ioutil.WriteFile(path, data, os.ModePerm)
	if cache.Config.ShouldDisableRemoteBackup { // 暂无对应处理

	}
}

func (cache *diskCache) ContainDataForKey(key string) bool {
	path := cache.CachePathForKey(key)
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func (cache *diskCache) GetDataForKey(key string) []byte {
	path := cache.CachePathForKey(key)
	data, _ := ioutil.ReadFile(path)
	return data
}

func (cache *diskCache) SetExtendedDataForKey(key string, data []byte) {

}

func (cache *diskCache) RemoveAllData() {
	if len(cache.diskCachePath) == 0 {
		return
	}
	_ = os.RemoveAll(cache.diskCachePath)
	_ = os.MkdirAll(cache.diskCachePath, os.ModeDir)
}
