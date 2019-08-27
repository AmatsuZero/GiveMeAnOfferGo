package tests

import (
	"GiveMeAnOfferGo/array"
	"testing"
)

func TestRemoveDuplicates(t *testing.T) {
	input := []int{1, 1, 2}
	if array.RemoveDuplicates(input) == 2 {
		t.Log("通过")
	} else {
		t.Error("不通过")
	}

	input = []int{0, 0, 1, 1, 1, 2, 2, 3, 3, 4}
	if array.RemoveDuplicates(input) == 5 {
		t.Log("通过")
	} else {
		t.Error("不通过")
	}
}
