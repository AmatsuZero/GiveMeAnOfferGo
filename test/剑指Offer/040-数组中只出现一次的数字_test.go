package 剑指Offer

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/剑指Offer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindAppearOnce(t *testing.T) {
	num1, num2 := 剑指Offer.FindNumsAppearOnce([]int{2, 4, 3, 6, 3, 2, 5, 5})
	assert.Equal(t, 6, num1)
	assert.Equal(t, 4, num2)
}
