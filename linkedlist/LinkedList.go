package linkedlist

import (
	"fmt"
)

type CompareResult int

const InvalidResult = -1

const (
	OrderedAscending CompareResult = iota
	OrderedSame
	OrderedDescending
)

type Comparable interface {
	fmt.Stringer
	Compare(obj interface{}) CompareResult
}

type Node struct {
	Value Comparable
	Next  *Node
}

func NewNode(value Comparable, next *Node) *Node {
	return &Node{
		Value: value,
		Next:  next,
	}
}

func (node *Node) String() string {
	if node.Next == nil {
		return fmt.Sprint(node.Value)
	}
	return fmt.Sprintf("%v -> %v ", node.Value, node.Next)
}

type LinkedList struct {
	Head *Node
	Tail *Node
}

func (list *LinkedList) IsEmpty() bool {
	return list.Head == nil
}

func (list *LinkedList) String() string {
	if list.IsEmpty() {
		return "Empty List"
	}
	return fmt.Sprintf("%v", list.Head)
}

func (list *LinkedList) Push(value Comparable) {
	list.copyNodes()
	list.Head = NewNode(value, list.Head)
	if list.Tail == nil {
		list.Tail = list.Head
	}
}

func (list *LinkedList) Append(value Comparable) {
	list.copyNodes()
	if list.IsEmpty() {
		list.Push(value)
		return
	}

	list.Tail.Next = NewNode(value, nil)
	list.Tail = list.Tail.Next
}

func (list *LinkedList) NodeAt(index int) *Node {
	currentNode := list.Head
	currentIndex := 0

	for currentNode != nil && currentIndex < index {
		currentNode = currentNode.Next
		currentIndex += 1
	}

	return currentNode
}

func (list *LinkedList) Insert(value Comparable, after *Node) *Node {
	list.copyNodes()
	if list.Tail == after {
		list.Append(value)
		return list.Tail
	}
	after.Next = NewNode(value, after.Next)
	return after.Next
}

func (list *LinkedList) Pop() Comparable {
	list.copyNodes()
	defer func() {
		list.Head = list.Head.Next
		if list.IsEmpty() {
			list.Tail = nil
		}
	}()
	return list.Head.Value
}

func (list *LinkedList) RemoveLast() Comparable {
	list.copyNodes()
	if list.Head == nil {
		return nil
	}
	if list.Head.Next == nil {
		return list.Pop()
	}
	prev := list.Head
	current := list.Head

	next := current.Next
	for next != nil {
		prev = current
		current = next
		next = current.Next
	}

	prev.Next = nil
	list.Tail = prev
	return current.Value
}

func (list *LinkedList) RemoveAfter(node *Node) Comparable {
	list.copyNodes()
	defer func() {
		if node.Next == list.Tail {
			list.Tail = node
		}
		node.Next = node.Next.Next
	}()

	return node.Next.Value
}

func (list *LinkedList) copyNodes() {
	oldNode := list.Head
	if oldNode == nil {
		return
	}
	list.Head = NewNode(oldNode.Value, nil)
	newNode := list.Head
	nextOldNode := oldNode.Next
	for nextOldNode != nil {
		newNode.Next = NewNode(nextOldNode.Value, nil)
		newNode = newNode.Next
		oldNode = nextOldNode
		nextOldNode = oldNode.Next
	}
	list.Tail = newNode
}

func (list *LinkedList) Copy() *LinkedList {
	newList := &LinkedList{
		Head: list.Head,
		Tail: list.Tail,
	}
	list.copyNodes()
	return newList
}

func (list *LinkedList) Traverse(block func(val Comparable)) {
	if block == nil || list.IsEmpty() {
		return
	}
	current := list.Head
	for current != nil {
		block(current.Value)
		current = current.Next
	}
}

func (list *LinkedList) ReverseTraverse(block func(val Comparable)) {
	if block == nil || list.IsEmpty() {
		return
	}

	newList := list.Copy()
	val := newList.RemoveLast()
	for val != nil {
		block(val)
		val = newList.RemoveLast()
	}
}

func (list *LinkedList) Length() (length int) {
	list.Traverse(func(val Comparable) {
		length += 1
	})
	return
}

func (list *LinkedList) MiddleValue() Comparable {
	if list.IsEmpty() {
		return nil
	}
	index := list.Length() / 2
	return list.NodeAt(index).Value
}
