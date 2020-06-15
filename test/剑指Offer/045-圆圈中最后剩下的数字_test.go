package 剑指Offer

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/剑指Offer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLastRemain(t *testing.T) {
	assert.Equal(t, 0, 剑指Offer.LastRemain(5, 2))
	assert.Equal(t, 3, 剑指Offer.LastRemain(5, 3))
	assert.Equal(t, 0, 剑指Offer.LastRemain(6, 6))
	assert.Equal(t, 2, 剑指Offer.LastRemain(6, 7))
}