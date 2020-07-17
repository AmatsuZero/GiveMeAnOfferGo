package Core

import (
	"context"
	"github.com/reactivex/rxgo/v2"
	"testing"
)

func TestMoveDir(t *testing.T) {
	result := false
	isExecuting := true
	ob := rxgo.Create([]rxgo.Producer{func(ctx context.Context, next chan<- rxgo.Item) {
		result = true
		isExecuting = false
		next <- rxgo.Of(true)
	}}).First()
	ob.Run()
	t.Log("None blocking")
	for isExecuting {

	}
	t.Log(result)
}
