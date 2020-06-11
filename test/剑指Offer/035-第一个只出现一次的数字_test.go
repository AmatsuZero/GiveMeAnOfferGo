package 剑指Offer

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/剑指Offer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFirstNotRepeatingChar(t *testing.T) {
	assert.Equal(t, "b", 剑指Offer.FirstNotRepeatingChar("abaccdeff"))
	assert.Equal(t, "a", 剑指Offer.FirstNotRepeatingChar("abcdef"))
	assert.Equal(t, "", 剑指Offer.FirstNotRepeatingChar("aabbccddeeff"))
	assert.Equal(t, "", 剑指Offer.FirstNotRepeatingChar(""))
}

func TestRemoveCharacters(t *testing.T) {
	assert.Equal(t, "W r stdnts", 剑指Offer.RemoveCharacters("We are students", "aeiou"))
	assert.Equal(t, "We are students", 剑指Offer.RemoveCharacters("We are students", ""))
	assert.Equal(t, "", 剑指Offer.RemoveCharacters("", "---"))
}

func TestRemoveRepeatCharacters(t *testing.T) {
	assert.Equal(t, "abc", 剑指Offer.RemoveRepeatCharacters("aabbcc"))
}
