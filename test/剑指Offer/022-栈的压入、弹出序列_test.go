package 剑指Offer

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/剑指Offer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsPopOrder(t *testing.T) {
	assert.True(t, 剑指Offer.IsPopOrder([]int{1, 2, 3, 4, 5}, []int{4, 5, 3, 2, 1}))
	assert.False(t, 剑指Offer.IsPopOrder([]int{1, 2, 3, 4, 5}, []int{4, 3, 5, 1, 2}))
	assert.True(t, 剑指Offer.IsPopOrder([]int{1}, []int{1}))
}
