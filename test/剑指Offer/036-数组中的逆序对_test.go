package 剑指Offer

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/剑指Offer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInReversePairs(t *testing.T) {
	assert.Equal(t, 5, 剑指Offer.InversePairs([]int{7, 5, 6, 4}))
	assert.Equal(t, 0, 剑指Offer.InversePairs([]int{1}))
	assert.Equal(t, 1, 剑指Offer.InversePairs([]int{2, 1}))
	assert.Equal(t, 0, 剑指Offer.InversePairs([]int{}))
	assert.Equal(t, 0, 剑指Offer.InversePairs([]int{1, 2, 3, 4}))
	assert.Equal(t, 6, 剑指Offer.InversePairs([]int{4, 3, 2, 1}))
}
