package linkedlist

import "fmt"

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
	list.Head = NewNode(value, list.Head)
	if list.Tail == nil {
		list.Tail = list.Head
	}
}

func (list *LinkedList) Append(value Comparable) {
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
	if list.Tail == after {
		list.Append(value)
		return list.Tail
	}
	after.Next = NewNode(value, after.Next)
	return after.Next
}

func (list *LinkedList) Pop() Comparable {
	defer func() {
		list.Head = list.Head.Next
		if list.IsEmpty() {
			list.Tail = nil
		}
	}()
	return list.Head.Value
}

func (list *LinkedList) RemoveLast() Comparable {
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
	defer func() {
		if node.Next == list.Tail {
			list.Tail = node
		}
		node.Next = node.Next.Next
	}()

	return node.Next.Value
}
