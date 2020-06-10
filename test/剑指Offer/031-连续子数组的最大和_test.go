package 剑指Offer

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/剑指Offer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindGreatestSumOfSubArray(t *testing.T) {
	input := []int{1, -2, 3, 10, -4, 7, 2, -5}
	assert.Equal(t, 18, 剑指Offer.FindGreatestSumOfSubArray(input))

	input = []int{1, 2, 3, 4, 5}
	assert.Equal(t, 15, 剑指Offer.FindGreatestSumOfSubArray(input))

	input = []int{-1, -2, -3, -4, -5}
	assert.Equal(t, -1, 剑指Offer.FindGreatestSumOfSubArray(input))
}
