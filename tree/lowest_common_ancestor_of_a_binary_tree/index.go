package lowest_common_ancestor_of_a_binary_tree

import "GiveMeAnOffer/defines"

// LowestCommonAncestor https://leetcode-cn.com/problems/lowest-common-ancestor-of-a-binary-tree/
func LowestCommonAncestor(root, p, q *defines.TreeNode) *defines.TreeNode {
	if root == nil || root == q || root == p {
		return root
	}
	left := LowestCommonAncestor(root.Left, p, q)
	right := LowestCommonAncestor(root.Right, p, q)
	if left != nil {
		if right != nil {
			return root
		}
		return left
	}
	return right
}
