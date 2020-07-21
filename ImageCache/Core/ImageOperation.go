package Core

import (
	"bytes"
	"context"
	"fmt"
	"github.com/reactivex/rxgo/v2"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"sync"
	"time"
)

type ImageDownloaderProgressBlock func(receivedSize, expectedSize int64, targetURL *url.URL)
type ImageDownloaderCompletedBlock func(data []byte, err error, finished bool)

var (
	ImageOperationCancelError                    error
	ImageOperationInvalidError                   error
	ImageOperationInvalidResponseError           error
	ImageOperationInvalidDownloadStatusCodeError error
	ImageOperationCacheNoModified                error
)

func init() {
	ImageOperationCancelError = fmt.Errorf("operation cancelled by user before sending the request")
	ImageOperationInvalidError = fmt.Errorf("task can't be initialized")
	ImageOperationInvalidResponseError = fmt.Errorf("download marked as failed because response is nil")
	ImageOperationInvalidDownloadStatusCodeError = fmt.Errorf("download marked as failed because response status code is not in 200-400")
	ImageOperationCacheNoModified = fmt.Errorf("download response status code is 304 not modified and ignored")
}

type WebImageOperationProtocol interface {
	Cancel()
}

type WebImageOperation struct {
	cancel     rxgo.Disposable
	task       rxgo.Observable
	ctx        context.Context
	IsCanceled bool
	IsRunning  bool
	IsFinished bool
}

func (op *WebImageOperation) Cancel() {
	if op == nil || op.cancel == nil || op.IsCanceled {
		return
	}
	op.cancel()
}

func (op *WebImageOperation) Start() {
	if op == nil || op.task == nil {
		return
	}
	ctx, cancel := op.task.Connect()
	op.IsRunning = true
	op.cancel = cancel
	op.ctx = ctx
	go func(ctx context.Context) {
		select {
		case <-ctx.Done():
			op.IsRunning = false
			op.IsCanceled = true
		default:
			op.IsRunning = true
			op.IsCanceled = false
		}
	}(ctx)
}

func newWebImageOperation(task rxgo.Observable) *WebImageOperation {
	return &WebImageOperation{task: task}
}

type ImageDownloadOperationProtocol interface {
	Cancel(token interface{}) bool
	AddHandlersForProgressAndCompletion(
		progressCb ImageDownloaderProgressBlock,
		completion ImageDownloaderCompletedBlock) interface{}
	GetResponse() *http.Response
	GetRequest() *http.Request
	Start()
}

type kCallbacksDictionary map[string]interface{}

const kProgressCallbackKey = "progress"
const kCompletedCallbackKey = "completed"

type ImageDownloadOperation struct {
	WebImageOperation
	req              *http.Request
	client           *http.Client
	options          ImageDownloaderOptions
	context          ImageContext
	callbackBlocks   []kCallbacksDictionary
	lock             sync.Mutex
	isFinished       bool
	cachedData       []byte
	total            int64 // Total # of bytes transferred
	expectedSize     int64
	taskErr          error
	responseModifier ImageDownloaderResponseModifier // 用来修改原来的相应体
	decryption       ImageDownloaderDecryptor        // 用来解密图像数据
	response         *http.Response
}

func (op *ImageDownloadOperation) GetRequest() *http.Request {
	if op == nil || op.req == nil {
		return nil
	}
	return op.req
}

func (op *ImageDownloadOperation) GetResponse() *http.Response {
	if op == nil || op.response == nil {
		return nil
	}
	return op.response
}

func (op *ImageDownloadOperation) Cancel(token interface{}) (shouldCancel bool) {
	if op == nil || token == nil {
		return false
	}
	op.lock.Lock()
	tmpCallbacks := op.callbackBlocks[:0]
	for _, cb := range op.callbackBlocks {
		if !reflect.DeepEqual(cb, token) {
			tmpCallbacks = append(tmpCallbacks, cb)
		}
	}
	shouldCancel = len(tmpCallbacks) == 0
	op.lock.Unlock()
	if shouldCancel {
		op.cancel() // 取消最后一个正在运行中的任务，并唤醒最后一个回调
	} else {
		op.lock.Lock()
		defer op.lock.Unlock()
		op.callbackBlocks = tmpCallbacks
		t, ok := token.(kCallbacksDictionary)
		if !ok {
			return
		}
		cb, ok := t[kCompletedCallbackKey]
		if !ok {
			return
		}
		block, ok := cb.(ImageDownloaderCompletedBlock)
		if !ok {
			return
		}
		go block(nil, ImageOperationCancelError, true)
	}
	return
}

func (op *ImageDownloadOperation) AddHandlersForProgressAndCompletion(
	progressCb ImageDownloaderProgressBlock,
	completionCb ImageDownloaderCompletedBlock) interface{} {
	if op == nil {
		return nil
	}
	callbacks := kCallbacksDictionary{}
	if progressCb != nil {
		callbacks[kProgressCallbackKey] = progressCb
	}
	if completionCb != nil {
		callbacks[kCompletedCallbackKey] = completionCb
	}
	op.lock.Lock()
	op.callbackBlocks = append(op.callbackBlocks, callbacks)
	op.lock.Unlock()
	return callbacks
}

func (op *ImageDownloadOperation) callbacksForKey(key string) []interface{} {
	if op == nil {
		return nil
	}
	callbacks := make([]interface{}, 0)
	op.lock.Lock()
	for _, cb := range op.callbackBlocks {
		if fn, ok := cb[key]; ok {
			callbacks = append(callbacks, fn)
		}
	}
	op.lock.Unlock()
	return callbacks
}

func (op *ImageDownloadOperation) cancel() {
	if op == nil {
		return
	}
	op.lock.Lock()
	op.cancelInternal()
	op.lock.Unlock()
}

func (op *ImageDownloadOperation) cancelInternal() {
	if op == nil || op.isFinished {
		return
	}
	if op.task != nil {
		op.WebImageOperation.Cancel()
		if op.IsRunning {
			op.IsRunning = false
		}
		if !op.isFinished {
			op.isFinished = true
		}
	} else {
		op.callCompletionBlocksWithError(ImageOperationCancelError)
	}
	op.reset()
}

func (op *ImageDownloadOperation) Start() {
	if op == nil {
		return
	}
	if op.IsCanceled {
		op.isFinished = true
		op.callCompletionBlocksWithError(ImageOperationCancelError)
		op.reset()
		return
	}
	if op.client == nil {
		op.client = &http.Client{
			Timeout: time.Second * 15,
		}
	}
	if BitsHas(BitsType(op.options), BitsType(ImageDownloaderIgnoreCachedResponse)) {
		resp, ok := getURLCache().GetCachedResponseForRequest(op.req)
		if ok {
			op.cachedData = resp.BodyData
		}
	}
	op.task = rxgo.Create([]rxgo.Producer{func(ctx context.Context, next chan<- rxgo.Item) {
		data, err := op.download()
		next <- rxgo.Item{
			V: data,
			E: err,
		}
	}})
	op.task.DoOnCompleted(op.done)
	op.WebImageOperation.Start()
}

func (op *ImageDownloadOperation) download() (data []byte, err error) {
	defer func() {
		op.cachedData = data
		op.taskErr = err
	}()
	size, err := op.headRequest()
	if err != nil {
		return
	}
	op.expectedSize = int64(size)
	resp, err := op.client.Do(op.req)
	if err != nil {
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	var dst bytes.Buffer
	src := io.TeeReader(resp.Body, op)
	_, err = io.Copy(&dst, src)
	op.response = resp
	// 放入缓存
	getURLCache().AddCache(op.req, resp)
	return dst.Bytes(), err
}

// Head 请求，获取大小, 验证 response
func (op *ImageDownloadOperation) headRequest() (int, error) {
	if op == nil || op.req == nil {
		return 0, ImageOperationInvalidError
	}
	headResp, err := op.client.Head(op.req.URL.String())
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = headResp.Body.Close()
	}()
	if op.responseModifier != nil {
		headResp = op.responseModifier.ModifiedResponseWithResponse(headResp)
		if headResp == nil {
			return 0, ImageOperationInvalidResponseError
		}
	}
	statusCode := headResp.StatusCode
	statusCodeValid := statusCode >= 200 && statusCode < 400
	if !statusCodeValid {
		return 0, ImageOperationInvalidDownloadStatusCodeError
	}
	if statusCode == 304 && len(op.cachedData) == 0 {
		return 0, ImageOperationCacheNoModified
	}
	// 进度 0 开始
	_, err = op.Write(nil)
	return strconv.Atoi(headResp.Header.Get("Content-Length"))
}

func (op *ImageDownloadOperation) Write(p []byte) (n int, err error) {
	n = len(p)
	op.total += int64(n)
	for _, cb := range op.callbacksForKey(kProgressCallbackKey) {
		block, ok := cb.(ImageDownloaderProgressBlock)
		if !ok {
			continue
		}
		block(op.total, op.expectedSize, op.req.URL)
	}
	return
}

func (op *ImageDownloadOperation) reset() {
	if op == nil {
		return
	}
	op.lock.Lock()
	defer op.lock.Unlock()
	op.callbackBlocks = op.callbackBlocks[:0]
	op.task = nil
	op.client = nil
	op.total = 0
	op.expectedSize = 0
	op.cachedData = nil
	op.taskErr = nil
	op.response = nil
}

func (op *ImageDownloadOperation) done() {
	if op == nil {
		return
	}
	op.isFinished = true
	op.IsRunning = false
	data, err := op.cachedData, op.taskErr
	if data != nil && op.decryption != nil {
		data, err = op.decryption.DecryptedWithResponse(data, op.response)
	}
	op.callCompletionBlocksWithImageData(data, err, true)
	op.reset()
}

func (op *ImageDownloadOperation) callCompletionBlocksWithError(err error) {
	op.callCompletionBlocksWithImageData(nil, err, true)
}

func (op *ImageDownloadOperation) callCompletionBlocksWithImageData(data []byte, err error, finished bool) {
	if op == nil {
		return
	}
	rxgo.Just(op.callbacksForKey(kCompletedCallbackKey))().DoOnNext(func(i interface{}) {
		if block, ok := i.(ImageDownloaderCompletedBlock); ok {
			block(data, err, finished)
		}
	})
}

func NewImageDownloadOperation(req *http.Request, client *http.Client,
	options ImageDownloaderOptions, ctx ImageContext) *ImageDownloadOperation {
	var modifier ImageDownloaderResponseModifier
	if v, ok := ctx[kImageContextDownloadResponseModifier]; ok {
		cb, ok := v.(ImageDownloaderResponseModifier)
		if ok {
			modifier = cb
		}
	}
	var decryption ImageDownloaderDecryptor
	if v, ok := ctx[kImageContextDownloadDecryptor]; ok {
		cb, ok := v.(ImageDownloaderDecryptor)
		if ok {
			decryption = cb
		}
	}
	return &ImageDownloadOperation{
		req:              req,
		client:           client,
		options:          options,
		context:          ctx,
		callbackBlocks:   make([]kCallbacksDictionary, 0),
		responseModifier: modifier,
		decryption:       decryption,
	}
}
