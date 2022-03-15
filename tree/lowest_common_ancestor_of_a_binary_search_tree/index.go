package lowest_common_ancestor_of_a_binary_search_tree

import "GiveMeAnOffer/defines"

// LowestCommonAncestor https://leetcode-cn.com/problems/lowest-common-ancestor-of-a-binary-search-tree/
func LowestCommonAncestor(root, p, q *defines.TreeNode) *defines.TreeNode {
	if p == nil || q == nil || root == nil {
		return nil
	}
	if p.Val < root.Val && q.Val < root.Val {
		return LowestCommonAncestor(root.Left, p, q)
	}
	if p.Val > root.Val && q.Val > root.Val {
		return LowestCommonAncestor(root.Right, p, q)
	}
	return root
}
