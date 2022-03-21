package subtree_of_another_tree

import (
	"GiveMeAnOffer/defines"
	"GiveMeAnOffer/tree/same_tree"
)

// IsSubtree https://leetcode-cn.com/problems/subtree-of-another-tree/
func IsSubtree(s, t *defines.TreeNode) bool {
	if same_tree.IsSameTree(s, t) {
		return true
	}
	if s == nil {
		return false
	}
	if IsSubtree(s.Left, t) || IsSubtree(s.Right, t) {
		return true
	}
	return false
}
