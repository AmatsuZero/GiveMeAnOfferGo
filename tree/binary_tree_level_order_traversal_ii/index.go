package binary_tree_level_order_traversal_ii

import (
	"GiveMeAnOffer/defines"
	"GiveMeAnOffer/tree/binary_tree_level_order_traversal"
)

// LevelOrderBottom https://leetcode-cn.com/problems/binary-tree-level-order-traversal-ii/
func LevelOrderBottom(root *defines.TreeNode) [][]int {
	res := binary_tree_level_order_traversal.LevelOrder(root)
	for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
		res[i], res[j] = res[j], res[i]
	}
	return res
}
