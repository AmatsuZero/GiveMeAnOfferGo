package trim_a_binary_search_tree

import "GiveMeAnOffer/defines"

func TrimBST(root *defines.TreeNode, low, high int) *defines.TreeNode {
	if root == nil {
		return nil
	}
	if root.Val < low {
		return TrimBST(root.Right, low, high)
	}
	if root.Val > high {
		return TrimBST(root.Left, low, high)
	}
	root.Left = TrimBST(root.Left, low, high)
	root.Right = TrimBST(root.Right, low, high)
	return root
}
