package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"github.com/grafov/m3u8"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.org/x/net/context"
	"io"
	"net/http"
	"net/url"
)

type Cipher struct {
	KeyReq *http.Request
	IV     string
	Method string
	Ctx    context.Context
	len    int

	MyKeyIV  string
	block    cipher.Block
	queryKey func(u string) ([]byte, bool)
	setKey   func(u string, key []byte)
}

func NewCipherFromKey(config *ParserTask, key *m3u8.Key, queryKey func(u string) ([]byte, bool), setKey func(u string, key []byte)) (*Cipher, error) {
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
		IV:      key.IV,
		MyKeyIV: config.KeyIV,
	}
	req, err := http.NewRequest("GET", key.URI, nil)
	req = req.WithContext(SharedApp.ctx)
	if err != nil {
		runtime.LogError(SharedApp.ctx, fmt.Sprintf("生成 密钥 Key 请求出粗：%v", err))
		return nil, err
	}
	decrypt.KeyReq = req
	decrypt.queryKey = queryKey
	decrypt.setKey = setKey
	return decrypt, err
}

func (c *Cipher) Decrypt(body io.Reader) (*bytes.Buffer, error) {
	// cbc解密模式
	src, err := io.ReadAll(body)
	blockSize := c.block.BlockSize()

	if len(src)%blockSize != 0 {
		return nil, fmt.Errorf("crypto/cipher: input not full blocks, len: %v, block size: %v", len(src), blockSize)
	}

	key, _ := c.queryKey(c.KeyReq.RequestURI)
	var iv []byte
	if len(c.IV) == 0 {
		iv = key[:blockSize]
	} else {
		iv = []byte(c.IV)
	}
	blockMode := cipher.NewCBCDecrypter(c.block, iv)
	dst := make([]byte, len(src))
	blockMode.CryptBlocks(dst, src)
	buffer := bytes.NewBuffer(dst)
	return buffer, err
}

func (c *Cipher) Generate() error {
	if c.block != nil {
		return nil
	}
	var b []byte
	if len(c.MyKeyIV) > 0 {
		b = []byte(c.MyKeyIV)
	} else if _, ok := c.queryKey(c.KeyReq.RequestURI); !ok {
		req := c.KeyReq
		if c.Ctx != nil {
			req = req.WithContext(c.Ctx)
		}

		// 下载 Key
		resp, err := SharedApp.client.Do(req)
		if err != nil {
			runtime.LogError(SharedApp.ctx, fmt.Sprintf("下载密钥失败：%v", err))
			return err
		}
		defer func(Body io.ReadCloser) {
			err = Body.Close()
			if err != nil {
				runtime.LogError(SharedApp.ctx, err.Error())
			}
		}(resp.Body)

		if resp.StatusCode != 200 {
			return fmt.Errorf("下载失败：Received HTTP %v for %v", resp.StatusCode, c.KeyReq.URL.String())
		}

		b, err = io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		c.setKey(c.KeyReq.RequestURI, b)
	} else {
		b, _ = c.queryKey(c.KeyReq.RequestURI)
	}
	block, err := aes.NewCipher(b)
	if err != nil {
		return err
	}
	c.block = block
	return nil
}
