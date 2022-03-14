package path_sum_ii

import "GiveMeAnOffer/defines"

// PathSum https://leetcode-cn.com/problems/path-sum-ii/
func PathSum(root *defines.TreeNode, targetSum int) [][]int {
	return findPath(root, targetSum, [][]int{}, []int{})
}

func findPath(n *defines.TreeNode, sum int, slice [][]int, stack []int) [][]int {
	if n == nil {
		return slice
	}
	sum -= n.Val
	stack = append(stack, n.Val)
	if sum == 0 && n.Left == nil && n.Right == nil {
		slice = append(slice, append([]int{}, stack...))
		stack = stack[:len(stack)-1]
	}
	slice = findPath(n.Left, sum, slice, stack)
	slice = findPath(n.Right, sum, slice, stack)
	return slice
}
