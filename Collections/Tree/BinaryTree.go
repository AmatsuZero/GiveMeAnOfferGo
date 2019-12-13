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

func (btn *BinaryTreeNode) String() string {
	return btn.diagram(btn, "", "", "")
}

func (btn *BinaryTreeNode) diagram(node *BinaryTreeNode, top string, root string, bottom string) string {
	if node == nil {
		return root + "nil\n"
	}
	if node.LeftChild == nil && node.RightChild == nil {
		return root + fmt.Sprintf("%v\n", node.Value)
	}
	return btn.diagram(node.RightChild, top+" ", top+"┌──", top+"│ ") +
		root +
		fmt.Sprintf("%v\n", node.Value) +
		btn.diagram(node.LeftChild, bottom+"│ ", bottom+"└──", bottom+" ")
}

func (btn *BinaryTreeNode) ForEachInOrder(visit func(val Objects.ComparableObject)) {
	if visit == nil {
		return
	}
	if btn.LeftChild != nil {
		btn.LeftChild.ForEachInOrder(visit)
	}
	visit(btn.Value)
	if btn.RightChild != nil {
		btn.RightChild.ForEachInOrder(visit)
	}
}

func (btn *BinaryTreeNode) ForEachPreOrder(visit func(val Objects.ComparableObject)) {
	if visit == nil {
		return
	}
	visit(btn.Value)
	if btn.LeftChild != nil {
		btn.LeftChild.ForEachPreOrder(visit)
	}
	if btn.RightChild != nil {
		btn.RightChild.ForEachPreOrder(visit)
	}
}

func (btn *BinaryTreeNode) ForEachPostOrder(visit func(val Objects.ComparableObject)) {
	if visit == nil {
		return
	}
	if btn.LeftChild != nil {
		btn.LeftChild.ForEachPostOrder(visit)
	}
	if btn.RightChild != nil {
		btn.RightChild.ForEachPostOrder(visit)
	}
	visit(btn.Value)
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
