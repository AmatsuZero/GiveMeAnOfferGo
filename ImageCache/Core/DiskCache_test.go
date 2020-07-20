package Core

import (
	"net/http"
	"net/url"
	"testing"
)

func TestMoveDir(t *testing.T) {
	isExecuting := true
	req, _ := http.NewRequest("GET", "http://img.chewen.com/pics/2011/09/29/4112011092915584230_b.jpg", nil)
	op := &imageDownloadOperation{req: req}
	op.AddHandlersForProgressAndCompletion(func(receivedSize, expectedSize int64, targetURL *url.URL) {
		t.Logf("Read %d bytes for a total of %d\n", receivedSize, expectedSize)
	}, func(data []byte, err error, finished bool) {
		isExecuting = false
		if err != nil {
			t.Fail()
		} else {
			t.Logf("done: %v", len(data))
		}
	})
	op.Start()
	for isExecuting {

	}
}
