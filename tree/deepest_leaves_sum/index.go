package deepest_leaves_sum

import "GiveMeAnOffer/defines"

// DeepestLeavesSum https://leetcode.cn/problems/deepest-leaves-sum/
func DeepestLeavesSum(root *defines.TreeNode) int {
	sum := make(map[int]int)
	maxLv := 0
	var dfs func(lv int, node *defines.TreeNode)
	dfs = func(lv int, node *defines.TreeNode) {
		if node == nil {
			return
		}
		sum[lv] += node.Val
		if lv > maxLv {
			maxLv = lv
		}
		dfs(lv+1, node.Left)
		dfs(lv+1, node.Right)
	}
	dfs(0, root)
	return sum[maxLv]
}
