package string

import (
	string2 "GiveMeAnOfferGo/string"
	"testing"
)

func TestLengthOfLongestSubstring(t *testing.T) {
	if string2.LengthOfLongestSubstring("abcabcbb") != 3 {
		t.Fail()
	}
	if string2.LengthOfLongestSubstring("bbbbb") != 1 {
		t.Fail()
	}
	if string2.LengthOfLongestSubstring("pwwkew") != 3 {
		t.Fail()
	}
}
