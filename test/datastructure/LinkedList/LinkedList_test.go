package LinkedList

import (
	"GiveMeAnOfferGo/linkedlist"
	"fmt"
	"testing"
)

func TestCreateNode(t *testing.T) {
	node1 := linkedlist.NewNode(linkedlist.NewNumberWithInt(1), nil)
	node2 := linkedlist.NewNode(linkedlist.NewNumberWithInt(2), nil)
	node3 := linkedlist.NewNode(linkedlist.NewNumberWithInt(3), nil)
	node1.Next = node2
	node2.Next = node3
	fmt.Print(node1)
}

func TestPushNode(t *testing.T) {
	list := new(linkedlist.LinkedList)
	list.Push(linkedlist.NewNumberWithInt(3))
	list.Push(linkedlist.NewNumberWithInt(2))
	list.Push(linkedlist.NewNumberWithInt(1))
	fmt.Print(list)
}

func TestAppendNode(t *testing.T) {
	list := new(linkedlist.LinkedList)
	list.Append(linkedlist.NewNumberWithInt(1))
	list.Append(linkedlist.NewNumberWithInt(2))
	list.Append(linkedlist.NewNumberWithInt(3))
	fmt.Print(list)
}

func TestInsertNode(t *testing.T) {
	list := new(linkedlist.LinkedList)
	list.Push(linkedlist.NewNumberWithInt(3))
	list.Push(linkedlist.NewNumberWithInt(2))
	list.Push(linkedlist.NewNumberWithInt(1))

	fmt.Printf("Before Inserting: %v\n", list)
	middleNode := list.NodeAt(1)
	for i := 1; i < 4; i++ {
		middleNode = list.Insert(linkedlist.NewNumberWithInt(-1), middleNode)
	}
	fmt.Printf("After Inserting: %v\n", list)
}

func TestPopNode(t *testing.T) {
	list := new(linkedlist.LinkedList)
	list.Push(linkedlist.NewNumberWithInt(3))
	list.Push(linkedlist.NewNumberWithInt(2))
	list.Push(linkedlist.NewNumberWithInt(1))

	fmt.Printf("Before poping list: %v\n", list)
	poppedValue := list.Pop()
	fmt.Printf("After popping list: %v\n", list)
	fmt.Printf("Popped Value: %v\n", poppedValue)
}

func TestRemoveLastNode(t *testing.T) {
	list := new(linkedlist.LinkedList)
	list.Push(linkedlist.NewNumberWithInt(3))
	list.Push(linkedlist.NewNumberWithInt(2))
	list.Push(linkedlist.NewNumberWithInt(1))

	fmt.Printf("Before removing last node: %v\n", list)
	removedValue := list.RemoveLast()

	fmt.Printf("After removing last node: %v\n", list)
	fmt.Printf("Removed Value: %v\n", removedValue)
}

func TestRemoveAfterNode(t *testing.T) {
	list := new(linkedlist.LinkedList)
	list.Push(linkedlist.NewNumberWithInt(3))
	list.Push(linkedlist.NewNumberWithInt(2))
	list.Push(linkedlist.NewNumberWithInt(1))
	fmt.Printf("Before removing at particular index: %v\n", list)

	index := 1
	node := list.NodeAt(index - 1)
	removedValue := list.RemoveAfter(node)
	fmt.Printf("After removing at index %v: %v \n", index, list)
	fmt.Printf("Removed Value: %v\n", removedValue)
}

func TestCOW(t *testing.T) {
	list1 := new(linkedlist.LinkedList)
	list1.Append(linkedlist.NewNumberWithInt(1))
	list1.Append(linkedlist.NewNumberWithInt(2))

	list2 := list1.Copy()
	fmt.Printf("List1 :%v\n", list1)
	fmt.Printf("List2 :%v\n", list2)

	fmt.Println("After appending 3 to list2")
	list2.Append(linkedlist.NewNumberWithInt(3))
	fmt.Printf("List1 :%v\n", list1)
	fmt.Printf("List2 :%v\n", list2)
}

func TestIsSameObject(t *testing.T) {
	list1 := new(linkedlist.LinkedList)
	list1.Append(linkedlist.NewNumberWithInt(1))
	list1.Append(linkedlist.NewNumberWithInt(2))

	fmt.Printf("list1: %p\n", list1.Head)

	_ = list1.Copy()
	fmt.Printf("list1: %p\n", list1.Head)
}

func TestTraverse(t *testing.T) {
	list1 := new(linkedlist.LinkedList)
	for i := 1; i < 11; i++ {
		val := linkedlist.NewNumberWithInt(i)
		list1.Append(val)
	}
	list1.Traverse(func(val linkedlist.Comparable) {
		fmt.Println(val)
	})
}

func TestReverseTraverse(t *testing.T) {
	list1 := new(linkedlist.LinkedList)
	for i := 1; i < 1001; i++ {
		val := linkedlist.NewNumberWithInt(i)
		list1.Append(val)
	}
	list1.ReverseTraverse(func(val linkedlist.Comparable) {
		fmt.Println(val)
	})
}

func TestMiddleNode(t *testing.T) {
	list1 := new(linkedlist.LinkedList)
	for i := 1; i < 4; i++ {
		val := linkedlist.NewNumberWithInt(i)
		list1.Append(val)
	}

	val := linkedlist.NewNumberWithInt(2)
	if list1.MiddleValue().Compare(val) != linkedlist.OrderedSame {
		t.Fail()
	}
	list1.Append(linkedlist.NewNumberWithInt(4))
	val = linkedlist.NewNumberWithInt(3)
	if list1.MiddleValue().Compare(val) != linkedlist.OrderedSame {
		t.Fail()
	}
}
