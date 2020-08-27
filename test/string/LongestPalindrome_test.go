package string

import (
	string2 "github.com/AmatsuZero/GiveMeAnOfferGo/string"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLongestPalindrome(t *testing.T) {
	assert.Equal(t, "bb", string2.LongestPalindrome("cbbd"))
	assert.Equal(t, "bab", string2.LongestPalindrome("babad"))
}
