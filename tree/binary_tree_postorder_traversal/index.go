package binary_tree_postorder_traversal

import "GiveMeAnOffer/defines"

func PostorderTraversal(root *defines.TreeNode) (res []int) {
	var stack []*defines.TreeNode
	var prev *defines.TreeNode
	for root != nil || len(stack) > 0 {
		for root != nil {
			stack = append(stack, root)
			root = root.Left
		}
		root = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if root.Right == nil || root.Right == prev {
			res = append(res, root.Val)
			prev = root
			root = nil
		} else {
			stack = append(stack, root)
			root = root.Right
		}
	}
	return
}
