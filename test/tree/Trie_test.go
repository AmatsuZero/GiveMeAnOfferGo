package tree

import (
	"GiveMeAnOfferGo/Collections"
	"GiveMeAnOfferGo/Objects"
	"fmt"
	"testing"
)

func TestCreateATree(t *testing.T) {
	beverages := Collections.NewTreeNode(getString("Beverages"))
	hot := Collections.NewTreeNode(getString("hot"))
	cold := Collections.NewTreeNode(getString("cold"))
	beverages.Add(hot)
	beverages.Add(cold)
}

func TestDFSForEach(t *testing.T) {
	tree := makeTree()
	tree.ForEachDFS(func(node *Collections.TreeNode) {
		fmt.Println(node.Value)
	})
}

func makeTree() *Collections.TreeNode {
	beverages := Collections.NewTreeNode(getString("Beverages"))
	hot := Collections.NewTreeNode(getString("hot"))
	cold := Collections.NewTreeNode(getString("cold"))
	beverages.Add(hot)
	beverages.Add(cold)

	tea := Collections.NewTreeNode(getString("tea"))
	coffee := Collections.NewTreeNode(getString("coffee"))
	chocolate := Collections.NewTreeNode(getString("chocolate"))
	hot.Add(tea)
	hot.Add(coffee)
	hot.Add(chocolate)

	blackTea := Collections.NewTreeNode(getString("black"))
	greenTea := Collections.NewTreeNode(getString("green"))
	chaiTea := Collections.NewTreeNode(getString("chai"))
	tea.Add(blackTea)
	tea.Add(greenTea)
	tea.Add(chaiTea)

	soda := Collections.NewTreeNode(getString("soda"))
	milk := Collections.NewTreeNode(getString("milk"))
	cold.Add(soda)
	cold.Add(milk)

	gingerAle := Collections.NewTreeNode(getString("ginger ale"))
	bitterLemon := Collections.NewTreeNode(getString("bitter lemon"))
	soda.Add(gingerAle)
	soda.Add(bitterLemon)

	return beverages
}

func TestLevelOrderTraverse(t *testing.T) {
	tree := makeTree()
	tree.ForEachLevelOrder(func(node *Collections.TreeNode) {
		fmt.Println(node)
	})
}

func TestSearchNode(t *testing.T) {
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

func getString(str string) *Objects.StringObject {
	return &Objects.StringObject{GoString: &str}
}
