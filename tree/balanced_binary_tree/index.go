package balanced_binary_tree

import (
	"GiveMeAnOffer/defines"
	"GiveMeAnOffer/tree/maximum_depth_of_binary_tree"
)

// IsBalanced https://leetcode-cn.com/problems/balanced-binary-tree/
func IsBalanced(root *defines.TreeNode) bool {
	if root == nil {
		return true
	}
	leftHeight := maximum_depth_of_binary_tree.MaxDepth(root.Left)
	rightHeight := maximum_depth_of_binary_tree.MaxDepth(root.Right)
	diff := leftHeight - rightHeight
	if diff < 0 {
		diff = -diff
	}
	return diff <= 1 && IsBalanced(root.Left) && IsBalanced(root.Right)
}
