package 剑指Offer

/*
题目：输入两个递增排序的链表，合并这两个链表并使新链表中的结点仍然是按照递增排序的。
*/

func (node *ListNode) Merge(head *ListNode) *ListNode {
	return Merge(node, head)
}

func Merge(pHead1, pHead2 *ListNode) (mergedHead *ListNode) {
	if pHead1 == nil {
		return pHead2
	} else if pHead2 == nil {
		return pHead1
	}
	if pHead1.Val < pHead2.Val {
		mergedHead = pHead1
		mergedHead.Next = Merge(pHead1.Next, pHead2)
	} else {
		mergedHead = pHead2
		mergedHead.Next = Merge(pHead1, pHead2.Next)
	}
	return
}
