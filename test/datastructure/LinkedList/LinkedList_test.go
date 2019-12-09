package LinkedList

import (
	"GiveMeAnOfferGo/Collections"
	"GiveMeAnOfferGo/Objects"
	"fmt"
	"testing"
)

func TestCreateNode(t *testing.T) {
	node1 := Collections.NewNode(Objects.NewNumberWithInt(1), nil)
	node2 := Collections.NewNode(Objects.NewNumberWithInt(2), nil)
	node3 := Collections.NewNode(Objects.NewNumberWithInt(3), nil)
	node1.Next = node2
	node2.Next = node3
	fmt.Print(node1)
}

func TestPushNode(t *testing.T) {
	list := new(Collections.LinkedList)
	list.Push(Objects.NewNumberWithInt(3))
	list.Push(Objects.NewNumberWithInt(2))
	list.Push(Objects.NewNumberWithInt(1))
	fmt.Print(list)
}

func TestAppendNode(t *testing.T) {
	list := new(Collections.LinkedList)
	list.Append(Objects.NewNumberWithInt(1))
	list.Append(Objects.NewNumberWithInt(2))
	list.Append(Objects.NewNumberWithInt(3))
	fmt.Print(list)
}

func TestInsertNode(t *testing.T) {
	list := new(Collections.LinkedList)
	list.Push(Objects.NewNumberWithInt(3))
	list.Push(Objects.NewNumberWithInt(2))
	list.Push(Objects.NewNumberWithInt(1))

	fmt.Printf("Before Inserting: %v\n", list)
	middleNode := list.NodeAt(1)
	for i := 1; i < 4; i++ {
		middleNode = list.Insert(Objects.NewNumberWithInt(-1), middleNode)
	}
	fmt.Printf("After Inserting: %v\n", list)
}

func TestPopNode(t *testing.T) {
	list := new(Collections.LinkedList)
	list.Push(Objects.NewNumberWithInt(3))
	list.Push(Objects.NewNumberWithInt(2))
	list.Push(Objects.NewNumberWithInt(1))

	fmt.Printf("Before poping list: %v\n", list)
	poppedValue := list.Pop()
	fmt.Printf("After popping list: %v\n", list)
	fmt.Printf("Popped Value: %v\n", poppedValue)
}

func TestRemoveLastNode(t *testing.T) {
	list := new(Collections.LinkedList)
	list.Push(Objects.NewNumberWithInt(3))
	list.Push(Objects.NewNumberWithInt(2))
	list.Push(Objects.NewNumberWithInt(1))

	fmt.Printf("Before removing last node: %v\n", list)
	removedValue := list.RemoveLast()

	fmt.Printf("After removing last node: %v\n", list)
	fmt.Printf("Removed Value: %v\n", removedValue)
}

func TestRemoveAfterNode(t *testing.T) {
	list := new(Collections.LinkedList)
	list.Push(Objects.NewNumberWithInt(3))
	list.Push(Objects.NewNumberWithInt(2))
	list.Push(Objects.NewNumberWithInt(1))
	fmt.Printf("Before removing at particular index: %v\n", list)

	index := 1
	node := list.NodeAt(index - 1)
	removedValue := list.RemoveAfter(node)
	fmt.Printf("After removing at index %v: %v \n", index, list)
	fmt.Printf("Removed Value: %v\n", removedValue)
}

func TestCOW(t *testing.T) {
	list1 := new(Collections.LinkedList)
	list1.Append(Objects.NewNumberWithInt(1))
	list1.Append(Objects.NewNumberWithInt(2))

	list2 := list1.Copy()
	fmt.Printf("List1 :%v\n", list1)
	fmt.Printf("List2 :%v\n", list2)

	fmt.Println("After appending 3 to list2")
	list2.Append(Objects.NewNumberWithInt(3))
	fmt.Printf("List1 :%v\n", list1)
	fmt.Printf("List2 :%v\n", list2)
}

func TestIsSameObject(t *testing.T) {
	list1 := new(Collections.LinkedList)
	list1.Append(Objects.NewNumberWithInt(1))
	list1.Append(Objects.NewNumberWithInt(2))

	fmt.Printf("list1: %p\n", list1.Head)

	_ = list1.Copy()
	fmt.Printf("list1: %p\n", list1.Head)
}

func TestTraverse(t *testing.T) {
	list1 := new(Collections.LinkedList)
	for i := 1; i < 11; i++ {
		val := Objects.NewNumberWithInt(i)
		list1.Append(val)
	}
	list1.Traverse(func(val Objects.Comparable) {
		fmt.Println(val)
	})
}

func TestReverseTraverse(t *testing.T) {
	list1 := new(Collections.LinkedList)
	for i := 1; i < 1001; i++ {
		val := Objects.NewNumberWithInt(i)
		list1.Append(val)
	}
	list1.ReverseTraverse(func(val Objects.Comparable) {
		fmt.Println(val)
	})
}

func TestMiddleNode(t *testing.T) {
	list1 := new(Collections.LinkedList)
	for i := 1; i < 4; i++ {
		val := Objects.NewNumberWithInt(i)
		list1.Append(val)
	}

	val := Objects.NewNumberWithInt(2)
	if list1.MiddleValue().Compare(val) != Objects.OrderedSame {
		t.Fail()
	}
	list1.Append(Objects.NewNumberWithInt(4))
	val = Objects.NewNumberWithInt(3)
	if list1.MiddleValue().Compare(val) != Objects.OrderedSame {
		t.Fail()
	}
}

func TestUnique(t *testing.T) {
	list1 := new(Collections.LinkedList)
	list1.Append(Objects.NewNumberWithInt(1))
	list1.Append(Objects.NewNumberWithInt(3))
	list1.Append(Objects.NewNumberWithInt(3))
	list1.Append(Objects.NewNumberWithInt(3))
	list1.Append(Objects.NewNumberWithInt(4))

	list2 := list1.Unique()
	fmt.Println(list2)
}

func TestConvenientInit(t *testing.T) {
	s := Collections.NewStack(
		Objects.NewNumberWithInt(1),
		Objects.NewNumberWithInt(2),
		Objects.NewNumberWithInt(3),
		Objects.NewNumberWithInt(4))

	fmt.Println(s)
}
