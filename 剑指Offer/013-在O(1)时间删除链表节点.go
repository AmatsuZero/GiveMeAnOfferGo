package 剑指Offer

/*
给定单向链表的头指针和一个结点指针，定义一个函数在O（1）时间删除该结点
*/
func DeleteNode(pListHead, pToBeDeleted *ListNode) {
	if pListHead == nil || pToBeDeleted == nil {
		return
	}
	if pToBeDeleted.Next != nil { // 要删除的节点不是尾节点
		pNext := pToBeDeleted.Next
		pToBeDeleted.Val = pNext.Val
		pToBeDeleted.Next = pNext.Next
	} else if pListHead != pToBeDeleted {
		pNode := pListHead
		for pNode.Next != pToBeDeleted {
			pNode = pNode.Next
		}
		pNode.Next = nil
	}
	// 隐含的只有一个节点的情况，在 Golang GC 的限制，不再引用即删除
}

func (node *ListNode) DeleteNode(ptr *ListNode) {
	if node == nil {
		return
	}
	DeleteNode(node, ptr)
}
