package Collections

import "github.com/AmatsuZero/GiveMeAnOfferGo/Objects"

type PriorityQueue struct {
	heap *Heap
}

func NewPriorityQueue(sort func(lhs Objects.ComparableObject, rhs Objects.ComparableObject) bool, elements ...Objects.ComparableObject) *PriorityQueue {
	return &PriorityQueue{
		heap: NewHeap(sort, elements),
	}
}

func (pq *PriorityQueue) IsEmpty() bool {
	return pq.heap.IsEmpty()
}

func (pq *PriorityQueue) Peek() Objects.ComparableObject {
	return pq.heap.Peek()
}

func (pq *PriorityQueue) Enqueue(element Objects.ComparableObject) bool {
	pq.heap.Insert(element)
	return true
}

func (pq *PriorityQueue) Dequeue() Objects.ComparableObject {
	return pq.heap.Remove()
}

func (pq *PriorityQueue) Length() int {
	return pq.heap.Count()
}
