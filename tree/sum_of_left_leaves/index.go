package sum_of_left_leaves

import "GiveMeAnOffer/defines"

func SumOfLeftLeaves(root *defines.TreeNode) (ans int) {
	if root == nil {
		return 0
	}
	if root.Left != nil {
		if isLeafNode(root.Left) {
			ans += root.Left.Val
		} else {
			ans += SumOfLeftLeaves(root.Left)
		}
	}
	if root.Right != nil && !isLeafNode(root.Right) {
		ans += SumOfLeftLeaves(root.Right)
	}
	return
}

func isLeafNode(node *defines.TreeNode) bool {
	return node.Left == nil && node.Right == nil
}
