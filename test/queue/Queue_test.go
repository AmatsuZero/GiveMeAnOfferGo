package queue

import (
	"fmt"
	"github.com/AmatsuZero/GiveMeAnOfferGo/Collections"
	"github.com/AmatsuZero/GiveMeAnOfferGo/Objects"
	"github.com/AmatsuZero/GiveMeAnOfferGo/test/Utils"
	"testing"
)

func TestQueueDoublyLinkedList(t *testing.T) {
	queue := Collections.NewQueueLinedList()
	queue.Enqueue(getString("Ray"))
	queue.Enqueue(getString("Brain"))
	queue.Enqueue(getString("Eric"))
	fmt.Println(queue)
	queue.Dequeue()
	fmt.Println(queue)
	fmt.Println(queue.Peek())
	if queue.Length() != 2 {
		t.Fail()
	}
}

func TestQueueArray(t *testing.T) {
	queue := Collections.NewQueueArray()
	queue.Enqueue(getString("Ray"))
	queue.Enqueue(getString("Brain"))
	queue.Enqueue(getString("Eric"))
	fmt.Println(queue)
	queue.Dequeue()
	fmt.Println(queue)
	fmt.Println(queue.Peek())
	if queue.Length() != 2 {
		t.Fail()
	}
}

func TestQueueRingBuffer(t *testing.T) {
	queue := Collections.NewQueueRingBuffer(10)
	queue.Enqueue(getString("Ray"))
	queue.Enqueue(getString("Brain"))
	queue.Enqueue(getString("Eric"))
	fmt.Println(queue)
	queue.Dequeue()
	fmt.Println(queue)
	fmt.Println(queue.Length())
	fmt.Println(queue.Peek())
}

func TestQueueStack(t *testing.T) {
	queue := Collections.NewQueueStack()
	queue.Enqueue(getString("Ray"))
	queue.Enqueue(getString("Brain"))
	queue.Enqueue(getString("Eric"))
	fmt.Println(queue)
	queue.Dequeue()
	fmt.Println(queue)
	fmt.Println(queue.Length())
	fmt.Println(queue.Peek())
}

func TestReveredQueueArray(t *testing.T) {
	getInt := Utils.GetInt
	queue := Collections.NewQueueArray()
	queue.Enqueue(getInt(1))
	queue.Enqueue(getInt(21))
	queue.Enqueue(getInt(18))
	queue.Enqueue(getInt(42))

	fmt.Printf("Before :%v\n", queue)
	fmt.Printf("After :%v\n", queue.Reversed())
	queue.Reverse()
	fmt.Printf("Reverse Self: %v\n", queue)
}

func getString(str string) *Objects.StringObject {
	return &Objects.StringObject{GoString: &str}
}
