package Core

import (
	"context"
	"github.com/reactivex/rxgo/v2"
	"testing"
)

func TestMoveDir(t *testing.T) {
	var cb ImageNoParamsBlock
	isExecuting := true
	t.Log("out block")
	rxgo.Concat([]rxgo.Observable{rxgo.Empty(), rxgo.Create([]rxgo.Producer{func(ctx context.Context, next chan<- rxgo.Item) {
		t.Log("in block")
		isExecuting = false
		next <- rxgo.Of(true)
	}})}).DoOnCompleted(rxgo.CompletedFunc(cb))
	for isExecuting {

	}
	t.Log("end")
}
