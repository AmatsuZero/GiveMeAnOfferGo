package Collections

import (
	"GiveMeAnOfferGo/Objects"
)

type Queue interface {
	Enqueue(element Objects.Comparable) bool
	Dequeue() Objects.Comparable
	IsEmpty() bool
	Peek() Objects.Comparable
}

type QueueArray struct {
	array []Objects.Comparable
}

func (qa *QueueArray) IsEmpty() bool {
	return len(qa.array) == 0
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
