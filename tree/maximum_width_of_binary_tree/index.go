package maximum_width_of_binary_tree

import "GiveMeAnOffer/defines"

type annotatedNode struct {
	node       *defines.TreeNode
	depth, pos int
}

// WidthOfBinaryTree https://leetcode-cn.com/problems/maximum-width-of-binary-tree/
func WidthOfBinaryTree(root *defines.TreeNode) int {
	queue := []annotatedNode{{
		root,
		0,
		0,
	}}
	curDepth, left, ans := 0, 0, 0
	for len(queue) != 0 {
		a := queue[0]
		queue = queue[1:]
		if a.node != nil {
			queue = append(queue, annotatedNode{
				node:  a.node.Left,
				depth: a.depth + 1,
				pos:   a.pos * 2,
			})
			queue = append(queue, annotatedNode{
				node:  a.node.Right,
				depth: a.depth + 1,
				pos:   a.pos*2 + 1,
			})
			if curDepth != a.depth {
				curDepth = a.depth
				left = a.pos
			}
			ans = max(ans, a.pos-left+1)
		}
	}
	return ans
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
