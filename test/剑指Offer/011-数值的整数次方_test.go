package 剑指Offer

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/剑指Offer"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestPower(t *testing.T) {
	assert.Equal(t, 4.0, 剑指Offer.Power(2, 2))
	assert.Equal(t, 1.0/4.0, 剑指Offer.Power(2, -2))
	assert.Equal(t, -8.0, 剑指Offer.Power(-2, 3))
	assert.Equal(t, -1.0/8.0, 剑指Offer.Power(-2, -3))
	assert.Equal(t, 0.0, 剑指Offer.Power(0, 2))

	math.Pow(12, 12)
}
