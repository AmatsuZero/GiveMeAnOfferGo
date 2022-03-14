package binary_tree_maximum_path_sum

import (
	"GiveMeAnOffer/defines"
	"math"
)

// MaxPathSum https://leetcode-cn.com/problems/binary-tree-maximum-path-sum/submissions/
func MaxPathSum(root *defines.TreeNode) int {
	if root == nil {
		return 0
	}
	max := math.MinInt32
	getPathSum(root, &max)
	return max
}

func getPathSum(root *defines.TreeNode, maxSum *int) int {
	if root == nil {
		return math.MinInt32
	}
	left := getPathSum(root.Left, maxSum)
	right := getPathSum(root.Right, maxSum)
	currMax := max(left+root.Val, right+root.Val, root.Val)
	*maxSum = max(*maxSum, currMax, left+right+root.Val)
	return currMax
}

func max(num ...int) (ret int) {
	ret = math.MinInt32
	for _, n := range num {
		if ret < n {
			ret = n
		}
	}
	return
}
