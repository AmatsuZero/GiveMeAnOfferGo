package flatten_binary_tree_to_linked_list

import "GiveMeAnOffer/defines"

func Flatten(root *defines.TreeNode) {
	if root == nil || (root.Left == nil && root.Right == nil) {
		return
	}
	Flatten(root.Left)
	Flatten(root.Right)
	currentRight := root.Right
	root.Right = root.Left
	root.Left = nil
	for root.Right != nil {
		root = root.Right
	}
	root.Right = currentRight
}
