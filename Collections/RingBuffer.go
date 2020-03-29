package Collections

import (
	"fmt"
	"github.com/AmatsuZero/GiveMeAnOfferGo/Objects"
)

type RingBuffer struct {
	array      []Objects.ObjectProtocol
	readIndex  int
	writeIndex int
}

func NewBufferRing(count int) *RingBuffer {
	return &RingBuffer{array: make([]Objects.ObjectProtocol, count)}
}

func (buffer *RingBuffer) Write(element Objects.ObjectProtocol) bool {
	if buffer.IsFull() {
		return false
	}
	buffer.array[buffer.writeIndex%len(buffer.array)] = element
	buffer.writeIndex += 1
	return true
}

func (buffer *RingBuffer) Read() Objects.ObjectProtocol {
	if buffer.IsEmpty() {
		return nil
	}
	element := buffer.array[buffer.readIndex%len(buffer.array)]
	buffer.readIndex += 1
	return element
}

func (buffer *RingBuffer) availableSpaceForReading() int {
	return buffer.writeIndex - buffer.readIndex
}

func (buffer *RingBuffer) IsEmpty() bool {
	return buffer.availableSpaceForReading() == 0
}

func (buffer *RingBuffer) availableSpaceForWriting() int {
	return len(buffer.array) - buffer.availableSpaceForReading()
}

func (buffer *RingBuffer) IsFull() bool {
	return buffer.availableSpaceForWriting() == 0
}

func (buffer *RingBuffer) String() string {
	str := fmt.Sprintf("=== RingBuffer: %p\n", buffer)
	if !buffer.IsEmpty() {
		buffer.ForEach(func(index int, val Objects.ObjectProtocol) {
			str += fmt.Sprintf("[%d]:%v\n", index, val)
		})
	}
	str += "=== end\n"
	return str
}

func (buffer *RingBuffer) ForEach(traverse func(index int, val Objects.ObjectProtocol)) {
	if traverse == nil || buffer.IsEmpty() {
		return
	}
	index := 0
	for i := buffer.readIndex; i < buffer.writeIndex; i++ {
		traverse(index, buffer.array[i])
		index++
	}
}
