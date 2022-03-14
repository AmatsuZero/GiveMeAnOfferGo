package binary_tree_preorder_traversal

import "GiveMeAnOffer/defines"

// PreorderTraversal https://leetcode-cn.com/problems/binary-tree-preorder-traversal/
func PreorderTraversal(root *defines.TreeNode) []int {
	if root == nil {
		return []int{}
	}
	var res []int
	stack := []*defines.TreeNode{root}
	for len(stack) != 0 {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if node != nil {
			res = append(res, node.Val)
		}
		if node.Right != nil {
			stack = append(stack, node.Right)
		}
		if node.Left != nil {
			stack = append(stack, node.Left)
		}
	}
	return res
}
