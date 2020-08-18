package Tree

import (
	"fmt"
	"github.com/AmatsuZero/GiveMeAnOfferGo/Objects"
)

type BinaryTreeNode struct {
	Value      Objects.ComparableObject
	LeftChild  *BinaryTreeNode
	RightChild *BinaryTreeNode
}

func NewBinaryNode(val Objects.ComparableObject) *BinaryTreeNode {
	return &BinaryTreeNode{Value: val}
}

func (bn *BinaryTreeNode) String() string {
	return bn.diagram(bn, "", "", "")
}

func (bn *BinaryTreeNode) diagram(node *BinaryTreeNode, top string, root string, bottom string) string {
	if node == nil {
		return root + "nil\n"
	}
	if node.LeftChild == nil && node.RightChild == nil {
		return root + fmt.Sprintf("%v\n", node.Value)
	}
	return bn.diagram(node.RightChild, top+" ", top+"┌──", top+"│ ") +
		root +
		fmt.Sprintf("%v\n", node.Value) +
		bn.diagram(node.LeftChild, bottom+"│ ", bottom+"└──", bottom+" ")
}

func (bn *BinaryTreeNode) ForEachInOrder(visit func(val Objects.ComparableObject)) {
	if visit == nil {
		return
	}
	if bn.LeftChild != nil {
		bn.LeftChild.ForEachInOrder(visit)
	}
	visit(bn.Value)
	if bn.RightChild != nil {
		bn.RightChild.ForEachInOrder(visit)
	}
}

func (bn *BinaryTreeNode) ForEachPreOrder(visit func(val Objects.ComparableObject)) {
	if visit == nil {
		return
	}
	visit(bn.Value)
	if bn.LeftChild != nil {
		bn.LeftChild.ForEachPreOrder(visit)
	}
	if bn.RightChild != nil {
		bn.RightChild.ForEachPreOrder(visit)
	}
}

func (bn *BinaryTreeNode) ForEachPostOrder(visit func(val Objects.ComparableObject)) {
	if visit == nil {
		return
	}
	if bn.LeftChild != nil {
		bn.LeftChild.ForEachPostOrder(visit)
	}
	if bn.RightChild != nil {
		bn.RightChild.ForEachPostOrder(visit)
	}
	visit(bn.Value)
}

func HeightOfTree(node *BinaryTreeNode) int {
	if node == nil {
		return 0
	}
	lhs := HeightOfTree(node.LeftChild)
	rhs := HeightOfTree(node.RightChild)
	if lhs > rhs {
		return 1 + lhs
	} else {
		return 1 + rhs
	}
}

func (bn *BinaryTreeNode) min() *BinaryTreeNode {
	if bn.LeftChild == nil {
		return bn
	}
	return bn.LeftChild.min()
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

func (bn *BinaryTreeNode) IsBinarySearchTree() bool {
	return bn.isBinarySearchTree(bn, nil, nil)
}

func (bn *BinaryTreeNode) isBinarySearchTree(tree *BinaryTreeNode, min Objects.ComparableObject, max Objects.ComparableObject) bool {
	if tree == nil {
		return true
	}
	if min != nil && tree.Value.Compare(min) != Objects.OrderedAscending {
		return false
	} else if max != nil && tree.Value.Compare(max) == Objects.OrderedDescending {
		return false
	}

	return bn.isBinarySearchTree(tree.LeftChild, min, tree.Value) &&
		bn.isBinarySearchTree(tree.RightChild, tree.Value, max)
}
