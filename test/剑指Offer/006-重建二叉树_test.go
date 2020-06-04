package 剑指Offer

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/剑指Offer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConstructBinaryTree(t *testing.T) {
	preOrder := []int{1, 2, 4, 7, 3, 5, 6, 8}
	inOrder := []int{4, 7, 2, 1, 5, 3, 8, 6}

	tree := 剑指Offer.Construct(preOrder, inOrder)
	lhs := make([]int, 0)
	tree.EnumerateByPreorder(func(value int) {
		lhs = append(lhs, value)
	})
	assert.Equal(t, lhs, preOrder)

	lhs = lhs[:0]
	tree.EnumerateByInorder(func(value int) {
		lhs = append(lhs, value)
	})
	assert.Equal(t, lhs, inOrder)
}
