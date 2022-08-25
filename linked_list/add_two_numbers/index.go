package add_two_numbers

import "GiveMeAnOffer/defines"

func AddTwoNumbers(l1 *defines.ListNode, l2 *defines.ListNode) *defines.ListNode {
	l1, l2 = reverse(l1), reverse(l2)
	dummy := &defines.ListNode{}
	pre := dummy
	for h1, h2 := l1, l2; h1 != nil && h2 != nil; {
		var tmp *defines.ListNode
		for sum := h1.Val + h2.Val; sum/10 != 0; sum /= 10 {
			tmp = &defines.ListNode{Val: sum % 10}
		}
		pre.Next = tmp
		pre = tmp
		h1 = h1.Next
		h2 = h2.Next
	}

	return dummy.Next
}

func reverse(head *defines.ListNode) *defines.ListNode {
	if head == nil || head.Next == nil {
		return head
	}

	newHead := reverse(head.Next)
	head.Next.Next = head
	head.Next = nil
	return newHead
}
