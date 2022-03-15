package binary_search_tree_iterator

import (
	"GiveMeAnOffer/defines"
	"GiveMeAnOffer/tree/binary_tree_inorder_traversal"
)

// BSTIterator https://leetcode-cn.com/problems/binary-search-tree-iterator/
type BSTIterator struct {
	arr []int
}

func Constructor(root *defines.TreeNode) BSTIterator {
	return BSTIterator{arr: binary_tree_inorder_traversal.InorderTraversal(root)}
}

func (it *BSTIterator) Next() int {
	val := it.arr[0]
	it.arr = it.arr[1:]
	return val
}

func (it *BSTIterator) HasNext() bool {
	return len(it.arr) > 0
}
