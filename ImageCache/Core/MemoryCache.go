package Core

type MemoryCache interface {
	GetImageCacheConfig() ImageCacheConfig
	ObjectForKey() interface{}
	SetObjectWithKey(key, value interface{})
}
