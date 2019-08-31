package string

import (
	string2 "GiveMeAnOfferGo/string"
	"testing"
)

func TestLongestCommonPrefix(t *testing.T) {
	input := []string{"flower", "flow", "flight"}
	if string2.LongestCommonPrefix(input) == "fl" {
		t.Log("通过")
	} else {
		t.Error("不通过")
	}

	input = []string{"dog", "racecar", "car"}
	if string2.LongestCommonPrefix(input) == "" {
		t.Log("通过")
	} else {
		t.Error("不通过")
	}
}
