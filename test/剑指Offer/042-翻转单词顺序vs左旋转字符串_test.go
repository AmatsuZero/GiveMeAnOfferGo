package 剑指Offer

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/剑指Offer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReverseSentence(t *testing.T) {
	assert.Equal(t, "student. a am I", 剑指Offer.ReverseSentence("I am a student."))
	assert.Equal(t, "No.", 剑指Offer.ReverseSentence(".No"))
}

func TestLeftRotateString(t *testing.T) {
	assert.Equal(t, "abcdefg", 剑指Offer.LeftRotateString("abcdefg", 0))
	assert.Equal(t, "bcdefga", 剑指Offer.LeftRotateString("abcdefg", 1))
	assert.Equal(t, "cdefgab", 剑指Offer.LeftRotateString("abcdefg", 2))
	assert.Equal(t, "abcdefg", 剑指Offer.LeftRotateString("abcdefg", 7))
}
