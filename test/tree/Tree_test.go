package tree

import (
	"GiveMeAnOfferGo/Collections/Tree"
	"GiveMeAnOfferGo/Objects"
	"GiveMeAnOfferGo/test/Utils"
	"fmt"
	"testing"
)

func TestCreateATree(t *testing.T) {
	getString := Utils.GetString
	beverages := Tree.NewTreeNode(getString("Beverages"))
	hot := Tree.NewTreeNode(getString("hot"))
	cold := Tree.NewTreeNode(getString("cold"))
	beverages.Add(hot)
	beverages.Add(cold)
}

func TestDFSForEach(t *testing.T) {
	tree := makeTree()
	tree.ForEachDFS(func(node *Tree.Node) {
		fmt.Println(node.Value)
	})
}

func makeTree() *Tree.Node {
	getString := Utils.GetString
	beverages := Tree.NewTreeNode(getString("Beverages"))
	hot := Tree.NewTreeNode(getString("hot"))
	cold := Tree.NewTreeNode(getString("cold"))
	beverages.Add(hot)
	beverages.Add(cold)

	tea := Tree.NewTreeNode(getString("tea"))
	coffee := Tree.NewTreeNode(getString("coffee"))
	chocolate := Tree.NewTreeNode(getString("chocolate"))
	hot.Add(tea)
	hot.Add(coffee)
	hot.Add(chocolate)

	blackTea := Tree.NewTreeNode(getString("black"))
	greenTea := Tree.NewTreeNode(getString("green"))
	chaiTea := Tree.NewTreeNode(getString("chai"))
	tea.Add(blackTea)
	tea.Add(greenTea)
	tea.Add(chaiTea)

	soda := Tree.NewTreeNode(getString("soda"))
	milk := Tree.NewTreeNode(getString("milk"))
	cold.Add(soda)
	cold.Add(milk)

	gingerAle := Tree.NewTreeNode(getString("ginger ale"))
	bitterLemon := Tree.NewTreeNode(getString("bitter lemon"))
	soda.Add(gingerAle)
	soda.Add(bitterLemon)

	return beverages
}

func TestLevelOrderTraverse(t *testing.T) {
	tree := makeTree()
	tree.ForEachLevelOrder(func(node *Tree.Node) {
		fmt.Println(node)
	})
}

func TestSearchNode(t *testing.T) {
	getString := Utils.GetString
	tree := makeTree()
	node := tree.Search(getString("ginger ale"))
	fmt.Printf("Find node: \n%v\n", node)

	node = tree.Search(getString("WKD Blue"))
	if node != nil {
		fmt.Println(node)
	} else {
		fmt.Println("Couldn't find WKD Blue")
	}
}

func TestTreeDiagram(t *testing.T) {
	fmt.Println(makeBinaryTree())
}

func TestTraverseInOrder(t *testing.T) {
	tree := makeBinaryTree()
	tree.ForEachInOrder(func(val Objects.ComparableObject) {
		fmt.Println(val)
	})
}

func TestTraversePreOrder(t *testing.T) {
	tree := makeBinaryTree()
	tree.ForEachPreOrder(func(val Objects.ComparableObject) {
		fmt.Println(val)
	})
}

func TestTraversePostOrder(t *testing.T) {
	tree := makeBinaryTree()
	tree.ForEachPostOrder(func(val Objects.ComparableObject) {
		fmt.Println(val)
	})
}

func TestHeightOfTree(t *testing.T) {
	tree := makeBinaryTree()
	if Tree.HeightOfTree(tree) != 2 {
		t.Fail()
	}
}

func makeBinaryTree() *Tree.BinaryTreeNode {
	getInt := Utils.GetInt
	zero := Tree.NewBinaryNode(getInt(0))
	one := Tree.NewBinaryNode(getInt(1))
	five := Tree.NewBinaryNode(getInt(5))
	seven := Tree.NewBinaryNode(getInt(7))
	eight := Tree.NewBinaryNode(getInt(8))
	nine := Tree.NewBinaryNode(getInt(9))

	seven.LeftChild = one
	one.LeftChild = zero
	one.RightChild = five
	seven.RightChild = nine
	nine.LeftChild = eight

	return seven
}
