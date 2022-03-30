package ktsmallestelementinabst

import "GiveMeAnOffer/defines"

// https://leetcode-cn.com/problems/kth-smallest-element-in-a-bst/
func kthSmallest(root *defines.TreeNode, k int) int {
	res, cnt := 0, 0
	inorder(root, k, &cnt, &res)
	return res
}

func inorder(node *defines.TreeNode, k int, cnt, ans *int) {
	if node == nil {
		return
	}
	*cnt++
	inorder(node.Left, k, cnt, ans)
	if *cnt == k {
		*ans = node.Val
	}
	inorder(node.Right, k, cnt, ans)
}
