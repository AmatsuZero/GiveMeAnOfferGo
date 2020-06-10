package 剑指Offer

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/剑指Offer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOnes(t *testing.T) {
	assert.Equal(t, 2, 剑指Offer.Ones(10))
	assert.Equal(t, 5, 剑指Offer.Ones(12))
	assert.Equal(t, 57, 剑指Offer.Ones(132))

	assert.Equal(t, 1, 剑指Offer.Ones(5))
	assert.Equal(t, 0, 剑指Offer.Ones(0))
	assert.Equal(t, 1, 剑指Offer.Ones(1))

	assert.Equal(t, 10000, 10000)
}
