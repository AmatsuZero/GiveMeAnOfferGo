package binary_tree_right_side_view

import "GiveMeAnOffer/defines"

// RightSideView https://leetcode-cn.com/problems/binary-tree-right-side-view/
func RightSideView(root *defines.TreeNode) []int {
	if root == nil {
		return []int{}
	}
	queue := []*defines.TreeNode{root}
	curNum, nextLevelNum := 1, 0
	var res, tmp []int
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
			res = append(res, tmp[len(tmp)-1])
			curNum = nextLevelNum
			nextLevelNum = 0
			tmp = []int{}
		}
	}
	return res
}
