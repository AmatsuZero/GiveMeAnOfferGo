package Core

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

type ImageCacheConfig struct {
	ShouldDisableRemoteBackup            bool
	ShouldCacheImageInMemory             bool
	ShouldUseWeakMemoryCache             bool
	ShouldRemoveExpiredDataWhenAvailable bool
	DiskCacheReadingOption               DataReadingOption
	DiskCacheWritingOption               DataWritingOption
	MaxDiskAge                           float64
	MaxMemoryCost                        uint64
	DiskCacheExpireType                  ImageCacheConfigExpireType
	MemoryCache                          interface{}
	DiskCache                            interface{}
}

const kDefaultCacheMaxDiskAge = 60 * 60 * 24 * 7 // 1 week

var defaultImageCache ImageCacheConfig

func init() {
	defaultImageCache = NewImageCacheConfig()
}

func DefaultImageCache() ImageCacheConfig {
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
		MaxMemoryCost:                        0,
		DiskCacheExpireType:                  ImageCacheConfigExpireTypeModificationDate,
		MemoryCache:                          nil,
		DiskCache:                            nil,
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
		MaxMemoryCost:                        config.MaxMemoryCost,
		DiskCacheExpireType:                  config.DiskCacheExpireType,
		MemoryCache:                          config.MemoryCache,
		DiskCache:                            config.DiskCacheExpireType,
	}
}
