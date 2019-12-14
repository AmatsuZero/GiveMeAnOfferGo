package Tree

import (
	"GiveMeAnOfferGo/Objects"
	"fmt"
)

type BinaryTreeNode struct {
	Value      Objects.ComparableObject
	LeftChild  *BinaryTreeNode
	RightChild *BinaryTreeNode
}

func NewBinaryNode(val Objects.ComparableObject) *BinaryTreeNode {
	return &BinaryTreeNode{Value: val}
}

func (node *BinaryTreeNode) String() string {
	return diagram(node, "", "", "")
}

func diagram(node *BinaryTreeNode, top string, root string, bottom string) string {
	if node == nil {
		return root + "nil\n"
	}
	if node.LeftChild == nil && node.RightChild == nil {
		return root + fmt.Sprintf("%v\n", node.Value)
	}
	return diagram(node.RightChild, top+" ", top+"┌──", top+"│ ") +
		root +
		fmt.Sprintf("%v\n", node.Value) +
		diagram(node.LeftChild, bottom+"│ ", bottom+"└──", bottom+" ")
}

func (node *BinaryTreeNode) ForEachInOrder(visit func(val Objects.ComparableObject)) {
	if visit == nil {
		return
	}
	if node.LeftChild != nil {
		node.LeftChild.ForEachInOrder(visit)
	}
	visit(node.Value)
	if node.RightChild != nil {
		node.RightChild.ForEachInOrder(visit)
	}
}

func (node *BinaryTreeNode) ForEachPreOrder(visit func(val Objects.ComparableObject)) {
	if visit == nil {
		return
	}
	visit(node.Value)
	if node.LeftChild != nil {
		node.LeftChild.ForEachPreOrder(visit)
	}
	if node.RightChild != nil {
		node.RightChild.ForEachPreOrder(visit)
	}
}

func (node *BinaryTreeNode) ForEachPostOrder(visit func(val Objects.ComparableObject)) {
	if visit == nil {
		return
	}
	if node.LeftChild != nil {
		node.LeftChild.ForEachPostOrder(visit)
	}
	if node.RightChild != nil {
		node.RightChild.ForEachPostOrder(visit)
	}
	visit(node.Value)
}

func HeightOfTree(node *BinaryTreeNode) int {
	if node == nil {
		return -1
	}
	lhs := HeightOfTree(node.LeftChild)
	rhs := HeightOfTree(node.RightChild)
	if lhs > rhs {
		return 1 + lhs
	} else {
		return 1 + lhs
	}
}

func (node *BinaryTreeNode) min() *BinaryTreeNode {
	if node.LeftChild == nil {
		return node
	}
	return node.LeftChild.min()
}

//func Serialize(node *BinaryTreeNode) (array []*Objects.ComparableObject) {
//	array = make([]*Objects.ComparableObject, 0)
//	if node == nil {
//		return
//	}
//	node.ForEachPreOrder(func(val Objects.ComparableObject) {
//		array = append(array, &val)
//	})
//	return
//}

func (node *BinaryTreeNode) IsBinarySearchTree() bool {
	return isBinarySearchTree(node, nil, nil)
}

func isBinarySearchTree(tree *BinaryTreeNode, min Objects.ComparableObject, max Objects.ComparableObject) bool {
	if tree == nil {
		return true
	}
	if min != nil && tree.Value.Compare(min) != Objects.OrderedAscending {
		return false
	} else if max != nil && tree.Value.Compare(max) == Objects.OrderedDescending {
		return false
	}

	return isBinarySearchTree(tree.LeftChild, min, tree.Value) &&
		isBinarySearchTree(tree.RightChild, tree.Value, max)
}
