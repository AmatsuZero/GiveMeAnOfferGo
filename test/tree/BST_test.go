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
	fmt.Println(bst)
}
