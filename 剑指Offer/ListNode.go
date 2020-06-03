package 剑指Offer

type ListNode struct {
	Val  int
	Next *ListNode
}

func AddToTail(pHead **ListNode, value int) {
	pNew := &ListNode{
		Val:  value,
		Next: nil,
	}
	/*
		我们要特别注意函数的第一个参数pHead是一个指向指针的指针。当我们往一个空链表中插入一个结点时，新插入的结点就是链表的头指针。
		由于此时会改动头指针，因此必须把pHead参数设为指向指针的指针，否则出了这个函数pHead仍然是一个空指针。
	*/
	if *pHead == nil {
		*pHead = pNew
	} else {
		pNode := *pHead
		for pNode.Next != nil {
			pNode = pNode.Next
		}
		pNode.Next = pNew
	}
}

func (node *ListNode) AddToTail(value int) {
	AddToTail(&node, value)
}

func (node *ListNode) IntArray() []int {
	array := make([]int, 0)
	node.Enumerate(func(value int, stop *bool) {
		array = append(array, value)
	})
	return array
}

func (node *ListNode) Enumerate(block func(value int, stop *bool)) {
	head := node
	flag := false
	for head != nil && !flag {
		block(head.Val, &flag)
		head = head.Next
	}
}

func RemoveNode(pHead **ListNode, value int) (pToBeDeleted *ListNode) {
	if pHead == nil || *pHead == nil {
		return
	}
	if (*pHead).Val == value {
		*pHead = (*pHead).Next
	} else {
		pNode := *pHead
		for pNode.Next != nil && pNode.Next.Val != value {
			pNode = pNode.Next
		}
		if pNode.Next != nil && pNode.Next.Val == value {
			pToBeDeleted = pNode.Next
			pNode.Next = pNode.Next.Next
		}
	}
	return
}

func (node *ListNode) RemoveNode(value int) *ListNode {
	return RemoveNode(&node, value)
}
