package symmetric_tree

import (
	"GiveMeAnOffer/defines"
	invert_binary_tree "GiveMeAnOffer/tree/invert_binary_tree"
	same_tree "GiveMeAnOffer/tree/same_tree"
)

// IsSymmetric https://leetcode-cn.com/leetbook/read/leetcode-cookbook/5gt176/
func IsSymmetric(root *defines.TreeNode) bool {
	if root == nil {
		return true
	}
	return same_tree.IsSameTree(invert_binary_tree.InvertTree(root.Left), root.Right)
}
