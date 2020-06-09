package 剑指Offer

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/剑指Offer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMoreThanHalfNum(t *testing.T) {
	input := []int{1, 2, 3, 2, 2, 2, 5, 4, 2}
	assert.Equal(t, 2, 剑指Offer.MoreThanHalfNum(input))

	input = []int{1, 2, 3, 4, 5}
	assert.Equal(t, 0, 剑指Offer.MoreThanHalfNum(input))
}
