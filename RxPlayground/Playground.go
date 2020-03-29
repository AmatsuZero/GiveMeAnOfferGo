package RxPlayground

import (
	"fmt"
	"github.com/reactivex/rxgo/v2"
)

func init() {
	observable := rxgo.Just("Hello, World!")()
	ch := observable.Observe()
	item := <-ch
	fmt.Println(item.V)
}