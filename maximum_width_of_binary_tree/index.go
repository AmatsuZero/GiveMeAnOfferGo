package maximumwidthofbinarytree

import (
	"GiveMeAnOffer/defines"
	"math"
)

// https://leetcode-cn.com/problems/maximum-width-of-binary-tree/
func WidthOfBinaryTree(root *defines.TreeNode) (res int) {
	if root == nil {
		return 0
	}
	if root.Left == nil && root.Right == nil {
		return 1
	}

	queue := []*defines.TreeNode{&defines.TreeNode{0, root.Left, root.Right}}
	for len(queue) != 0 {
		left, right := math.MinInt32, math.MaxInt32
		// 这里需要注意，先保存 queue 的个数，相当于拿到此层的总个数
		qLen := len(queue)
		for i := 0; i < qLen; i++ {
			node := queue[0]
			queue = queue[1:]
			if node.Left != nil {
				newVal := node.Val * 2
				queue = append(queue, &defines.TreeNode{newVal, node.Left.Left, node.Left.Right})
				if left == math.MaxInt32 || left > newVal {
					left = newVal
				}
				if right == math.MinInt32 || right < newVal {
					right = newVal
				}
			}
			if node.Right != nil {
				newVal := node.Val*2 + 1
				queue = append(queue, &defines.TreeNode{newVal, node.Right.Left, node.Right.Right})
				if left == math.MinInt32 || left > newVal {
					left = newVal
				}
				if right == math.MinInt32 || right < newVal {
					right = newVal
				}
			}
		}
	}

	return
}
