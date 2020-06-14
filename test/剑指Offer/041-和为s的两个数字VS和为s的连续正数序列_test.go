package 剑指Offer

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/剑指Offer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindNumberWithSum(t *testing.T) {
	found, num1, num2 := 剑指Offer.FindNumberWithSum([]int{1, 2, 4, 7, 11, 15}, 15)
	assert.True(t, found)
	assert.Equal(t, 4, num1)
	assert.Equal(t, 11, num2)

	found, num1, num2 = 剑指Offer.FindNumberWithSum([]int{1, 2, 4, 7, 11}, 17)
	assert.False(t, found)
}

func TestFindContinuousSequence(t *testing.T) {
	assert.Equal(t, [][]int{
		{1, 2, 3, 4, 5},
		{4, 5, 6},
		{7, 8},
	}, 剑指Offer.FindContinuousSequence(15))

	assert.Nil(t, 剑指Offer.FindContinuousSequence(0))
	assert.Nil(t, 剑指Offer.FindContinuousSequence(4))
}
