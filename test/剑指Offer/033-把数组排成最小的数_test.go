package 剑指Offer

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/剑指Offer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetMinNumber(t *testing.T) {
	assert.Equal(t, "321323", 剑指Offer.GetMinNumber([]int{3, 32, 321}))
}
