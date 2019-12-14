package Tree

import (
	"GiveMeAnOfferGo/Collections"
	"GiveMeAnOfferGo/Objects"
	"fmt"
)

type Node struct {
	Value    Objects.ComparableObject
	Children []*Node
}

func (node *Node) IsNil() bool {
	return node.Value == nil
}

func (node *Node) String() string {
	str := fmt.Sprintf("=== TreeNode :%p\n", node)
	str += fmt.Sprintf("Value: %v\n", node.Value)
	return str + "=== end"
}

func NewTreeNode(val Objects.ComparableObject) *Node {
	return &Node{
		Value:    val,
		Children: make([]*Node, 0),
	}
}

func (node *Node) Add(child *Node) {
	if child == nil {
		return
	}
	node.Children = append(node.Children, child)
}

func (node *Node) Search(val Objects.ComparableObject) (result *Node) {
	node.ForEachLevelOrder(func(n *Node) {
		if n.Value.IsEqualTo(val) {
			result = n
		}
	})
	return
}

func (node *Node) ForEachDFS(visit func(node *Node)) {
	if visit == nil {
		return
	}
	visit(node)
	for _, v := range node.Children {
		v.ForEachDFS(visit)
	}
}

func (node *Node) ForEachLevelOrder(visit func(node *Node)) {
	if visit == nil {
		return
	}
	visit(node)
	queue := Collections.NewQueueLinedList()
	for _, v := range node.Children {
		queue.Enqueue(v)
	}
	n, ok := queue.Dequeue().(*Node)
	for ok && n != nil {
		visit(n)
		for _, cn := range n.Children {
			queue.Enqueue(cn)
		}
		n, ok = queue.Dequeue().(*Node)
	}
}
