package Core

import (
	"context"
	"github.com/AmatsuZero/GiveMeAnOfferGo/Collections"
	"github.com/AmatsuZero/GiveMeAnOfferGo/Objects"
	"github.com/reactivex/rxgo/v2"
)

type ImageDownloadQueue struct {
	priorityQueue    *Collections.PriorityQueue
	tasks            chan rxgo.Item
	ob               rxgo.Observable
	cancel           rxgo.Disposable
	ctx              context.Context
	cancelTokens     map[ImageDownloadOperationProtocol]interface{}
	maxConcurrentNum int
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
		ctx:              ctx,
		cancel:           cancel,
		maxConcurrentNum: maxConcurrentNum,
	}
}

func (queue *ImageDownloadQueue) CancelAllOperations() {
	if queue == nil {
		return
	}
	for !queue.priorityQueue.IsEmpty() {
		task := queue.priorityQueue.Dequeue().(ImageDownloadOperationProtocol)
		token, ok := queue.cancelTokens[task]
		if ok {
			task.CancelWithToken(token)
		}
	}
	queue.cancel()
}

func (queue *ImageDownloadQueue) GetOperations() (ops []ImageDownloadOperationProtocol) {
	if queue == nil || len(queue.cancelTokens) == 0 {
		return ops
	}
	ops = make([]ImageDownloadOperationProtocol, 0, len(queue.cancelTokens))
	for k := range queue.cancelTokens {
		ops = append(ops, k)
	}
	return
}

func (queue *ImageDownloadQueue) AddOperation(op ImageDownloadOperationProtocol) {
	if queue == nil || op == nil {
		return
	}
	queue.priorityQueue.Enqueue(op)
	task := queue.priorityQueue.Dequeue().(ImageDownloadOperationProtocol)
	// 任务完成时，移除任务
	token := task.AddHandlersForProgressAndCompletion(nil, func(data []byte, err error, finished bool) {
		delete(queue.cancelTokens, task)
	})
	queue.cancelTokens[task] = token
	go func(task ImageDownloadOperationProtocol) {
		queue.tasks <- rxgo.Of(task)
	}(task)
}
