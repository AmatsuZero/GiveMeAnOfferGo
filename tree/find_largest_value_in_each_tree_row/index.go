package find_largest_value_in_each_tree_row

import (
	"GiveMeAnOffer/defines"
	"GiveMeAnOffer/tree/binary_tree_level_order_traversal"
	"sort"
)

// LargestValues https://leetcode-cn.com/problems/find-largest-value-in-each-tree-row/
func LargestValues(root *defines.TreeNode) (res []int) {
	tmp := binary_tree_level_order_traversal.LevelOrder(root)
	for _, nums := range tmp {
		sort.Ints(nums)
		res = append(res, nums[len(nums)-1])
	}
	return
}
