package Core

import (
	"context"
	"github.com/reactivex/rxgo/v2"
)

type WebImageOperation interface {
	Cancel()
}

type webImageOperation struct {
	cancel rxgo.Disposable
	ctx    context.Context
}

func (op *webImageOperation) Cancel() {
	op.cancel()
}

func newWebImageOperation(task rxgo.Observable) *webImageOperation {
	ctx, cancel := task.Connect()
	return &webImageOperation{
		cancel: cancel,
		ctx:    ctx,
	}
}
