package maximum_depth_of_binary_tree

import "GiveMeAnOffer/defines"

// MaxDepth https://leetcode-cn.com/problems/maximum-depth-of-binary-tree/
func MaxDepth(root *defines.TreeNode) int {
	if root == nil {
		return 0
	}
	depth := MaxDepth(root.Left)
	if rDepth := MaxDepth(root.Right); rDepth > depth {
		depth = rDepth
	}
	return depth + 1
}
