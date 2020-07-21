package Core

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"sync"
)

type ImageDownloaderDecryptorBlock func(data []byte, response *http.Response) ([]byte, error)

type ImageDownloaderDecryptor interface {
	DecryptedWithResponse(data []byte, response *http.Response) ([]byte, error)
}

type imageDownloaderDecryptor struct {
	decrptorCb ImageDownloaderDecryptorBlock
}

func (decrypter *imageDownloaderDecryptor) DecryptedWithResponse(data []byte, response *http.Response) ([]byte, error) {
	if decrypter == nil || decrypter.decrptorCb == nil {
		return nil, fmt.Errorf("check input")
	}
	return decrypter.decrptorCb(data, response)
}

func (decrypter *imageDownloaderDecryptor) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func newImageDownloaderDecryptor(cb ImageDownloaderDecryptorBlock) *imageDownloaderDecryptor {
	return &imageDownloaderDecryptor{decrptorCb: cb}
}

var kBase64DownloaderToken sync.Once
var kBase64DownloaderDecryptor ImageDownloaderDecryptor

func GetBase64DownloaderDecryptor() ImageDownloaderDecryptor {
	kBase64DownloaderToken.Do(func() {
		if kBase64DownloaderDecryptor == nil {
			kBase64DownloaderDecryptor = newImageDownloaderDecryptor(func(data []byte, response *http.Response) (dst []byte, err error) {
				enc := base64.StdEncoding
				dst = make([]byte, enc.DecodedLen(len(data)))
				_, err = base64.StdEncoding.Decode(dst, data)
				return
			})
		}
	})
	return kBase64DownloaderDecryptor
}
