package binary_tree_paths

import (
	"GiveMeAnOffer/defines"
	"strconv"
)

// BinaryTreePaths https://leetcode-cn.com/problems/binary-tree-paths/
func BinaryTreePaths(root *defines.TreeNode) []string {
	if root == nil {
		return []string{}
	}
	var res []string
	if root.Left == nil && root.Right == nil {
		return []string{strconv.Itoa(root.Val)}
	}
	tmpLeft := BinaryTreePaths(root.Left)
	for _, node := range tmpLeft {
		res = append(res, strconv.Itoa(root.Val)+"->"+node)
	}
	tmpRight := BinaryTreePaths(root.Right)
	for _, node := range tmpRight {
		res = append(res, strconv.Itoa(root.Val)+"->"+node)
	}
	return res
}
