package Collections

import (
	"GiveMeAnOfferGo/Objects"
	"fmt"
)

type DoublyLinkedListNode struct {
	Next     *DoublyLinkedListNode
	Previous *DoublyLinkedListNode
	Value    Objects.Comparable
}

func NewDoublyLinkedListNode(val Objects.Comparable) *DoublyLinkedListNode {
	return &DoublyLinkedListNode{Value: val}
}

func (node *DoublyLinkedListNode) String() string {
	return fmt.Sprint(node.Value)
}

type DoublyLinkedList struct {
	Head *DoublyLinkedListNode
	Tail *DoublyLinkedListNode
}

func (dl *DoublyLinkedList) First() *DoublyLinkedListNode {
	return dl.Head
}

func (dl *DoublyLinkedList) IsEmpty() bool {
	return dl.Head == nil
}

func (dl *DoublyLinkedList) Append(val Objects.Comparable) {
	newNode := NewDoublyLinkedListNode(val)
	if dl.Tail == nil {
		dl.Head = newNode
		dl.Tail = newNode
		return
	}

	newNode.Previous = dl.Tail
	dl.Tail.Next = newNode
	dl.Tail = newNode
}

func (dl *DoublyLinkedList) Remove(node *DoublyLinkedListNode) Objects.Comparable {
	prev := node.Previous
	next := node.Next

	if prev != nil {
		prev.Next = next
	} else {
		dl.Head = next
	}

	next.Previous = prev

	if next == nil {
		dl.Tail = prev
	}

	node.Previous = nil
	node.Next = nil

	return node.Value
}

func (dl *DoublyLinkedList) String() string {
	current := dl.Head
	str := ""
	for current != nil {
		str += fmt.Sprintf("%v -> ", current.Value)
		current = current.Next
	}

	return str + "end"
}

func (dl *DoublyLinkedList) Length() (length int) {
	current := dl.Head
	for current != nil {
		length += 1
		current = current.Next
	}
	return
}
