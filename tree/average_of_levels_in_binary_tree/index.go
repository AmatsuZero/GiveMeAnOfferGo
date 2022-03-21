package averageoflevelsinbinarytree

import "GiveMeAnOffer/defines"

// AverageOfLevels https://leetcode-cn.com/problems/average-of-levels-in-binary-tree/
func AverageOfLevels(root *defines.TreeNode) (res []float64) {
	if root == nil {
		return []float64{0}
	}
	queue := []*defines.TreeNode{root}
	curNum, nextLevelNum, count, sum := 1, 0, 1, 0
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
			sum += node.Val
			queue = queue[1:]
		}
		if curNum == 0 {
			res = append(res, float64(sum)/float64(count))
			curNum, count, nextLevelNum, sum = nextLevelNum, nextLevelNum, 0, 0
		}
	}
	return
}
