package 剑指Offer

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/剑指Offer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsStraightCards(t *testing.T) {
	assert.True(t, 剑指Offer.IsStraightCards([]int{5, 4, 3, 2, 1}))  // 顺子
	assert.False(t, 剑指Offer.IsStraightCards([]int{9, 7, 5, 3, 1})) // 不是顺子
    assert.False(t, 剑指Offer.IsStraightCards([]int{5, 4, 3, 3, 2})) // 有对子
	assert.True(t, 剑指Offer.IsStraightCards([]int{5, 4, 3, 2, 0})) // 有 Joker
	assert.True(t, 剑指Offer.IsStraightCards([]int{5, 4, 0, 2, 0})) // 有 Joker
	assert.False(t, 剑指Offer.IsStraightCards(nil))
}
