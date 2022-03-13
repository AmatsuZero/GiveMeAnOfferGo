package binary_tree_level_order_traversal

import "GiveMeAnOffer/defines"

// LevelOrder https://leetcode-cn.com/problems/binary-tree-level-order-traversal/
func LevelOrder(root *defines.TreeNode) [][]int {
	return dfsLevel(root, -1, [][]int{})
}

func dfsLevel(node *defines.TreeNode, level int, res [][]int) [][]int {
	if node == nil {
		return res
	}
	currentLevel := level + 1
	if len(res) <= currentLevel {
		res = append(res, []int{})
	}
	res[currentLevel] = append(res[currentLevel], node.Val)
	res = dfsLevel(node.Left, currentLevel, res)
	res = dfsLevel(node.Right, currentLevel, res)
	return res
}
