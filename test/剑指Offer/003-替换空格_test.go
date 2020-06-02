package 剑指Offer

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/剑指Offer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReplaceBlank(t *testing.T) {
	input := []rune("hello world")
	剑指Offer.ReplaceBlank(&input)
	assert.Equal(t, input, []rune("hello%20world"))

	input = []rune("hello")
	剑指Offer.ReplaceBlank(&input)
	assert.Equal(t, input, []rune("hello"))

	input = []rune(" ")
	剑指Offer.ReplaceBlank(&input)
	assert.Equal(t, input, []rune("%20"))

	input = []rune("")
	剑指Offer.ReplaceBlank(&input)
	assert.Equal(t, input, []rune(""))

	input = []rune("   ")
	剑指Offer.ReplaceBlank(&input)
	assert.Equal(t, input, []rune("%20%20%20"))

	input = nil
	剑指Offer.ReplaceBlank(&input)
	assert.Nil(t, input)
}
