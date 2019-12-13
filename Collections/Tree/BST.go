package Tree

import (
	"GiveMeAnOfferGo/Objects"
	"fmt"
)

type BinarySearchTree struct {
	root *BinaryTreeNode
}

func (bst *BinarySearchTree) String() string {
	if bst.root == nil {
		return "empty tree"
	}
	return fmt.Sprint(bst.root)
}

func (bst *BinarySearchTree) Insert(value Objects.ComparableObject) {
	bst.root = bst.insert(bst.root, value)
}

func (bst *BinarySearchTree) insert(node *BinaryTreeNode, value Objects.ComparableObject) *BinaryTreeNode {
	if node == nil {
		return NewBinaryNode(value)
	}
	if value.Compare(node.Value) == Objects.OrderedAscending {
		node.LeftChild = bst.insert(node.LeftChild, value)
	} else {
		node.RightChild = bst.insert(node.RightChild, value)
	}
	return node
}
