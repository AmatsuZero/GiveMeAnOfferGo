package 剑指Offer

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/剑指Offer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCustomAtoi(t *testing.T) {
	err, result := 剑指Offer.Atoi("+12")
	assert.Equal(t, 12, result)

	err, result = 剑指Offer.Atoi("+1a")
	assert.Error(t, err)

	err, result = 剑指Offer.Atoi("-12")
	assert.Equal(t, -12, result)

	err, result = 剑指Offer.Atoi("")
	assert.Error(t, err)
}
