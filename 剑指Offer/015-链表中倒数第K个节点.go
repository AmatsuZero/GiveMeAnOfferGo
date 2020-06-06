package 剑指Offer

/*
输入一个链表，输出该链表中倒数第k个结点。为了符合大多数人的习惯，本题从1开始计数，即链表的尾结点是倒数第1个结点。
例如一个链表有6个结点，从头结点开始它们的值依次是1、2、3、4、5、6。这个链表的倒数第3个结点是值为4的结点
*/

func (node *ListNode) FindKthToTail(k uint) *ListNode {
	if k == 0 { // 从1开始，输入0，没有意义
		return nil
	}
	head := node
	var behind *ListNode
	// head先前进到k-1
	for i := uint(0); i < k-1; i++ {
		if head.Next != nil {
			head = head.Next
		} else {
			return nil
		}
	}
	// 后面的指针指向头指针
	behind = node
	// 当前面的指针走到尾部的时候，后面的指针正好指向倒数第k个
	for head.Next != nil {
		head = head.Next
		behind = behind.Next
	}
	return behind
}

/*
求链表的中间结点。如果链表中结点总数为奇数，返回中间结点；
如果结点总数是偶数，返回中间两个结点的任意一个
*/
func (node *ListNode) MidNode() *ListNode {
	if node.Next == nil {
		return node
	}
	// 慢指针一次前进一步
	slow := node
	// 快指针一次前进两步
	fast := slow
	// 当快指针走到最后的时候，慢指针正好指向中间
	for fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
		if fast == nil {
			return slow
		}
	}
	return slow
}

/*
是否是环形结构
*/
func (node *ListNode) IsCircleList() bool {
	if node.Next == nil {
		return false
	}
	// 慢指针一次前进一步
	slow := node
	// 快指针一次前进两步
	fast := slow
	// 当快指针走到最后的时候，慢指针正好指向中间
	for fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
		if fast == nil { // 快指针走到了最后，说明不是环形链表
			return false
		} else if fast == slow { // 快慢指针重合了，说明是环形链表
			return true
		}
	}
	return false
}
