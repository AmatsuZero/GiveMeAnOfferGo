package Core

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

type ImageDownloaderDecryptorBlock func(response *http.Response) ([]byte, error)

type ImageDownloaderDecryptor interface {
	DecryptedWithResponse(response *http.Response) ([]byte, error)
}

type imageDownloaderDecryptor struct {
	decrptorCb ImageDownloaderDecryptorBlock
}

func (decrypter *imageDownloaderDecryptor) DecryptedWithResponse(response *http.Response) ([]byte, error) {
	if decrypter == nil || decrypter.decrptorCb == nil {
		return nil, fmt.Errorf("check input")
	}
	return decrypter.decrptorCb(response)
}

func newImageDownloaderDecryptor(cb ImageDownloaderDecryptorBlock) *imageDownloaderDecryptor {
	return &imageDownloaderDecryptor{decrptorCb: cb}
}

var kBase64DownloaderToken sync.Once
var kBase64DownloaderDecryptor ImageDownloaderDecryptor

func GetBase64DownloaderDecryptor() ImageDownloaderDecryptor {
	kBase64DownloaderToken.Do(func() {
		if kBase64DownloaderDecryptor == nil {
			kBase64DownloaderDecryptor = newImageDownloaderDecryptor(func(r *http.Response) (dst []byte, err error) {
				src, err := ioutil.ReadAll(r.Body)
				if err != nil {
					return nil, err
				}
				enc := base64.StdEncoding
				dst = make([]byte, enc.DecodedLen(len(src)))
				_, err = base64.StdEncoding.Decode(dst, src)
				return
			})
		}
	})
	return kBase64DownloaderDecryptor
}
