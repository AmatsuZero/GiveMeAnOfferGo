package 剑指Offer

/*
	题目：输入一个链表的头结点，从尾到头反过来打印出每个结点的值。
*/
func (node *ListNode) TraverseReversely(block func(val int)) {
	TraverseReversely(node, block)
}

func TraverseReversely(pHead *ListNode, block func(val int)) {
	if pHead == nil {
		return
	}
	if pHead.Next != nil {
		TraverseReversely(pHead.Next, block)
	}
	block(pHead.Val)
}
