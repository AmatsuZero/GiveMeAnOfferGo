package 剑指Offer

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/剑指Offer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTreeDepth(t *testing.T) {
	tree := 剑指Offer.Construct([]int{1, 2, 4, 5, 7, 3, 6}, []int{4, 2, 7, 5, 1, 3, 6})
	assert.Equal(t, 4, tree.Depth())

	// 只有左子树
	tree = 剑指Offer.Construct([]int{8, 6, 5}, []int{5, 6, 8})
	assert.Equal(t, 3, tree.Depth())

	// 只有右子树
	tree = 剑指Offer.Construct([]int{8, 10, 11}, []int{8, 10, 11})
	assert.Equal(t, 3, tree.Depth())

	// 只有一个节点
	tree = &剑指Offer.BinaryTreeNode{Value: 1}
	assert.Equal(t, 1, tree.Depth())
}

func TestIsBalancedTree(t *testing.T) {
	// 平衡二叉树
	tree := 剑指Offer.Construct([]int{1, 2, 4, 5, 3, 7, 6}, []int{4, 2, 5, 1, 7, 3, 6})
	assert.True(t, tree.IsBalanced())

	// 只有左子树
	tree = 剑指Offer.Construct([]int{8, 6, 5}, []int{5, 6, 8})
	assert.False(t, tree.IsBalanced())

	// 只有右子树
	tree = 剑指Offer.Construct([]int{8, 10, 11}, []int{8, 10, 11})
	assert.False(t, tree.IsBalanced())

	// 只有一个节点
	tree = &剑指Offer.BinaryTreeNode{Value: 1}
	assert.True(t, tree.IsBalanced())
}
