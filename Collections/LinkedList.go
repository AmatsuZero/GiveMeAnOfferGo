package Collections

import (
	"fmt"
	"github.com/AmatsuZero/GiveMeAnOfferGo/Objects"
)

type Node struct {
	Value Objects.ComparableObject
	Next  *Node
}

func NewNode(value Objects.ComparableObject, next *Node) *Node {
	return &Node{
		Value: value,
		Next:  next,
	}
}

func (node *Node) String() string {
	if node.Next == nil {
		return fmt.Sprint(node.Value)
	}
	return fmt.Sprintf("%v ComparableObject> %v ", node.Value, node.Next)
}

func (node *Node) IsEqual(target interface{}) bool {
	return node.Value.Compare(target) == Objects.OrderedSame
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

func (list *LinkedList) Push(value Objects.ComparableObject) {
	list.copyNodes()
	list.Head = NewNode(value, list.Head)
	if list.Tail == nil {
		list.Tail = list.Head
	}
}

func (list *LinkedList) Append(value Objects.ComparableObject) {
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

func (list *LinkedList) Insert(value Objects.ComparableObject, after *Node) *Node {
	list.copyNodes()
	if list.Tail == after {
		list.Append(value)
		return list.Tail
	}
	after.Next = NewNode(value, after.Next)
	return after.Next
}

func (list *LinkedList) Pop() Objects.ComparableObject {
	list.copyNodes()
	defer func() {
		list.Head = list.Head.Next
		if list.IsEmpty() {
			list.Tail = nil
		}
	}()
	return list.Head.Value
}

func (list *LinkedList) RemoveLast() Objects.ComparableObject {
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

func (list *LinkedList) RemoveAfter(node *Node) Objects.ComparableObject {
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

func (list *LinkedList) ForEach(block func(index int, val Objects.ComparableObject)) {
	if block == nil || list.IsEmpty() {
		return
	}
	current := list.Head
	i := 0
	for current != nil {
		block(i, current.Value)
		current = current.Next
		i++
	}
}

func (list *LinkedList) ReverseTraverse(block func(val Objects.ComparableObject)) {
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
	list.ForEach(func(i int, val Objects.ComparableObject) {
		length += 1
	})
	return
}

func (list *LinkedList) MiddleValue() Objects.ComparableObject {
	if list.IsEmpty() {
		return nil
	}
	index := list.Length() / 2
	return list.NodeAt(index).Value
}

func (list *LinkedList) Last() Objects.ComparableObject {
	if list.IsEmpty() {
		return nil
	}
	return list.Tail.Value
}

func (list *LinkedList) Unique() (unique *LinkedList) {
	if list.IsEmpty() {
		return nil
	}
	unique = new(LinkedList)
	list.ForEach(func(i int, val Objects.ComparableObject) {
		if unique.IsEmpty() {
			unique.Push(val)
		} else if !unique.Last().IsEqualTo(val) {
			unique.Append(val)
		}
	})
	return
}
