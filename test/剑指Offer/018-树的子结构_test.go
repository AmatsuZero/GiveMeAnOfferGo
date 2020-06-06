package 剑指Offer

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/剑指Offer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHasSubtree(t *testing.T) {
	lhs := 剑指Offer.Construct([]int{1, 8, 9, 2, 4, 7, 6}, []int{9, 8, 4, 2, 7, 1, 6})
	rhs := 剑指Offer.Construct([]int{8, 9, 2}, []int{9, 8, 2})
	assert.True(t, lhs.HasSubtree(rhs))
	rhs = 剑指Offer.Construct([]int{8, 8, 7}, []int{8, 8, 7})
	assert.False(t, lhs.HasSubtree(rhs))
}
