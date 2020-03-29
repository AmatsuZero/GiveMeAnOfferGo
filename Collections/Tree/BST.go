package Tree

import (
	"fmt"
	"github.com/AmatsuZero/GiveMeAnOfferGo/Objects"
)

type BinarySearchTree struct {
	Root *BinaryTreeNode
}

func (bst *BinarySearchTree) String() string {
	if bst.Root == nil {
		return "empty tree"
	}
	return fmt.Sprint(bst.Root)
}

func (bst *BinarySearchTree) Insert(value Objects.ComparableObject) {
	bst.Root = bst.insert(bst.Root, value)
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

func (bst *BinarySearchTree) Contains(value Objects.ComparableObject) (found bool) {
	current := bst.Root
	for current != nil {
		if current.Value.IsEqualTo(value) {
			return true
		}
		if value.Compare(current.Value) == Objects.OrderedAscending {
			current = current.LeftChild
		} else {
			current = current.RightChild
		}
	}
	return
}

func (bst *BinarySearchTree) Remove(value Objects.ComparableObject) {
	bst.Root = bst.remove(bst.Root, value)
}

func (bst *BinarySearchTree) remove(node *BinaryTreeNode, value Objects.ComparableObject) *BinaryTreeNode {
	if node == nil {
		return nil
	}
	if value.IsEqualTo(node.Value) {
		if node.LeftChild == nil && node.RightChild == nil {
			return nil
		}
		if node.LeftChild == nil {
			return node.RightChild
		}
		if node.RightChild == nil {
			return node.LeftChild
		}
		node.Value = node.RightChild.min().Value
		node.RightChild = bst.remove(node.RightChild, node.Value)
	} else if value.Compare(node.Value) == Objects.OrderedAscending {
		node.LeftChild = bst.remove(node.LeftChild, value)
	} else {
		node.RightChild = bst.remove(node.RightChild, value)
	}
	return node
}

func (bst *BinarySearchTree) IsEqualTo(tree interface{}) bool {
	rhs, ok := tree.(*BinarySearchTree)
	if !ok {
		return false
	}
	return bst.isEqualTo(bst.Root, rhs.Root)
}

func (bst *BinarySearchTree) isEqualTo(node1 *BinaryTreeNode, node2 *BinaryTreeNode) bool {
	if node1 == nil || node2 == nil {
		return node1 == nil && node2 == nil
	}
	return node1.Value.IsEqualTo(node2.Value) &&
		bst.isEqualTo(node1.LeftChild, node2.LeftChild) &&
		bst.isEqualTo(node1.RightChild, node2.RightChild)
}
