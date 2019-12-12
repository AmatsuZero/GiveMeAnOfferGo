package Collections

import "GiveMeAnOfferGo/Objects"

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
