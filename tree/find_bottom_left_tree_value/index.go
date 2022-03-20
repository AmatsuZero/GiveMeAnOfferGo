package find_bottom_left_tree_value

import "GiveMeAnOffer/defines"

// FindBottomLeftValue https://leetcode-cn.com/problems/find-bottom-left-tree-value/
func FindBottomLeftValue(root *defines.TreeNode) int {
	var queue []*defines.TreeNode
	for len(queue) > 0 {
		var next []*defines.TreeNode
		for _, node := range queue {
			if node.Left != nil {
				next = append(next, node.Left)
			}
			if node.Right != nil {
				next = append(next, node.Right)
			}
		}
		if len(next) == 0 {
			return queue[0].Val
		}
		queue = next
	}
	return 0
}
