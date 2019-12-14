package tree

import (
	"GiveMeAnOfferGo/Collections/Tree"
	"fmt"
	"testing"
)

func TestBST(t *testing.T) {
	bst := new(Tree.BinarySearchTree)
	for i := 0; i < 5; i++ {
		bst.Insert(getInt(i))
	}
	t.Log(bst)
}

func TestBSTContains(t *testing.T) {
	tree := makeBST()
	if tree.Contains(getInt(5)) {
		t.Log("Found 5")
	} else {
		t.Log("Couldn't find 5")
	}
}

func TestRemoveBSTNode(t *testing.T) {
	tree := makeBST()
	fmt.Println("Tree Before removal:")
	fmt.Println(tree)
	tree.Remove(getInt(3))
	fmt.Println("Tree after removing root:")
	fmt.Println(tree)
}

func makeBST() *Tree.BinarySearchTree {
	bst := new(Tree.BinarySearchTree)
	bst.Insert(getInt(3))
	bst.Insert(getInt(1))
	bst.Insert(getInt(4))
	bst.Insert(getInt(0))
	bst.Insert(getInt(2))
	bst.Insert(getInt(5))
	return bst
}

func TestIsBST(t *testing.T) {
	tree := makeBST()
	if !tree.Root.IsBinarySearchTree() {
		t.Fail()
	}
}

func TestIsBSTEqual(t *testing.T) {
	bst := new(Tree.BinarySearchTree)
	bst.Insert(getInt(3))
	bst.Insert(getInt(1))
	bst.Insert(getInt(4))
	bst.Insert(getInt(0))
	bst.Insert(getInt(2))
	bst.Insert(getInt(5))

	bst2 := new(Tree.BinarySearchTree)
	bst2.Insert(getInt(2))
	bst2.Insert(getInt(5))
	bst2.Insert(getInt(3))
	bst2.Insert(getInt(1))
	bst2.Insert(getInt(0))
	bst2.Insert(getInt(4))

	if bst.IsEqualTo(bst2) {
		t.Fail()
	}
}
