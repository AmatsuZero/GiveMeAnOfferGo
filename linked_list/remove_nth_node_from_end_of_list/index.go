package remove_nth_node_from_end_of_list

import "GiveMeAnOffer/defines"

// RemoveNthFromEnd https://leetcode.cn/problems/remove-nth-node-from-end-of-list/
func RemoveNthFromEnd(head *defines.ListNode, n int) *defines.ListNode {
	dummy := &defines.ListNode{Next: head}
	first, second := head, dummy
	for i := 0; i < n; i++ {
		first = first.Next
	}

	for ; first != nil; first = first.Next {
		second = second.Next
	}
	second.Next = second.Next.Next
	return dummy.Next
}
