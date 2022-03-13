package invert_binary_tree

import "GiveMeAnOffer/defines"

// InvertTree https://leetcode-cn.com/problems/invert-binary-tree/
func InvertTree(root *defines.TreeNode) *defines.TreeNode {
	if root == nil {
		return root
	}
	left := InvertTree(root.Left)
	right := InvertTree(root.Right)
	root.Left = right
	root.Right = left
	return root
}
