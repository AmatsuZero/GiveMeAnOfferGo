package Core

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const kDiskCacheExtendedAttributeName = "com.daubert.imageCache"

type DishCache interface {
	GetCachePath() string
	GetCacheConfig() *ImageCacheConfig
	ContainDataForKey(key string) bool
	GetDataForKey(key string) []byte
	SetDataForKey(key string, data []byte)
	GetExtendedDataForKey(key string) []byte
	SetExtendedDataForKey(key string, data []byte)
	RemoveDataForKey(key string)
	RemoveAllData()
	RemoveExpiredData()
	GetCachePathForKey(key string) string
	/// 缓存目录下的缓存数量
	GetTotalCount() uint64
	/// 缓存目录下缓存大小
	GetTotalSize() int64
}

type diskCache struct {
	Config        *ImageCacheConfig
	diskCachePath string
}

func newDiskCache(config *ImageCacheConfig, path string) *diskCache {
	return &diskCache{
		Config:        config,
		diskCachePath: path,
	}
}

func (cache *diskCache) GetCacheConfig() *ImageCacheConfig {
	return cache.Config
}

func (cache *diskCache) GetCachePath() string {
	return cache.diskCachePath
}

func (cache *diskCache) GetCachePathForKey(key string) string {
	return cache.cachePathForKeyInPath(key, cache.GetCachePath())
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
	if len(from) == 0 ||
		len(to) == 0 ||
		from == to {
		return
	}
	// 检查原来的是否是文件夹
	fileInfo, err := os.Stat(from)
	if err != nil || !fileInfo.IsDir() {
		return
	}
	// 检查新路径是否是文件夹
	fileInfo, err = os.Stat(to)
	if err != nil && os.IsNotExist(err) {
		if err = os.MkdirAll(to, os.ModePerm); err != nil {
			return
		}
	} else if !fileInfo.IsDir() { // 说明是文件，删除并重新创建
		if err = os.Remove(to); err != nil {
			return
		}
		if err = os.MkdirAll(to, os.ModePerm); err != nil {
			return
		}
	}
	files, err := ioutil.ReadDir(from)
	for _, file := range files {
		_ = os.Rename(filepath.Join(from, file.Name()), filepath.Join(to, file.Name()))
	}
	// 删除原来的文件夹
	_ = os.RemoveAll(from)
}

func (cache *diskCache) GetTotalCount() uint64 {
	files, _ := ioutil.ReadDir(cache.GetCachePath())
	return uint64(len(files))
}

func (cache *diskCache) GetTotalSize() (size int64) {
	files, _ := ioutil.ReadDir(cache.GetCachePath())
	for _, file := range files {
		size += file.Size()
	}
	return
}

func (cache *diskCache) SetDataForKey(key string, data []byte) {
	if len(key) == 0 || data == nil || len(cache.GetCachePath()) == 0 {
		return
	}
	// 检查目录是否存在, 不存在，则创建
	if _, err := os.Stat(cache.GetCachePath()); os.IsNotExist(err) {
		err = os.MkdirAll(cache.GetCachePath(), os.ModeDir)
		if err != nil {
			return
		}
	}
	path := cache.GetCachePathForKey(key)
	_ = ioutil.WriteFile(path, data, os.ModePerm)
	if cache.Config.ShouldDisableRemoteBackup { // 暂无对应处理

	}
}

func (cache *diskCache) ContainDataForKey(key string) bool {
	path := cache.GetCachePathForKey(key)
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func (cache *diskCache) GetDataForKey(key string) []byte {
	path := cache.GetCachePathForKey(key)
	data, _ := ioutil.ReadFile(path)
	return data
}

func (cache *diskCache) GetExtendedDataForKey(key string) []byte {
	if len(key) == 0 {
		return nil
	}
	path := cache.GetCachePathForKey(key)
	data, _ := getExtendedAttribute(kDiskCacheExtendedAttributeName, path, false)
	return data
}

func (cache *diskCache) SetExtendedDataForKey(key string, data []byte) {
	if len(key) == 0 {
		return
	}
	path := cache.GetCachePathForKey(key)
	if len(data) == 0 { // Remove
		_ = removeExtendedAttribute(kDiskCacheExtendedAttributeName, path, false)
	} else { // Override
		_ = setExtendedAttribute(kDiskCacheExtendedAttributeName, path, data, false, true)
	}
}

func (cache *diskCache) RemoveDataForKey(key string) {
	if len(key) == 0 {
		return
	}
	path := cache.GetCachePathForKey(key)
	_ = os.Remove(path)
}

func (cache *diskCache) RemoveAllData() {
	if len(cache.GetCachePath()) == 0 {
		return
	}
	_ = os.RemoveAll(cache.GetCachePath())
	_ = os.MkdirAll(cache.GetCachePath(), os.ModeDir)
}

func (cache *diskCache) RemoveExpiredData() {

}
