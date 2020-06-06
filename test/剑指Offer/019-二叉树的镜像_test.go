package 剑指Offer

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/剑指Offer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMirrorRecursively(t *testing.T) {
	// 正常翻转
	lhs := 剑指Offer.Construct([]int{8, 6, 10}, []int{6, 8, 10})
	rhs := 剑指Offer.Construct([]int{8, 10, 6}, []int{10, 8, 6})
	lhs.MirrorIteratively()
	assert.Equal(t, lhs, rhs)

	// 只有左子树
	lhs = 剑指Offer.Construct([]int{8, 6, 5}, []int{5, 6, 8})
	rhs = 剑指Offer.Construct([]int{8, 6, 5}, []int{8, 6, 5})
	lhs.MirrorRecursively()
	assert.Equal(t, lhs, rhs)

	// 只有右子树
	lhs = 剑指Offer.Construct([]int{8, 10, 11}, []int{8, 10, 11})
	rhs = 剑指Offer.Construct([]int{8, 10, 11}, []int{11, 10, 8})
	lhs.MirrorIteratively()
	assert.Equal(t, lhs, rhs)

	// 只有一个节点
	lhs = &剑指Offer.BinaryTreeNode{Value: 1}
	rhs = lhs
	lhs.MirrorRecursively()
	assert.Equal(t, lhs, rhs)
}
