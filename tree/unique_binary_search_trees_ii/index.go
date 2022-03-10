package unique_binary_search_trees_ii

import (
	"GiveMeAnOffer/defines"
)

// GenerateTrees https://leetcode-cn.com/problems/unique-binary-search-trees-ii/
func GenerateTrees(n int) []*defines.TreeNode {
	if n == 0 {
		return []*defines.TreeNode{}
	}
	return generateBSTree(1, n)
}

func generateBSTree(start, end int) []*defines.TreeNode {
	var tree []*defines.TreeNode
	if start > end {
		tree = append(tree, nil)
		return tree
	}

	for i := start; i <= end; i++ {
		left := generateBSTree(start, i-1)
		right := generateBSTree(i+1, end)
		for _, l := range left {
			for _, r := range right {
				root := &defines.TreeNode{
					Val:   i,
					Left:  l,
					Right: r,
				}
				tree = append(tree, root)
			}
		}
	}
	return tree
}
