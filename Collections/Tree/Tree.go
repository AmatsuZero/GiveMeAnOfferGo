package Tree

import (
	"GiveMeAnOfferGo/Collections"
	"GiveMeAnOfferGo/Objects"
	"fmt"
)

type TreeNode struct {
	Value    Objects.ComparableObject
	Children []*TreeNode
}

func (node *TreeNode) IsNil() bool {
	return node.Value == nil
}

func (node *TreeNode) String() string {
	str := fmt.Sprintf("=== TreeNode :%p\n", node)
	str += fmt.Sprintf("Value: %v\n", node.Value)
	return str + "=== end"
}

func NewTreeNode(val Objects.ComparableObject) *TreeNode {
	return &TreeNode{
		Value:    val,
		Children: make([]*TreeNode, 0),
	}
}

func (node *TreeNode) Add(child *TreeNode) {
	if child == nil {
		return
	}
	node.Children = append(node.Children, child)
}

func (node *TreeNode) Search(val Objects.ComparableObject) (result *TreeNode) {
	node.ForEachLevelOrder(func(n *TreeNode) {
		if n.Value.IsEqualTo(val) {
			result = n
		}
	})
	return
}

func (node *TreeNode) ForEachDFS(visit func(node *TreeNode)) {
	if visit == nil {
		return
	}
	visit(node)
	for _, v := range node.Children {
		v.ForEachDFS(visit)
	}
}

func (node *TreeNode) ForEachLevelOrder(visit func(node *TreeNode)) {
	if visit == nil {
		return
	}
	visit(node)
	queue := Collections.NewQueueLinedList()
	for _, v := range node.Children {
		queue.Enqueue(v)
	}
	n, ok := queue.Dequeue().(*TreeNode)
	for ok && n != nil {
		visit(n)
		for _, cn := range n.Children {
			queue.Enqueue(cn)
		}
		n, ok = queue.Dequeue().(*TreeNode)
	}
}
