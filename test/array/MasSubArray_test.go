package array

import (
	"GiveMeAnOfferGo/array"
	"testing"
)

func TestMaxSubArray(t *testing.T) {
	input := []int{-2, 1, -3, 4, -1, 2, 1, -5, 4}
	if array.MaxSubArray(input) == 6 {
		t.Log("Pass")
	} else {
		t.Error("Fail")
	}
}
