package print_binary_tree

import (
	"GiveMeAnOffer/defines"
	"strconv"
)

// PrintTree
/*
https://leetcode.cn/problems/print-binary-tree/
给你一棵二叉树的根节点 root ，请你构造一个下标从 0 开始、大小为 m x n 的字符串矩阵 res ，用以表示树的 格式化布局 。构造此格式化布局矩阵需要遵循以下规则：

树的 高度 为 height ，矩阵的行数 m 应该等于 height + 1 。
矩阵的列数 n 应该等于 2height+1 - 1 。
根节点 需要放置在 顶行 的 正中间 ，对应位置为 res[0][(n-1)/2] 。
对于放置在矩阵中的每个节点，设对应位置为 res[r][c] ，将其左子节点放置在 res[r+1][c-2height-r-1] ，右子节点放置在 res[r+1][c+2height-r-1] 。
继续这一过程，直到树中的所有节点都妥善放置。
任意空单元格都应该包含空字符串 "" 。
返回构造得到的矩阵 res 。
*/
func PrintTree(root *defines.TreeNode) [][]string {
	if root == nil {
		return nil
	}

	var (
		height = heightOfTree(root) - 1
		m      = height + 1
		n      = pow(2, height+1) - 1
		dfs    func(r, c int, node *defines.TreeNode)
	)

	res := make([][]string, m)
	for i := 0; i < m; i++ {
		res[i] = make([]string, n)
	}

	dfs = func(r, c int, node *defines.TreeNode) {
		if node == nil {
			return
		}
		res[r][c] = strconv.Itoa(node.Val)
		if r == height { // 到达最后一层
			return
		}
		dfs(r+1, c-pow(2, height-r-1), node.Left)
		dfs(r+1, c+pow(2, height-r-1), node.Right)
	}

	dfs(0, (n-1)/2, root)
	return res
}

func heightOfTree(root *defines.TreeNode) int {
	if root == nil {
		return 0
	}
	return max(heightOfTree(root.Left), heightOfTree(root.Right)) + 1
}

func pow(base, rad int) (ans int) {
	ans = 1
	for i := 0; i < rad; i++ {
		ans *= base
	}
	return ans
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
