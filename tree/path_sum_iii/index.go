package path_sum_iii

import "GiveMeAnOffer/defines"

// PathSum https://leetcode-cn.com/problems/path-sum-iii/
func PathSum(root *defines.TreeNode, targetSum int) (ans int) {
	if root == nil {
		return 0
	}
	ans = findPath(root, targetSum)
	ans += PathSum(root.Left, targetSum)
	ans += PathSum(root.Right, targetSum)
	return
}

func findPath(root *defines.TreeNode, sum int) (res int) {
	if root == nil {
		return 0
	}
	if root.Val == sum {
		res += 1
	}
	res += findPath(root.Left, sum-root.Val)
	res += findPath(root.Right, sum-root.Val)
	return
}
