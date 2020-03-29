package tree

import (
	"fmt"
	"github.com/AmatsuZero/GiveMeAnOfferGo/Collections/Tree"
	"github.com/AmatsuZero/GiveMeAnOfferGo/test/Utils"
	"testing"
)

func TestInsertionsInSequence(t *testing.T) {
	getInt := Utils.GetInt
	tree := new(Tree.AVLTree)
	for i := 0; i < 15; i++ {
		tree.Insert(getInt(i))
	}
	fmt.Println(tree)
}

func TestRemoveAValue(t *testing.T) {
	getInt := Utils.GetInt
	tree := new(Tree.AVLTree)
	tree.Insert(getInt(15))
	tree.Insert(getInt(10))
	tree.Insert(getInt(16))
	tree.Insert(getInt(18))
	fmt.Println(tree)
	tree.Remove(getInt(10))
	fmt.Println(tree)
}
