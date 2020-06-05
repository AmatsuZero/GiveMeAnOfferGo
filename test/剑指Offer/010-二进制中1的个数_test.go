package 剑指Offer

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/剑指Offer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNumberOfOne(t *testing.T) {
	assert.Equal(t, 2, 剑指Offer.NumberOfOne(3))
	assert.Equal(t, 1, 剑指Offer.NumberOfOne(1))
	assert.Equal(t, 0, 剑指Offer.NumberOfOne(0))
	t.Log(剑指Offer.NumberOfOne(0xffffffff))
	t.Log(剑指Offer.NumberOfOne(0x80000000))
}
