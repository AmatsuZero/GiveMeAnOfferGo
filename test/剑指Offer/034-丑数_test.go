package 剑指Offer

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/剑指Offer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetUglyNumber(t *testing.T) {
	assert.Equal(t, 2, 剑指Offer.GetUglyNumber(1))
	assert.Equal(t, 0, 剑指Offer.GetUglyNumber(0))
	assert.Equal(t, 1719926784, 剑指Offer.GetUglyNumber(1500))
}
