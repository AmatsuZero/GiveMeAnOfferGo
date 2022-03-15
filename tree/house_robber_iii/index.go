package house_robber_iii

import (
	"GiveMeAnOffer/defines"
	"math"
)

// Rob https://leetcode-cn.com/problems/house-robber-iii/
func Rob(root *defines.TreeNode) int {
	a, b := dfsTreeRob(root)
	return max(a, b)
}

func dfsTreeRob(root *defines.TreeNode) (a, b int) {
	if root == nil {
		return 0, 0
	}
	l0, l1 := dfsTreeRob(root.Left)
	r0, r1 := dfsTreeRob(root.Right)
	// 当前节点没有被打劫
	tmp0 := max(l0, l1) + max(r0, r1)
	// 当前节点被打劫
	tmp1 := root.Val + l0 + r0
	return tmp0, tmp1
}

func max(nums ...int) (res int) {
	res = math.MinInt32
	for _, n := range nums {
		if res < n {
			res = n
		}
	}
	return
}
