package tests

import (
	"GiveMeAnOfferGo/goString"
	"testing"
)

func TestLongestCommonPrefix(t *testing.T) {
	input := []string{"flower", "flow", "flight"}
	if goString.LongestCommonPrefix(input) == "fl" {
		t.Log("通过")
	} else {
		t.Error("不通过")
	}

	input = []string{"dog", "racecar", "car"}
	if goString.LongestCommonPrefix(input) == "" {
		t.Log("通过")
	} else {
		t.Error("不通过")
	}
}
