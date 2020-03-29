package Tree

import (
	"fmt"
	"github.com/AmatsuZero/GiveMeAnOfferGo/Objects"
	"math"
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

func (avl *AVLNode) Min() *AVLNode {
	if avl.LeftChild == nil {
		return avl
	}
	return avl.LeftChild.Min()
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

type AVLTree struct {
	Root *AVLNode
}

func (avl *AVLTree) leftRotate(node *AVLNode) *AVLNode {
	if node == nil {
		return nil
	}
	pivot := node.RightChild
	node.RightChild = pivot.LeftChild
	pivot.LeftChild = node
	node.Height = int(math.Max(float64(node.LeftHeight()), float64(node.RightHeight())) + 1)
	pivot.Height = int(math.Max(float64(pivot.LeftHeight()), float64(pivot.RightHeight())) + 1)
	return pivot
}

func (avl *AVLTree) rightRotate(node *AVLNode) *AVLNode {
	pivot := node.LeftChild
	node.LeftChild = pivot.RightChild
	pivot.RightChild = node
	node.Height = int(math.Max(float64(node.LeftHeight()), float64(node.RightHeight())) + 1)
	pivot.Height = int(math.Max(float64(pivot.LeftHeight()), float64(pivot.RightHeight())) + 1)
	return pivot
}

func (avl *AVLTree) RightLeftRotate(node *AVLNode) *AVLNode {
	if node == nil || node.RightChild == nil {
		return node
	}
	node.RightChild = avl.rightRotate(node.RightChild)
	return avl.leftRotate(node)
}

func (avl *AVLTree) LeftRightRotate(node *AVLNode) *AVLNode {
	if node == nil || node.LeftChild == nil {
		return node
	}
	node.LeftChild = avl.leftRotate(node.LeftChild)
	return avl.rightRotate(node)
}

func (avl *AVLTree) balanced(node *AVLNode) *AVLNode {
	if node == nil {
		return node
	}
	switch node.BalanceFactor() {
	case 2:
		if node.LeftChild != nil && node.LeftChild.BalanceFactor() == -1 {
			return avl.LeftRightRotate(node)
		} else {
			return avl.rightRotate(node)
		}
	case -2:
		if node.RightChild != nil && node.RightChild.BalanceFactor() == 1 {
			return avl.RightLeftRotate(node)
		} else {
			return avl.leftRotate(node)
		}
	default:
		return node
	}
}

func (avl *AVLTree) Insert(value Objects.ComparableObject) {
	avl.Root = avl.insert(avl.Root, value)
}

func (avl *AVLTree) insert(node *AVLNode, value Objects.ComparableObject) *AVLNode {
	if node == nil {
		return NewAVLNode(value)
	}
	if value.Compare(node.Value) == Objects.OrderedAscending {
		node.LeftChild = avl.insert(node.LeftChild, value)
	} else {
		node.RightChild = avl.insert(node.RightChild, value)
	}
	balancedNode := avl.balanced(node)
	balancedNode.Height = int(math.Max(float64(balancedNode.LeftHeight()), float64(balancedNode.RightHeight())) + 1)
	return balancedNode
}

func (avl *AVLTree) String() string {
	if avl.Root == nil {
		return "empty tree"
	}
	return fmt.Sprint(avl.Root)
}

func (avl *AVLTree) Contains(value Objects.ComparableObject) bool {
	current := avl.Root
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
	return false
}

func (avl *AVLTree) Remove(value Objects.ComparableObject) {
	avl.Root = avl.remove(avl.Root, value)
}

func (avl *AVLTree) remove(node *AVLNode, value Objects.ComparableObject) *AVLNode {
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
		node.Value = node.RightChild.Min().Value
		node.RightChild = avl.remove(node.RightChild, node.Value)
	} else if value.Compare(node.Value) == Objects.OrderedAscending {
		node.LeftChild = avl.remove(node.LeftChild, value)
	} else {
		node.RightChild = avl.remove(node.RightChild, value)
	}
	balancedNode := avl.balanced(node)
	balancedNode.Height = int(math.Max(float64(balancedNode.LeftHeight()), float64(balancedNode.RightHeight())) + 1)
	return balancedNode
}
