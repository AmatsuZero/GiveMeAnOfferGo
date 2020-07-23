package Core

import (
	"context"
	"github.com/reactivex/rxgo/v2"
	"net/http"
	"sync"
	"testing"
	"time"
)

func TestDownload(t *testing.T) {
	wg := sync.WaitGroup{}
	GetSharedImageDownloader().DownloadImage("http://www.artslifenews.com/files/175148.40248437_1000X1000.jpg", 0, nil, nil, func(data []byte, err error, finished bool) {
		wg.Done()
	})
	wg.Wait()
}

func TestDownloadWithDependency(t *testing.T) {
	client := &http.Client{
		Timeout: 15 * time.Second,
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	req, _ := http.NewRequest("GET", "http://www.artslifenews.com/files/175148.40248437_1000X1000.jpg", nil)
	op1 := NewImageDownloadOperation(req, client, 0, nil)
	op1.Start()
	req2, _ := http.NewRequest("GET", "http://www.artslifenews.com/files/175148.40248437_1000X1000.jpg", nil)
	op2 := NewImageDownloadOperation(req2, client, 0, nil)
	op2.AddDependency(op1)
	op2.AddHandlersForProgressAndCompletion(nil, func(data []byte, err error, finished bool) {
		if err != nil {
			t.Fatal(err)
		}
		wg.Done()
	})
	op2.Start()
	wg.Wait()
	t.Log("Done")
}

func TestDownloadQueue(t *testing.T) {
	client := &http.Client{
		Timeout: 15 * time.Second,
	}
	wg := sync.WaitGroup{}
	wg.Add(2)
	req3, _ := http.NewRequest("GET", "http://www.artslifenews.com/files/175148.40248437_1000X1000.jpg", nil)
	op3 := NewImageDownloadOperation(req3, client, 0, nil)
	op3.AddHandlersForProgressAndCompletion(nil, func(data []byte, err error, finished bool) {
		t.Log("op3 done")
	})

	req, _ := http.NewRequest("GET", "http://www.artslifenews.com/files/175148.40248437_1000X1000.jpg", nil)
	op1 := NewImageDownloadOperation(req, client, 0, nil)
	op1.AddHandlersForProgressAndCompletion(nil, func(data []byte, err error, finished bool) {
		t.Log("op1 done")
		op3.Start()
		wg.Done()
	})
	req2, _ := http.NewRequest("GET", "http://www.artslifenews.com/files/175148.40248437_1000X1000.jpg", nil)
	op2 := NewImageDownloadOperation(req2, client, 0, nil)
	op2.AddHandlersForProgressAndCompletion(nil, func(data []byte, err error, finished bool) {
		t.Log("op2 done")
		wg.Done()
	})
	op2.AddDependency(op3)

	q := NewImageDownloadQueue(1)
	q.AddOperation(op1)
	q.AddOperation(op2)
	wg.Wait()
	t.Log("Done")
}

func TestDefer(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	ctx := context.TODO()
	ob := rxgo.Defer([]rxgo.Producer{func(ctx context.Context, next chan<- rxgo.Item) {

		next <- rxgo.Of(1)
	}}, rxgo.WithContext(ctx))
	ob.DoOnCompleted(func() {
		wg.Done()
	})
	wg.Wait()
}
