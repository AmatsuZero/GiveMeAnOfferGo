package 剑指Offer

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/剑指Offer"
	"testing"
)

func TestConstructBinaryTree(t *testing.T) {
	preOrder := []int{1, 2, 4, 7, 3, 5, 6, 8}
	inOrder := []int{4, 7, 2, 1, 5, 3, 8, 6}

	tree := 剑指Offer.Construct(preOrder, inOrder)
	tree.EnumerateByPreorder(func(value int) {
		t.Log(value)
	})
}
