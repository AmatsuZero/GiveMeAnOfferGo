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

func (qa *QueueArray) Reversed() (queue *QueueArray) {
	queue = qa.Copy()
	stack := NewStack()
	element := queue.Dequeue()
	for element != nil {
		stack.Push(element)
		element = queue.Dequeue()
	}
	element = stack.Pop()
	for element != nil {
		queue.Enqueue(element)
		element = stack.Pop()
	}
	return
}

func (qa *QueueArray) Reverse() {
	stack := NewStack()
	element := qa.Dequeue()
	for element != nil {
		stack.Push(element)
		element = qa.Dequeue()
	}
	element = stack.Pop()
	for element != nil {
		qa.Enqueue(element)
		element = stack.Pop()
	}
}

func (qa *QueueArray) Copy() *QueueArray {
	if qa.IsEmpty() {
		return NewQueueArray()
	}
	return &QueueArray{array: append(qa.array[:0:0], qa.array...)}
}

func (qa *QueueArray) String() string {
	str := fmt.Sprintf("=== Queue Array: %p\n", qa)
	for i, v := range qa.array {
		str += fmt.Sprintf("[%d] %v\n", i, v)
	}
	str += "=== end"
	return str
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

type QueueRingBuffer struct {
	ringBuffer *RingBuffer
}

func NewQueueRingBuffer(count int) *QueueRingBuffer {
	return &QueueRingBuffer{
		ringBuffer: NewBufferRing(count),
	}
}

func (qr *QueueRingBuffer) IsEmpty() bool {
	return qr.ringBuffer.IsEmpty()
}

func (qr *QueueRingBuffer) Peek() Objects.ObjectProtocol {
	if qr.IsEmpty() {
		return nil
	}
	return qr.ringBuffer.array[0]
}

func (qr *QueueRingBuffer) Enqueue(val Objects.ObjectProtocol) bool {
	return qr.ringBuffer.Write(val)
}

func (qr *QueueRingBuffer) Dequeue() Objects.ObjectProtocol {
	if qr.IsEmpty() {
		return nil
	}
	return qr.ringBuffer.Read()
}

func (qr *QueueRingBuffer) Length() int {
	return qr.ringBuffer.availableSpaceForReading()
}

func (qr *QueueRingBuffer) String() string {
	return fmt.Sprintln(qr.ringBuffer)
}

type QueueStack struct {
	leftStack  []Objects.ObjectProtocol
	rightStack []Objects.ObjectProtocol
}

func NewQueueStack() *QueueStack {
	return &QueueStack{
		leftStack:  make([]Objects.ObjectProtocol, 0),
		rightStack: make([]Objects.ObjectProtocol, 0),
	}
}

func (qs *QueueStack) IsEmpty() bool {
	return len(qs.leftStack) == 0 && len(qs.rightStack) == 0
}

func (qs *QueueStack) Peek() Objects.ObjectProtocol {
	if qs.IsEmpty() {
		return nil
	} else if len(qs.leftStack) > 0 {
		return qs.leftStack[len(qs.leftStack)-1]
	} else {
		return qs.rightStack[0]
	}
}

func (qs *QueueStack) Enqueue(element Objects.ObjectProtocol) bool {
	qs.rightStack = append(qs.rightStack, element)
	return true
}

func (qs *QueueStack) Dequeue() (popped Objects.ObjectProtocol) {
	if len(qs.leftStack) == 0 {
		qs.leftStack = make([]Objects.ObjectProtocol, 0)
		for i := len(qs.rightStack) - 1; i >= 0; i-- {
			qs.leftStack = append(qs.leftStack, qs.rightStack[i])
		}
		qs.rightStack = make([]Objects.ObjectProtocol, 0)
	}
	popped, qs.leftStack = qs.leftStack[len(qs.leftStack)-1], qs.leftStack[:len(qs.leftStack)-1]
	return
}

func (qs *QueueStack) Length() int {
	return len(qs.leftStack) + len(qs.rightStack)
}

func (qs *QueueStack) ForEach(traverse func(index int, val Objects.ObjectProtocol)) {
	if qs.IsEmpty() || traverse == nil {
		return
	}
	index := 0
	for i := len(qs.leftStack) - 1; i >= 0; i-- {
		traverse(index, qs.leftStack[i])
		index++
	}
	for i := 0; i < len(qs.rightStack); i++ {
		traverse(index, qs.rightStack[i])
		index++
	}
}

func (qs *QueueStack) String() string {
	str := fmt.Sprintf("=== Queue Stack: %p\n", qs)
	qs.ForEach(func(index int, val Objects.ObjectProtocol) {
		str += fmt.Sprintf("[%d] %v\n", index, val)
	})
	str += "=== End"
	return str
}
