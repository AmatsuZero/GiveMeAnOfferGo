package 剑指Offer

/*
题目：定义一个函数，输入一个链表的头结点，反转该链表并输出反转后链表的头结点
*/
func (node *ListNode) Reverse() (reversedHead *ListNode) {
	if node == nil {
		return
	}
	var prev *ListNode = nil
	pNode := node
	for pNode != nil {
		pNext := pNode.Next
		if pNext == nil { // 原来的尾节点，变成了头节点
			reversedHead = pNode
		}
		pNode.Next = prev
		prev = pNode
		pNode = pNext
	}
	return
}
