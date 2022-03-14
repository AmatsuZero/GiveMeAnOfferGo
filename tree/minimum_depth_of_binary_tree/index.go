package minimum_depth_of_binary_tree

import "GiveMeAnOffer/defines"

func MinDepth(root *defines.TreeNode) int {
	if root == nil {
		return 0
	}
	if root.Left == nil {
		return MinDepth(root.Right) + 1
	}
	if root.Right == nil {
		return MinDepth(root.Left) + 1
	}
	depth := MinDepth(root.Left)
	if rDepth := MinDepth(root.Right); rDepth < depth {
		depth = rDepth
	}
	return depth + 1
}
