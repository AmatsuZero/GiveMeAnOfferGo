package 剑指Offer

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/剑指Offer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetLeastNumbers(t *testing.T) {
	input := []int{4, 5, 1, 6, 2, 7, 3, 8}
	result := 剑指Offer.GetLeastNumbers(4, input)
	assert.Equal(t, []int{1, 2, 3, 4}, result)
}
