package Core

import (
	"net/http"
	"testing"
	"time"
)

func TestDownload(t *testing.T) {
	isExecuting := true
	GetSharedImageDownloader().DownloadImage("http://www.artslifenews.com/files/175148.40248437_1000X1000.jpg", 0, nil, nil, func(data []byte, err error, finished bool) {
		isExecuting = false
	})
	for isExecuting {

	}
}

func TestDownloadWithDependency(t *testing.T) {
	client := &http.Client{
		Timeout: 15 * time.Second,
	}
	isExecuting := true
	req, _ := http.NewRequest("GET", "http://www.artslifenews.com/files/175148.40248437_1000X1000.jpg", nil)
	op1 := NewImageDownloadOperation(req, client, 0, nil)
	op1.Start()
	req2, _ := http.NewRequest("GET", "http://www.artslifenews.com/files/175148.40248437_1000X1000.jpg", nil)
	op2 := NewImageDownloadOperation(req2, client, 0, nil)
	op2.AddDependency(op1)
	op2.AddHandlersForProgressAndCompletion(nil, func(data []byte, err error, finished bool) {
		isExecuting = false
		if err != nil {
			t.Fatal(err)
		} else {
			t.Log(data)
		}
	})
	op2.Start()
	for isExecuting {

	}
	t.Log("Done")
}
