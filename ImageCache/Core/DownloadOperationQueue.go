package Core

import (
	"context"
	"github.com/AmatsuZero/GiveMeAnOfferGo/Collections"
	"github.com/AmatsuZero/GiveMeAnOfferGo/Objects"
	"github.com/reactivex/rxgo/v2"
	"sync"
)

type ImageDownloadQueue struct {
	priorityQueue    *Collections.PriorityQueue
	tasks            chan rxgo.Item
	ob               rxgo.Observable
	cancel           rxgo.Disposable
	Context          context.Context
	cancelTokens     map[ImageDownloadOperationProtocol]interface{}
	maxConcurrentNum int
	queueLock        sync.Mutex
}

func NewImageDownloadQueue(maxConcurrentNum int) *ImageDownloadQueue {
	tasks := make(chan rxgo.Item)
	ob := rxgo.FromChannel(tasks, rxgo.WithBufferedChannel(maxConcurrentNum))
	ob.DoOnNext(func(i interface{}) {
		task := i.(ImageDownloadOperationProtocol)
		task.Start()
	})
	ctx, cancel := ob.Connect()
	return &ImageDownloadQueue{
		priorityQueue: Collections.NewPriorityQueue(func(lhs Objects.ComparableObject, rhs Objects.ComparableObject) bool {
			return lhs.Compare(rhs) == Objects.OrderedDescending
		}),
		tasks:            tasks,
		ob:               ob,
		cancelTokens:     map[ImageDownloadOperationProtocol]interface{}{},
		Context:          ctx,
		cancel:           cancel,
		maxConcurrentNum: maxConcurrentNum,
	}
}

func (queue *ImageDownloadQueue) CancelAllOperations() {
	if queue == nil {
		return
	}
	queue.queueLock.Lock()
	for !queue.priorityQueue.IsEmpty() {
		task := queue.priorityQueue.Dequeue().(ImageDownloadOperationProtocol)
		token, ok := queue.cancelTokens[task]
		if ok {
			task.CancelWithToken(token)
		}
	}
	queue.cancelTokens = map[ImageDownloadOperationProtocol]interface{}{} // 清空
	queue.queueLock.Unlock()
	queue.cancel()
}

func (queue *ImageDownloadQueue) GetOperations() (ops []ImageDownloadOperationProtocol) {
	if queue == nil || len(queue.cancelTokens) == 0 {
		return ops
	}
	ops = make([]ImageDownloadOperationProtocol, 0, len(queue.cancelTokens))
	queue.queueLock.Lock()
	defer queue.queueLock.Unlock()
	for k := range queue.cancelTokens {
		ops = append(ops, k)
	}
	return
}

func (queue *ImageDownloadQueue) AddOperation(op ImageDownloadOperationProtocol) {
	if queue == nil || op == nil {
		return
	}
	queue.queueLock.Lock()
	defer queue.queueLock.Unlock()
	token := op.AddHandlersForProgressAndCompletion(nil, func(data []byte, err error, finished bool) {
		delete(queue.cancelTokens, op) // 任务完成时，移除任务
	})
	queue.cancelTokens[op] = token
	queue.priorityQueue.Enqueue(op)
	go func() {
		task := queue.priorityQueue.Dequeue()
		if task == nil {
			return
		}
		queue.tasks <- rxgo.Of(task)
	}()
}
