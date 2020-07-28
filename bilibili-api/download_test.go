package bilibili_api

import (
	"context"
	"github.com/reactivex/rxgo/v2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestGetVideoInfo(t *testing.T) {
	req := VideoInfoRequest{}
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

func TestGetVideoDesc(t *testing.T) {
	req := VideoDescRequest{}
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

func TestFetchVideoPageList(t *testing.T) {
	req := VideoPageListRequest{}
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
	req := VideoStreamRequest{
		Bvid: "BV1y7411Q7Eq",
		Cid:  "171776208",
	}
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

func TestDownloadSegment(t *testing.T) {
	path := "/Users/jiangzhenhua/Desktop/download.flv"
	req := VideoStreamRequest{
		Bvid: "BV117411r7R1",
		Cid:  "146044693",
	}
	req.ProgressCb = func(progress float64) {
		t.Logf("下载进度 %f", progress*100)
	}
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

func TestReduce(t *testing.T) {
	ob := rxgo.Range(0, 100).Reduce(func(ctx context.Context, acc interface{}, elem interface{}) (interface{}, error) {
		var ret []int
		if acc == nil {
			ret = make([]int, 0)
		} else {
			ret = acc.([]int)
		}
		ret = append(ret, elem.(int))
		return ret, nil
	})
	item, _ := ob.Get()
	t.Log(item.V)
}
