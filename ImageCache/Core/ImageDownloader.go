package Core

type ImageDownloaderOptions BitsType

const (
	mageDownloaderLowPriority ImageDownloaderOptions = 1 << iota
	ImageDownloaderProgressiveLoad
	ImageDownloaderUseNSURLCache
	ImageDownloaderIgnoreCachedResponse
	ImageDownloaderContinueInBackground
	ImageDownloaderHandleCookies
	ImageDownloaderHighPriority
	ImageDownloaderScaleDownLargeImages
	ImageDownloaderAvoidDecodeImage
	ImageDownloaderDecodeFirstFrameOnly
	ImageDownloaderPreloadAllFrames
	ImageDownloaderMatchAnimatedImageClass
)
