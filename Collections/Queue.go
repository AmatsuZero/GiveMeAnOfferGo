package Collections

import (
	"GiveMeAnOfferGo/Objects"
	"fmt"
)

type Queue interface {
	Enqueue(element Objects.ObjectProtocol) bool
	Dequeue() Objects.ObjectProtocol
	IsEmpty() bool
	Peek() Objects.ObjectProtocol
	Length() int
}

type QueueArray struct {
	array []Objects.ObjectProtocol
}

func NewQueueArray() *QueueArray {
	return &QueueArray{array: make([]Objects.ObjectProtocol, 0)}
}

func (qa *QueueArray) IsEmpty() bool {
	return qa.Length() == 0
}

func (qa *QueueArray) Enqueue(element Objects.ObjectProtocol) bool {
	qa.array = append(qa.array, element)
	return true
}

func (qa *QueueArray) Dequeue() (x Objects.ObjectProtocol) {
	if qa.IsEmpty() {
		return nil
	}
	x, qa.array = qa.array[0], qa.array[1:]
	return
}

func (qa *QueueArray) Peek() Objects.ObjectProtocol {
	if qa.IsEmpty() {
		return nil
	}
	return qa.array[0]
}

func (qa *QueueArray) Length() int {
	return len(qa.array)
}

type QueueLinkedList struct {
	list *DoublyLinkedList
}

func NewQueueLinedList() *QueueLinkedList {
	return &QueueLinkedList{list: new(DoublyLinkedList)}
}

func (qd *QueueLinkedList) Enqueue(val Objects.ObjectProtocol) bool {
	qd.list.Append(val)
	return true
}

func (qd *QueueLinkedList) Dequeue() Objects.ObjectProtocol {
	element := qd.list.First()
	if qd.list.IsEmpty() || element == nil {
		return nil
	}
	return qd.list.Remove(element)
}

func (qd *QueueLinkedList) Peek() Objects.ObjectProtocol {
	first := qd.list.First()
	if first == nil {
		return nil
	}
	return first.Value
}

func (qd *QueueLinkedList) IsEmpty() bool {
	return qd.list.IsEmpty()
}

func (qd *QueueLinkedList) String() string {
	return fmt.Sprint(qd.list)
}

func (qd *QueueLinkedList) Length() int {
	return qd.list.Length()
}
