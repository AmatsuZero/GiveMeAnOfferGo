package linkedlist

import "fmt"

type ListNode struct {
	Val  int
	Next *ListNode
}

func (i *ListNode) Append(val int) {
	if i.Next == nil {
		i.Next = NewNode(val)
	} else {
		i.Next.Append(val)
	}
}

func NewNode(val int) *ListNode {
	return &ListNode{Val: val}
}

func BuildLinkedList(nums []int) (head *ListNode) {
	for _, val := range nums {
		if head == nil {
			head = NewNode(val)
		} else {
			head.Append(val)
		}
	}
	return
}

func (i *ListNode) IsEqual(another *ListNode) bool {
	if another == nil || i.Length() != another.Length() {
		return false
	}
	l1 := i
	l2 := another
	for l1 != nil && l2 != nil {
		if l1.Val != l2.Val {
			return false
		}
		l1 = l1.Next
		l2 = l2.Next
	}
	return true
}

type Iterator = func(index int, val int, stop *bool)

func (i *ListNode) Enumerate(iterator Iterator) {
	head := i
	index := 0
	stop := false
	for head != nil && !stop {
		iterator(index, head.Val, &stop)
		head = head.Next
		index += 1
	}
}

func (i *ListNode) String() (ret string) {
	i.Enumerate(func(index int, val int, stop *bool) {
		ret += fmt.Sprintf(" index: %v value: Val: %v", index, val)
	})
	return
}

func (i *ListNode) Length() (len int) {
	head := i
	for head != nil {
		len += 1
		head = head.Next
	}
	return
}

func (i *ListNode) Copy() (node *ListNode) {
	if i == nil {
		return nil
	}
	i.Enumerate(func(index int, val int, stop *bool) {
		if node == nil {
			node = NewNode(val)
		} else {
			node.Append(val)
		}
	})
	return
}

func (i *ListNode) Contains(target int) (ret bool) {
	i.Enumerate(func(index int, val int, stop *bool) {
		if target == val {
			ret = true
			*stop = true
		}
	})
	return
}

func (i *ListNode) AppendList(node ListNode) {
	node.Enumerate(func(index int, val int, stop *bool) {
		i.Append(val)
	})
}

func (i *ListNode) ToArray() (array []int) {
	i.Enumerate(func(index int, val int, stop *bool) {
		array = append(array, val)
	})
	return
}
