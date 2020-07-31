package bilibili_api

import (
	"bufio"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/reactivex/rxgo/v2"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

type VideoStreamResolution int

const (
	VideoStreamResolution360             VideoStreamResolution = 16
	VideoStreamResolution480                                   = VideoStreamResolution360 * 2
	VideoStreamResolution720                                   = VideoStreamResolution480 * 2
	VideoStreamResolutionVIP720                                = VideoStreamResolution720 + 10
	VideoStreamResolution1080                                  = VideoStreamResolutionVIP720 + 6
	VideoStreamResolutionVIP1080Plus                           = VideoStreamResolution1080 + 32
	VideoStreamResolutionVIP1080and60FPS                       = VideoStreamResolutionVIP1080Plus + 4
	VideoStreamResolutionVIP4K                                 = VideoStreamResolutionVIP1080and60FPS + 4
)

func (resolution VideoStreamResolution) String() string {
	switch resolution {
	case VideoStreamResolutionVIP4K:
		return "4K 超清（大会员）"
	case VideoStreamResolutionVIP1080and60FPS:
		return "1080P60 高清（大会员）"
	case VideoStreamResolutionVIP1080Plus:
		return "1080P+ 高清（大会员）"
	case VideoStreamResolution1080:
		return "1080P 高清（登录）"
	case VideoStreamResolutionVIP720:
		return "720P60 高清（大会员）"
	case VideoStreamResolution720:
		return "720P 高清（登录）"
	case VideoStreamResolution480:
		return "480P 清晰"
	default:
		return "360P 流畅"
	}
}

func (resolution VideoStreamResolution) Encode() string {
	return fmt.Sprintf("%d", resolution)
}

type videoStreamBaseRequest struct {
	baseRequest
	progressCb func(progress float64)
	tmpDir     string
}

func (request *videoStreamBaseRequest) SetProgressFunc(cb func(progress float64)) {
	request.progressCb = cb
}

func (request videoStreamBaseRequest) Download(to string, client *http.Client, opts ...rxgo.Option) rxgo.OptionalSingle {
	defaultClient := http.DefaultClient // 这里改用没有超时的默认 Client，避免任务被 Cancel
	ob := request.download(request.Fetch(client), defaultClient, opts...)
	return request.export(ob, to)
}

func (request *videoStreamBaseRequest) download(ob rxgo.Observable, client *http.Client, opts ...rxgo.Option) rxgo.Observable {
	return ob.FlatMap(func(item rxgo.Item) rxgo.Observable {
		if item.E != nil {
			return rxgo.Thrown(item.E)
		}
		info := item.V.(videoStreamInfoProtocol)
		tmpDir := info.GetTmpDirName()
		_, err := os.Stat(tmpDir)
		if os.IsNotExist(err) {
			err = os.MkdirAll(tmpDir, os.ModePerm)
		}
		if err != nil {
			return rxgo.Thrown(item.E)
		}
		request.tmpDir = tmpDir
		info.SetProgressFunc(request.progressCb)
		return info.download(tmpDir, client, opts...)
	})
}

func (request videoStreamBaseRequest) export(ob rxgo.Observable, to string) rxgo.OptionalSingle {
	return ob.Reduce(func(ctx context.Context, acc interface{}, elem interface{}) (interface{}, error) {
		var ret []string
		if acc == nil {
			ret = make([]string, 0)
		} else {
			ret = acc.([]string)
		}
		ret = append(ret, elem.(string))
		return ret, nil
	}).Map(func(ctx context.Context, i interface{}) (interface{}, error) {
		p := filepath.Join(request.tmpDir, "ff.txt")
		manifest, err := os.Create(p)
		if err != nil {
			_ = os.RemoveAll(request.tmpDir)
			return nil, err
		}
		defer func() {
			_ = manifest.Close()
		}()
		dataWriter := bufio.NewWriter(manifest)
		for _, file := range i.([]string) {
			_, _ = dataWriter.WriteString("file " + "'" + file + "'" + "\n")
		}
		err = dataWriter.Flush()
		return p, err
	}).Map(func(ctx context.Context, i interface{}) (interface{}, error) {
		defer func() {
			_ = os.RemoveAll(request.tmpDir)
		}()
		err := exec.CommandContext(ctx, "ffmpeg", "-f", "concat", "-safe", "0", "-i", i.(string), "-c", "copy", to).Run()
		return to, err
	})
}

type VideoStreamRequest struct {
	videoStreamBaseRequest
	Avid        string
	Bvid        string
	Cid         string
	Resolution  VideoStreamResolution
	ShouldUse4K bool
}

func (request VideoStreamRequest) IsParamsValid() bool {
	return (len(request.Avid) > 0 || len(request.Bvid) > 0) && len(request.Cid) > 0
}

func (request VideoStreamRequest) Request() (*http.Request, error) {
	if !request.IsParamsValid() {
		return nil, kInvalidParamError
	}
	base := *kBaseURL
	base.Path = "x/player/playurl"
	base.RawQuery = request.queryItems(base.Query()).Encode()
	return http.NewRequest("GET", base.String(), nil)
}

func (request VideoStreamRequest) queryItems(query url.Values) url.Values {
	if len(request.Avid) > 0 {
		query.Add("avid", request.Avid)
	} else {
		query.Add("bvid", request.Bvid)
	}
	query.Set("cid", request.Cid)
	if request.Resolution > 0 {
		query.Set("qn", request.Resolution.Encode())
	}
	if request.ShouldUse4K {
		query.Set("fourk", "1")
	}
	return query
}

func (request VideoStreamRequest) Download(to string, client *http.Client, opts ...rxgo.Option) rxgo.OptionalSingle {
	defaultClient := http.DefaultClient // 这里改用没有超时的默认 Client，避免任务被 Cancel
	ob := request.download(request.Fetch(client), defaultClient, opts...)
	return request.export(ob, to)
}

func (request VideoStreamRequest) Fetch(client *http.Client, opts ...rxgo.Option) rxgo.Observable {
	req, err := request.Request()
	if err != nil {
		return rxgo.Thrown(err)
	}
	return request.fetch(client, req, opts...).Map(func(ctx context.Context, i interface{}) (interface{}, error) {
		info := &VideoStreamInfo{} // interface 要使用，必须 cast 为指针类型，参考：https://forum.golangbridge.org/t/how-to-cast-interface-to-a-given-interface/13997
		data := i.([]byte)
		err := json.Unmarshal(data, info)
		if err != nil {
			return nil, err
		}
		info.md5Sign = md5.Sum(data)
		return info, nil
	})
}

type VideoSegments struct {
	baseRequest
	Order      int // 分段序号
	Length     int
	Size       int // 单位为Byte
	Ahead      string
	Vhead      string
	Url        string   // 视频流url
	BackupUrl  []string `json:"backup_url"`
	retryCount int      // 标志当前重试次数
	tmpDir     string   // 存放分片文件的文件夹
	progressCB func(segmentSize, acceptedSize int)
}

func (segments VideoSegments) Request() (req *http.Request, err error) {
	if segments.retryCount > len(segments.BackupUrl) {
		return nil, kOutOfTimesError
	}
	if segments.retryCount == 0 {
		req, err = http.NewRequest("GET", segments.Url, nil)
	} else {
		req, err = http.NewRequest("GET", segments.BackupUrl[segments.retryCount-1], nil)
	}
	if err != nil {
		return nil, err
	}
	req.Header.Set("user-agent", kFakeWebUA)
	req.Header.Set("referer", kFakeRefer)
	return
}

func (segments VideoSegments) IsParamsValid() bool {
	if segments.retryCount == 0 {
		return len(segments.Url) > 0
	}
	return len(segments.BackupUrl[segments.retryCount-1]) > 0
}

func (segments *VideoSegments) Fetch(client *http.Client, opts ...rxgo.Option) rxgo.Observable {
	if segments == nil {
		return rxgo.Empty()
	}
	if client == nil {
		client = http.DefaultClient
	}
	req, err := segments.Request()
	if err != nil {
		return rxgo.Thrown(err)
	}
	if segments.Session == nil {
		segments.Session = kDefaultSession
	}
	if client.Jar == nil {
		for _, c := range segments.Session.Cookies(req.URL) {
			req.AddCookie(c)
		}
	}
	return rxgo.Defer([]rxgo.Producer{func(ctx context.Context, next chan<- rxgo.Item) {
		req = req.WithContext(ctx)
		resp, err := client.Do(req)
		if err == nil && client.Jar == nil {
			kDefaultSession.SetCookies(resp.Request.URL, resp.Cookies())
			_ = kDefaultSession.Serialize(kDefaultSessionPath)
		}
		next <- rxgo.Item{
			V: resp,
			E: err,
		}
	}}, opts...).Retry(len(segments.BackupUrl), func(err error) bool {
		if e, ok := err.(BadResponseError); ok && e.Code == http.StatusForbidden { // 禁止请求，不必重试了
			return false
		}
		segments.retryCount++
		return true
	}).Map(segments.writeToFile, opts...)
}

func (segments *VideoSegments) writeToFile(ctx context.Context, i interface{}) (interface{}, error) {
	if segments == nil || len(segments.tmpDir) == 0 {
		return "", kInvalidParamError
	}
	r, ok := i.(*http.Response)
	if !ok {
		return nil, kInvalidParamError
	}
	defer func() {
		_ = r.Body.Close()
	}()
	if r.StatusCode == http.StatusForbidden {
		return nil, BadResponseError{
			Code:    r.StatusCode,
			Message: r.Status,
		}
	}
	dst := filepath.Join(segments.tmpDir, path.Base(r.Request.URL.Path))
	file, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = file.Close()
	}()
	_, err = io.Copy(file, readerFunc(func(p []byte) (n int, err error) {
		// golang non-blocking channel: https://gobyexample.com/non-blocking-channel-operations
		select {
		// if context has been canceled
		case <-ctx.Done():
			// stop process and propagate "context canceled" error
			return 0, ctx.Err()
		default:
			// otherwise just run default io.Reader implementation
			n, err = r.Body.Read(p)
			segments.progressCB(segments.Size, n)
			return n, err
		}
	}))
	return dst, err
}

type VideoStreamInfoData struct {
	From              string
	Result            string
	Message           string
	Quality           VideoStreamResolution
	Format            string
	TimeLength        int                     `json:"timelength"`
	AcceptFormat      string                  `json:"accept_format"`
	AcceptDescription []string                `json:"accept_description"`
	AcceptQuality     []VideoStreamResolution `json:"accept_quality"`
	VideoCodecid      int                     `json:"video_codecid"`
	SeekParam         string                  `json:"seek_param"`
	SeekType          string                  `json:"seek_type"`
	Durl              []VideoSegments
}

type videoStreamInfoProtocol interface {
	Downloadable
	GetTmpDirName() string
	GetSegments() []VideoSegments
	TotalSize() int64
	download(tmpDir string, client *http.Client, opts ...rxgo.Option) rxgo.Observable
}

type baseVideoStreamInfo struct {
	md5Sign    [16]byte
	progressCB func(progress float64)
}

func (info baseVideoStreamInfo) GetTmpDirName() string {
	return filepath.Join(os.TempDir(), fmt.Sprintf("%x", info.md5Sign))
}

func (info baseVideoStreamInfo) GetSegments() []VideoSegments {
	return []VideoSegments{}
}

func (info *baseVideoStreamInfo) SetProgressFunc(cb func(progress float64)) {
	info.progressCB = cb
}

func (info baseVideoStreamInfo) Download(to string, client *http.Client, opts ...rxgo.Option) rxgo.OptionalSingle {
	return info.download(to, client, opts...).First()
}

func (info baseVideoStreamInfo) download(tmpDir string, client *http.Client, opts ...rxgo.Option) rxgo.Observable {
	return info.operation(info.GetSegments(), info.TotalSize(), tmpDir, client, opts...)
}

func (info baseVideoStreamInfo) operation(segments []VideoSegments, size int64,
	tmpDir string, client *http.Client, opts ...rxgo.Option) rxgo.Observable {
	totalSize := float64(size)
	receivedData := float64(0)
	return rxgo.Just(segments)(rxgo.Serialize(func(i interface{}) int {
		return i.(VideoSegments).Order
	})).FlatMap(func(item rxgo.Item) rxgo.Observable {
		if item.E != nil {
			return rxgo.Thrown(item.E)
		}
		seg := item.V.(VideoSegments)
		seg.tmpDir = tmpDir
		seg.progressCB = func(segmentSize, acceptedSize int) {
			receivedData += float64(acceptedSize)
			if info.progressCB != nil {
				info.progressCB(receivedData / totalSize)
			}
		}
		return seg.Fetch(client, opts...)
	})
}

func (info baseVideoStreamInfo) TotalSize() (size int64) {
	for _, seg := range info.GetSegments() {
		size += int64(seg.Size)
	}
	return
}

type VideoStreamInfo struct {
	BaseResponse
	baseVideoStreamInfo
	Data VideoStreamInfoData
}

func (info VideoStreamInfo) TotalSize() (size int64) {
	for _, seg := range info.GetSegments() {
		size += int64(seg.Size)
	}
	return
}

func (info VideoStreamInfo) GetSegments() []VideoSegments {
	return info.Data.Durl
}

func (info VideoStreamInfo) SupportFormats() []string {
	return strings.Split(info.Data.AcceptFormat, ",")
}

func (info VideoStreamInfo) download(tmpDir string, client *http.Client, opts ...rxgo.Option) rxgo.Observable {
	return info.operation(info.GetSegments(), info.TotalSize(), tmpDir, client, opts...)
}

func (info VideoStreamInfo) Download(to string, client *http.Client, opts ...rxgo.Option) rxgo.OptionalSingle {
	return info.download(to, client, opts...).First()
}
