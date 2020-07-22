package Core

import (
	"context"
	"fmt"
	"github.com/reactivex/rxgo/v2"
	"testing"
	"time"
)

func TestMoveDir(t *testing.T) {
	isFinished := false
	ch := make(chan rxgo.Item)
	// Create a Connectable Observable
	observable := rxgo.FromChannel(ch, rxgo.WithPublishStrategy())
	disposed, cancel := observable.Connect()
	go func() {
		for i := 0; i < 10; i++ {
			ch <- rxgo.Of(i)
		}
	}()
	// Create the first Observer
	observable.DoOnNext(func(i interface{}) {
		t.Logf("First observer: %d\n", i)
	})

	// Create the second Observer
	observable.DoOnNext(func(i interface{}) {
		t.Logf("Second observer: %d\n", i)
	})

	observable.DoOnCompleted(func() {
		isFinished = true
		cancel()
		_, ok := <-ch
		if !ok {
			t.Log("closed")
		}
	})
	isExecuting := true
	go func() {
		// Do something
		time.Sleep(time.Second)
		// Then cancel the subscription
		cancel()
		isExecuting = false
	}()
	// Wait for the subscription to be disposed
	go func(ctx context.Context, cancel rxgo.Disposable) {
		for {
			select {
			case <-ctx.Done():
				if isFinished {
					t.Log("Done")
					return
				}
				isFinished = true
				t.Log("Cancel")
				return
			}
		}
	}(disposed, cancel)

	for isExecuting {

	}
}

func createRange() []int {
	fmt.Println("Called")
	return []int{1, 2, 3, 4, 5, 6}
}

func TestPerformance(t *testing.T) {
	for _, num := range createRange() {
		t.Log(num)
	}
}
