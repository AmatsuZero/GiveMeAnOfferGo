package Core

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

type ImageDownloaderRequestModifierBlock func(req *http.Request) *http.Request

type ImageDownloaderRequestModifier interface {
	ModifiedRequest(req *http.Request) *http.Request
}

type imageDownloaderRequestModifier struct {
	block ImageDownloaderRequestModifierBlock
}

func newImageDownloaderRequestModifier(method string, header http.Header, body []byte) *imageDownloaderRequestModifier {
	return newImageDownloaderRequestModifierWithBlock(func(req *http.Request) *http.Request {
		req.Method = method
		req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		for k := range header {
			req.Header.Set(k, header.Get(k))
		}
		return req
	})
}

func newImageDownloaderRequestModifierWithBlock(block ImageDownloaderRequestModifierBlock) *imageDownloaderRequestModifier {
	return &imageDownloaderRequestModifier{block: block}
}
