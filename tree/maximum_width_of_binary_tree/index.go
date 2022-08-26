package maximum_width_of_binary_tree

import (
	"GiveMeAnOffer/defines"
)

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

func solution(root *defines.TreeNode) int {
	lvMin := map[int]int{}
	var dfs func(lv, idx int, node *defines.TreeNode) int
	dfs = func(lv, idx int, node *defines.TreeNode) int {
		if node == nil {
			return 0
		}
		if _, ok := lvMin[lv]; !ok {
			lvMin[lv] = idx // 深度优先遍历，先访问到最左边节点，每一层编号的最小值
		}
		// 当前宽度为当前节点编号与最小的编号之差
		width := idx - lvMin[lv] + 1
		// 左子树节点为 2 * idx， 右子树为 2* idx + 1
		leftWidth := dfs(lv+1, idx*2, node.Left)
		rightWidth := dfs(lv+1, idx*2+1, node.Right)
		return max(width, max(leftWidth, rightWidth))
	}

	return dfs(1, 1, root)
}
