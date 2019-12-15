package Tree

import (
	"GiveMeAnOfferGo/Objects"
	"fmt"
)

type AVLNode struct {
	Value      Objects.ComparableObject
	LeftChild  *AVLNode
	RightChild *AVLNode
	Height     int
}

func NewAVLNode(val Objects.ComparableObject) *AVLNode {
	return &AVLNode{Value: val}
}

func (avl *AVLNode) String() string {
	return avl.diagram(avl, "", "", "")
}

func (avl *AVLNode) diagram(node *AVLNode, top string, root string, bottom string) string {
	if node == nil {
		return root + "nil\n"
	}
	if node.LeftChild == nil && node.RightChild == nil {
		return root + fmt.Sprintf("%v\n", node.Value)
	}
	return avl.diagram(node.RightChild, top+" ", top+"┌──", top+"│ ") +
		root +
		fmt.Sprintf("%v\n", node.Value) +
		avl.diagram(node.LeftChild, bottom+"│ ", bottom+"└──", bottom+" ")
}

func (avl *AVLNode) ForEachInOrder(visit func(val Objects.ComparableObject)) {
	if visit == nil {
		return
	}
	if avl.LeftChild != nil {
		avl.LeftChild.ForEachInOrder(visit)
	}
	visit(avl.Value)
	if avl.RightChild != nil {
		avl.RightChild.ForEachInOrder(visit)
	}
}

func (avl *AVLNode) ForEachPreOrder(visit func(val Objects.ComparableObject)) {
	if visit == nil {
		return
	}
	visit(avl.Value)
	if avl.LeftChild != nil {
		avl.LeftChild.ForEachPreOrder(visit)
	}
	if avl.RightChild != nil {
		avl.RightChild.ForEachPreOrder(visit)
	}
}

func (avl *AVLNode) ForEachPostOrder(visit func(val Objects.ComparableObject)) {
	if visit == nil {
		return
	}
	if avl.LeftChild != nil {
		avl.LeftChild.ForEachPostOrder(visit)
	}
	if avl.RightChild != nil {
		avl.RightChild.ForEachPostOrder(visit)
	}
	visit(avl.Value)
}

func (avl *AVLNode) BalanceFactor() int {
	return avl.LeftHeight() - avl.RightHeight()
}

func (avl *AVLNode) LeftHeight() int {
	if avl.LeftChild == nil {
		return -1
	}
	return avl.LeftChild.Height
}

func (avl *AVLNode) RightHeight() int {
	if avl.RightChild == nil {
		return -1
	}
	return avl.RightChild.Height
}
