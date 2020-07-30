package test

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/bilibili-api"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestGetVideoInfo(t *testing.T) {
	req := bilibili_api.VideoInfoRequest{}
	req.Aid = "85440373"
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	item := <-req.Fetch(client).Observe()
	if item.E != nil {
		t.Fatal(item.E)
	}
	assert.NotNil(t, item.V)
}

func TestGetVideoDesc(t *testing.T) {
	req := bilibili_api.VideoDescRequest{}
	req.Aid = "85440373"
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	item := <-req.Fetch(client).Observe()
	if item.E != nil {
		t.Fatal(item.E)
	}
	assert.NotNil(t, item.V)
}

func TestFetchVideoPageList(t *testing.T) {
	req := bilibili_api.VideoPageListRequest{}
	req.Aid = "85440373"
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	item, err := req.Fetch(client).First().Get()
	if err != nil {
		t.Fatal(err)
	} else if item.E != nil {
		t.Fatal(item.E)
	}
	assert.NotNil(t, item.V)
}

func TestFetchVideoStreamSingle(t *testing.T) {
	req := bilibili_api.VideoStreamRequest{
		Bvid: "BV1y7411Q7Eq",
		Cid:  "171776208",
	}
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	item := <-req.Fetch(client).Observe()
	if item.E != nil {
		t.Fatal(item.E)
	}
	assert.NotNil(t, item.V)
}

func TestDownloadSegment(t *testing.T) {
	path, _ := os.UserHomeDir()
	path = filepath.Join(path, "Desktop", "download.flv")
	req := bilibili_api.VideoStreamRequest{
		Bvid: "BV117411r7R1",
		Cid:  "146044693",
	}
	req.SetProgressFunc(func(progress float64) {
		t.Logf("下载进度 %f", progress*100)
	})
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	item, err := req.Download(path, client).Get()
	if err != nil {
		t.Fatal(err)
	} else if item.E != nil {
		t.Fatal(item.E)
	}
	t.Log(item.V)
}

func TestDownloadByVideoInfo(t *testing.T) {
	req := bilibili_api.VideoInfoRequest{}
	req.Aid = "85440373"
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	item, err := req.Fetch(client).First().Get()
	if err != nil {
		t.Fatal(err)
	} else if item.E != nil {
		t.Fatal(item.E)
	}
	info := item.V.(bilibili_api.VideoInfo)
	info.SetProgressFunc(func(progress float64) {
		t.Logf("下载进度 %f", progress*100)
	})
	path, _ := os.UserHomeDir()
	path = filepath.Join(path, "Desktop", "download.flv")
	item, err = info.Download(path, client).Get()
	if err != nil {
		t.Fatal(err)
	} else if item.E != nil {
		t.Fatal(item.E)
	}
	t.Log(item.V)
}

func TestLogin(t *testing.T) {
	req := bilibili_api.LoginRequest{}
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	session := req.Login(client)
	t.Log(session)
}

func TestDownloadDanmuku(t *testing.T) {
	req := bilibili_api.HistoryDanmukuRequest{
		Oid:  "144541892",
		Date: "2020-01-21",
	}
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	path, _ := os.UserHomeDir()
	path = filepath.Join(path, "Desktop", "danmuku.xml")
	req.SetProgressFunc(func(progress float64) {
		t.Logf("下载进度 %f", progress*100)
	})
	item, err := req.Download(path, client).Get()
	if err != nil {
		t.Fatal(err)
	} else if item.E != nil {
		t.Fatal(item.E)
	}
	t.Log(item.V)
}

func TestDownloadDanmukuFromIndex(t *testing.T) {
	req := bilibili_api.HistoryDanmukuIndexRequest{
		Oid:   "144541892",
		Month: "2020-01",
	}
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	path, _ := os.UserHomeDir()
	path = filepath.Join(path, "Desktop")
	req.SetProgressFunc(func(progress float64) {
		t.Logf("下载进度 %f", progress*100)
	})
	item, err := req.Download(path, client).Get()
	if err != nil {
		t.Fatal(err)
	} else if item.E != nil {
		t.Fatal(item.E)
	}
	t.Log(item.V)
}
