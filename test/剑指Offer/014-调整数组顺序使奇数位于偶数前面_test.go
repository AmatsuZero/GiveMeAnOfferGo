package 剑指Offer

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/剑指Offer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReorderOddOrEven(t *testing.T) {
	source := []int{1, 2, 3, 4, 5}
	剑指Offer.ReorderOddEven(&source)
	assert.Equal(t, []int{1, 5, 3, 4, 2}, source)

	source = []int{2, 4, 1, 3, 5}
	剑指Offer.ReorderOddEven(&source)
	assert.Equal(t, []int{5, 3, 1, 4, 2}, source)
}
