package twosumivinputisabst

import "GiveMeAnOffer/defines"

// https://leetcode-cn.com/problems/two-sum-iv-input-is-a-bst/
func FindTarget(root *defines.TreeNode, k int) bool {
	return findTargetDFS(root, k, map[int]int{})
}

func findTargetDFS(root *defines.TreeNode, k int, m map[int]int) bool {
	if root == nil {
		return false
	}
	if _, ok := m[k-root.Val]; ok {
		return ok
	}
	m[root.Val] += 1
	return findTargetDFS(root.Left, k, m) || findTargetDFS(root.Right, k, m)
}
