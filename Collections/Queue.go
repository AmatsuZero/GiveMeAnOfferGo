package Collections

import (
	"GiveMeAnOfferGo/Objects"
	"fmt"
)

type Queue interface {
	Enqueue(element Objects.Comparable) bool
	Dequeue() Objects.Comparable
	IsEmpty() bool
	Peek() Objects.Comparable
	Length() int
}

type QueueArray struct {
	array []Objects.Comparable
}

func NewQueueArray() *QueueArray {
	return &QueueArray{array: make([]Objects.Comparable, 0)}
}

func (qa *QueueArray) IsEmpty() bool {
	return qa.Length() == 0
}

func (qa *QueueArray) Enqueue(element Objects.Comparable) bool {
	qa.array = append(qa.array, element)
	return true
}

func (qa *QueueArray) Dequeue() (x Objects.Comparable) {
	if qa.IsEmpty() {
		return nil
	}
	x, qa.array = qa.array[0], qa.array[1:]
	return
}

func (qa *QueueArray) Peek() Objects.Comparable {
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

func (qd *QueueLinkedList) Enqueue(val Objects.Comparable) bool {
	qd.list.Append(val)
	return true
}

func (qd *QueueLinkedList) Dequeue() Objects.Comparable {
	element := qd.list.First()
	if qd.list.IsEmpty() || element == nil {
		return nil
	}
	return qd.list.Remove(element)
}

func (qd *QueueLinkedList) Peek() Objects.Comparable {
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
