package maximum_binary_tree

import (
	"GiveMeAnOffer/defines"
	"math"
)

// ConstructMaximumBinaryTree https://leetcode.cn/problems/maximum-binary-tree/
func ConstructMaximumBinaryTree(nums []int) *defines.TreeNode {
	if len(nums) == 0 {
		return nil
	}
	// 找到最大值，并分成两个数组
	maxNum := math.MinInt64
	pivot := 0
	for i, num := range nums {
		if num > maxNum {
			maxNum = num
			pivot = i
		}
	}
	root := &defines.TreeNode{Val: maxNum}
	lhs, rhs := nums[:pivot], nums[pivot+1:]
	root.Left = ConstructMaximumBinaryTree(lhs)
	root.Right = ConstructMaximumBinaryTree(rhs)
	return root
}
