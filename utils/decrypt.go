package utils

import (
	"GiveMeAnOffer/custom_error"
	"GiveMeAnOffer/parse"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/grafov/m3u8"
	"golang.org/x/net/context"
)

type Cipher struct {
	KeyReq  *http.Request
	Key, IV []byte
	Method  string
	Ctx     context.Context

	MyKeyIV  string
	queryKey func(u string) ([]byte, bool)
	setKey   func(u string, key []byte)
}

type DecryptError struct{}

func (e DecryptError) Error() string {
	return "crypto/cipher: input not full blocks"
}

func NewCipherFromKey(config *parse.ParserTask, key *m3u8.Key, queryKey func(u string) ([]byte, bool), setKey func(u string, key []byte)) (*Cipher, error) {
	if key == nil || key.Method == "NONE" {
		return nil, nil
	}

	u, err := url.Parse(key.URI)
	if err != nil {
		return nil, err
	}

	if u.Scheme == "" || u.Host == "" {
		u, _ = url.Parse(config.Url)
		key.URI = fmt.Sprintf("%v://%v%v", u.Scheme, u.Host, key.URI)
	}

	decrypt := &Cipher{
		Method:  key.Method,
		IV:      []byte(key.IV),
		MyKeyIV: config.KeyIV,
	}
	req, err := http.NewRequest("GET", key.URI, nil)
	req = req.WithContext(config.Ctx)
	if err != nil {
		config.Logger.LogError(fmt.Sprintf("生成 密钥 Key 请求出粗：%v", err))
		return nil, err
	}
	decrypt.KeyReq = req
	decrypt.queryKey = queryKey
	decrypt.setKey = setKey
	return decrypt, err
}

// cbc解密模式
func (c *Cipher) aES128Decrypt(crypted []byte) ([]byte, error) {
	block, err := aes.NewCipher(c.Key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	if len(c.IV) == 0 {
		c.IV = c.Key
	}
	blockMode := cipher.NewCBCDecrypter(block, c.IV[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = pkcs5UnPadding(origData)
	return origData, nil
}

func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}

func (c *Cipher) Decrypt(body io.Reader) (*bytes.Buffer, error) {
	src, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}
	dst, err := c.aES128Decrypt(src)
	if err != nil {
		return nil, err
	}
	buffer := bytes.NewBuffer(dst)
	return buffer, err
}

func (c *Cipher) Generate(config *parse.ParserTask) error {
	var b []byte
	key := c.KeyReq.URL.String()

	if len(c.MyKeyIV) > 0 {
		b = []byte(c.MyKeyIV)
	} else if _, ok := c.queryKey(key); !ok {
		req := c.KeyReq
		if c.Ctx != nil {
			req = req.WithContext(c.Ctx)
		}

		// 下载 Key
		resp, err := config.Client.Do(req)
		if err != nil {
			config.Logger.LogErrorf("下载密钥失败：%v", err)
			return err
		}
		defer func(Body io.ReadCloser) {
			err = Body.Close()
			if err != nil {
				config.Logger.LogError(err.Error())
			}
		}(resp.Body)

		if resp.StatusCode != 200 {
			return custom_error.NetworkError{
				Code: resp.StatusCode,
				URL:  c.KeyReq.URL.String(),
			}
		}

		b, err = io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		c.setKey(key, b)
	} else {
		b, _ = c.queryKey(key)
	}
	c.Key = b
	return nil
}
