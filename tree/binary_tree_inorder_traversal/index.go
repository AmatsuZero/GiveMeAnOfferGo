package binary_tree_inorder_traversal

import "GiveMeAnOffer/defines"

// InorderTraversal https://leetcode-cn.com/problems/binary-tree-inorder-traversal/
func InorderTraversal(root *defines.TreeNode) []int {
	var output []int
	if root == nil {
		return output
	}
	output = append(output, InorderTraversal(root.Left)...)
	output = append(output, root.Val)
	output = append(output, InorderTraversal(root.Right)...)
	return output
}
