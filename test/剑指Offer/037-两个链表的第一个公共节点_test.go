package 剑指Offer

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/剑指Offer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindFirstCommonNode(t *testing.T) {
	lhs := 剑指Offer.RandomIntList(20)
	node := lhs.NodeAt(0)
	rhs := node
	rhs.AddToTail(6, 7, 8)
	assert.Equal(t, node, 剑指Offer.FindFirstCommonNode(lhs, rhs))

	node = lhs.NodeAt(-1)
	rhs = node
	rhs.AddToTail(10, 11, 12)
	assert.Equal(t, node, 剑指Offer.FindFirstCommonNode(lhs, rhs))

	rhs = 剑指Offer.RandomIntList(10)
	assert.Nil(t, 剑指Offer.FindFirstCommonNode(lhs, rhs))

	assert.Nil(t, 剑指Offer.FindFirstCommonNode(lhs, nil))
}
