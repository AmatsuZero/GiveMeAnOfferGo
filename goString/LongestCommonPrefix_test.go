package goString

import (
	"testing"
)

func TestLongestCommonPrefix(t *testing.T) {
	input := []string{"flower", "flow", "flight"}
	if LongestCommonPrefix(input) == "fl" {
		t.Log("通过")
	} else {
		t.Error("不通过")
	}

	input = []string{"dog", "racecar", "car"}
	if LongestCommonPrefix(input) == "" {
		t.Log("通过")
	} else {
		t.Error("不通过")
	}
}
