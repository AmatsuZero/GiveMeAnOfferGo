package merge_k_sorted_lists

import "GiveMeAnOffer/defines"

// MergeKLists https://leetcode.cn/problems/merge-k-sorted-lists/
func MergeKLists(lists []*defines.ListNode) *defines.ListNode {
	if len(lists) == 0 {
		return nil
	} else if len(lists) == 1 {
		return lists[1]
	}

	head := lists[0]
	for i := 1; i < len(lists); i++ {
		head = mergeTwo(head, lists[i])
	}
	return head
}

func mergeTwo(lhs, rhs *defines.ListNode) *defines.ListNode {
	head := lhs
	if lhs.Val > rhs.Val {
		head = rhs
	}
	cur := head

	lhs = lhs.Next
	rhs = rhs.Next
	for lhs != nil {
		if lhs.Val < rhs.Val {
			cur.Next = lhs
		} else {
			cur.Next = rhs
		}
		cur = cur.Next
		lhs = lhs.Next
		rhs = rhs.Next
	}

	for lhs != nil {
		cur.Next = lhs
		cur = cur.Next
		lhs = lhs.Next
	}

	for rhs != nil {
		cur.Next = rhs
		cur = cur.Next
		rhs = rhs.Next
	}

	return head
}
