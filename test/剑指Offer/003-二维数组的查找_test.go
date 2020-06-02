package 剑指Offer

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/剑指Offer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindInPartiallySortedArray(t *testing.T) {
	testInput := [][]int{
		{1, 2, 8, 9},
		{2, 4, 9, 12},
		{4, 7, 10, 13},
		{6, 8, 11, 15},
	}
	assert.True(t, 剑指Offer.Find(testInput, 10))
	assert.False(t, 剑指Offer.Find(testInput, 3))
	assert.False(t, 剑指Offer.Find([][]int{{}}, 7))
}
