package Core

import "github.com/reactivex/rxgo/v2"

type WebImageOperation interface {
	Cancel()
}

type webImageOperation struct {
	cancel rxgo.Disposable
}

func (op *webImageOperation) Cancel() {
	op.cancel()
}
