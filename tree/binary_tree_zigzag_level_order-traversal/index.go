package binary_tree_zigzag_level_order_traversal

import "GiveMeAnOffer/defines"

// ZigzagLevelOrder https://leetcode-cn.com/problems/binary-tree-zigzag-level-order-traversal/
func ZigzagLevelOrder(root *defines.TreeNode) [][]int {
	if root == nil {
		return [][]int{}
	}

	var queue []*defines.TreeNode
	queue = append(queue, root)
	var res [][]int
	var tmp []int
	curNum, nextLevelNum, curDir := 1, 0, 0
	for len(queue) != 0 {
		if curNum > 0 {
			node := queue[0]
			if node.Left != nil {
				queue = append(queue, node.Left)
				nextLevelNum += 1
			}
			if node.Right != nil {
				queue = append(queue, node.Right)
				nextLevelNum += 1
			}
			curNum -= 1
			tmp = append(tmp, node.Val)
			queue = queue[1:]
		}
		if curNum == 0 {
			if curDir == 1 {
				for i, j := 0, len(tmp)-1; i < j; i, j = i+1, j-1 {
					tmp[i], tmp[j] = tmp[j], tmp[i]
				}
			}
			res = append(res, tmp)
			curNum = nextLevelNum
			nextLevelNum = 0
			tmp = []int{}
			if curDir == 0 {
				curDir = 1
			} else {
				curDir = 0
			}
		}
	}
	return res
}
