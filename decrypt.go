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

func NewCipherFromKey(key *m3u8.Key, myKeyIV string, queryKey func(u string) ([]byte, bool), setKey func(u string, key []byte)) (*Cipher, error) {
	if key == nil || key.Method == "NONE" {
		return nil, nil
	}

	decrypt := &Cipher{
		Method:  key.Method,
		IV:      key.IV,
		MyKeyIV: myKeyIV,
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

	b, ok := c.queryKey(c.KeyReq.RequestURI)
	if !ok {
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
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return fmt.Errorf("下载失败：Received HTTP %v for %v", resp.StatusCode, c.KeyReq.URL.String())
		}

		b, err = io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		c.setKey(c.KeyReq.RequestURI, b)
	}

	if b == nil && len(c.MyKeyIV) > 0 {
		b = []byte(c.MyKeyIV)
	}

	block, err := aes.NewCipher(b)
	if err != nil {
		return err
	}
	c.block = block
	return nil
}
