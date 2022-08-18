package design_circular_queue

type MyCircularDeque struct {
	front, rear int
	elements    []int
}

func Constructor(k int) MyCircularDeque {
	return MyCircularDeque{elements: make([]int, k+1)}
}

func (d *MyCircularDeque) InsertFront(value int) bool {
	if d.IsFull() {
		return false
	}
	cnt := len(d.elements)
	d.front = (d.front - 1 + cnt) % cnt
	d.elements[d.front] = value
	return true
}

func (d *MyCircularDeque) InsertLast(value int) bool {
	if d.IsFull() {
		return false
	}

	d.elements[d.rear] = value
	cnt := len(d.elements)
	d.rear = (d.rear + 1) % cnt
	return true
}

func (d *MyCircularDeque) DeleteFront() bool {
	if d.IsEmpty() {
		return false
	}
	d.front = (d.front + 1) % len(d.elements)
	return true
}

func (d *MyCircularDeque) DeleteLast() bool {
	if d.IsEmpty() {
		return false
	}
	cnt := len(d.elements)
	d.rear = (d.rear - 1 + cnt) % cnt
	return true
}

func (d *MyCircularDeque) GetFront() int {
	if d.IsEmpty() {
		return -1
	}
	return d.elements[d.front]
}

func (d *MyCircularDeque) GetRear() int {
	if d.IsEmpty() {
		return -1
	}
	return d.elements[(d.rear-1+len(d.elements))%len(d.elements)]
}

func (d *MyCircularDeque) IsEmpty() bool {
	return d.front == d.rear
}

func (d *MyCircularDeque) IsFull() bool {
	return (d.rear+1)%len(d.elements) == d.front
}
