package convert_sorted_list_to_binary_search_tree

import "GiveMeAnOffer/defines"

// SortedListToBST https://leetcode-cn.com/problems/convert-sorted-list-to-binary-search-tree/
func SortedListToBST(head *defines.ListNode) *defines.TreeNode {
	return buildTree(head, nil)
}

func getMedian(left, right *defines.ListNode) *defines.ListNode {
	fast, slow := left, left
	for fast != right && fast.Next != right {
		fast = fast.Next.Next
		slow = slow.Next
	}
	return slow
}

func buildTree(left, right *defines.ListNode) *defines.TreeNode {
	if left == right {
		return nil
	}
	mid := getMedian(left, right)
	root := &defines.TreeNode{Val: mid.Val}
	root.Left = buildTree(left, mid)
	root.Right = buildTree(mid.Next, right)
	return root
}
