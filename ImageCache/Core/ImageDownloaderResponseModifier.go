package Core

import "net/http"

type ImageDownloaderResponseModifierBlock func(response *http.Response) *http.Response

type ImageDownloaderResponseModifier interface {
	ModifiedResponseWithResponse(response *http.Response) *http.Response
}

type imageDownloaderResponseModifier struct {
	block ImageDownloaderResponseModifierBlock
}

func (modifier *imageDownloaderResponseModifier) ModifiedResponseWithResponse(response *http.Response) *http.Response {
	if modifier == nil || modifier.block == nil {
		return nil
	}
	return modifier.block(response)
}

func newImageDownloaderResponseModifier(statusCode int, version string, headers http.Header) *imageDownloaderResponseModifier {
	if len(version) == 0 {
		version = "HTTP/1.1"
	}
	return newImageDownloaderResponseModifierWithBlock(func(response *http.Response) *http.Response {
		if response == nil {
			return response
		}
		for k := range headers {
			response.Header.Set(k, headers.Get(k))
		}
		response.StatusCode = statusCode
		response.Proto = version
		return response
	})
}

func newImageDownloaderResponseModifierWithBlock(block ImageDownloaderResponseModifierBlock) *imageDownloaderResponseModifier {
	return &imageDownloaderResponseModifier{block: block}
}
