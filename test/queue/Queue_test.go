package queue

import (
	"GiveMeAnOfferGo/Collections"
	"GiveMeAnOfferGo/Objects"
	"fmt"
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

func getString(str string) *Objects.StringObject {
	return &Objects.StringObject{GoString: &str}
}
