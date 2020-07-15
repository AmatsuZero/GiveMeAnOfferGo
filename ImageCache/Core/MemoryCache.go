package Core

import (
	"github.com/hashicorp/golang-lru/simplelru"
	"sync"
)

type MemoryCache interface {
	GetImageCacheConfig() *ImageCacheConfig
	ObjectForKey(key interface{}) interface{}
	SetObjectWithKey(key, value interface{})
	RemoveObjectForKey(key interface{})
	RemoveAllObjects()
}

type memoryCache struct {
	*simplelru.LRU
	Config         *ImageCacheConfig
	mutex          sync.Mutex
	weakMap        *WeakMap
	totalCostLimit int
}

func newMemoryCacheWithConfig(config *ImageCacheConfig, cb simplelru.EvictCallback) (*memoryCache, error) {
	count := config.MaxMemoryCount
	if count == 0 {
		count = 1
	}
	lru, err := simplelru.NewLRU(int(count), cb)
	if err != nil {
		return nil, err
	}
	return &memoryCache{
		LRU:            lru,
		Config:         config,
		weakMap:        NewWeakMap(WeakMapStrongToWeak),
		totalCostLimit: int(count),
	}, nil
}

func (cache *memoryCache) GetImageCacheConfig() *ImageCacheConfig {
	return cache.Config
}

func (cache *memoryCache) ObjectForKey(key interface{}) interface{} {
	val, ok := cache.Get(key)
	if !cache.Config.ShouldUseWeakMemoryCache {
		return val
	}
	if key != nil && ok {
		cache.mutex.Lock()
		obj, ok := cache.weakMap.Get(key)
		cache.mutex.Unlock()
		if ok {
			cache.SetObjectWithKey(key, obj)
		}
	}
	return nil
}

func (cache *memoryCache) SetObjectWithKey(key, value interface{}) {
	if cache.Config.MaxMemoryCount == 0 && cache.Len() == cache.totalCostLimit { // 没有限制，重新分配空间
		cache.totalCostLimit *= 2
		cache.Resize(cache.totalCostLimit)
	}
	cache.Add(key, value)
}

func (cache *memoryCache) RemoveObjectForKey(key interface{}) {
	cache.Remove(key)
	if !cache.Config.ShouldUseWeakMemoryCache || key == nil {
		return
	}
	cache.mutex.Lock()
	cache.weakMap.Remove(key)
	cache.mutex.Unlock()
}

func (cache *memoryCache) RemoveAllObjects() {
	cache.Purge()
	if !cache.Config.ShouldUseWeakMemoryCache {
		return
	}
	cache.mutex.Lock()
	cache.weakMap.RemoveAll()
	cache.mutex.Unlock()
}
