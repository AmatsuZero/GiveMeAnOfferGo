package LeetCode解题

type MyCircularQueueInterface interface {
	Front() int
	Rear() int
	Enqueue(val int) bool
	Dequeue() bool
	IsEmpty() bool
	IsFull() bool
}

type MyCircularQueue struct {
	backStore  []int
	head, tail int
}

func NewCircularQueue(size int) MyCircularQueue {
	return MyCircularQueue{
		backStore: make([]int, size),
		head:      -1,
		tail:      -1,
	}
}

func (q MyCircularQueue) IsFull() bool {
	return (q.tail+1)%cap(q.backStore) == q.head
}

func (q MyCircularQueue) IsEmpty() bool {
	return q.head == -1
}

func (q MyCircularQueue) Front() int {
	if q.IsEmpty() {
		return -1
	}
	return q.backStore[q.head]
}

func (q MyCircularQueue) Rear() int {
	if q.IsEmpty() {
		return -1
	}
	return q.backStore[q.tail]
}

func (q *MyCircularQueue) Enqueue(val int) bool {
	if q.IsFull() {
		return false
	}
	if q.IsEmpty() {
		q.head = 0
	}
	q.tail = (q.tail + 1) % cap(q.backStore)
	q.backStore[q.tail] = val
	return true
}

func (q *MyCircularQueue) Dequeue() bool {
	if q.IsEmpty() {
		return false
	}
	if q.head == q.tail {
		q.head = -1
		q.tail = -1
		return true
	}
	q.head = (q.head + 1) % cap(q.backStore)
	return true
}
