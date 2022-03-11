package validate_binary_search_tree

import (
	"GiveMeAnOffer/defines"
	"GiveMeAnOffer/tree/binary_tree_inorder_traversal"
)

// IsValidBST https://leetcode-cn.com/problems/validate-binary-search-tree/
func IsValidBST(root *defines.TreeNode) bool {
	arr := binary_tree_inorder_traversal.InorderTraversal(root)
	for i := 1; i < len(arr); i++ {
		if arr[i-1] >= arr[i] {
			return false
		}
	}
	return true
}
