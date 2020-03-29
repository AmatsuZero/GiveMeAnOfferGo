package array

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/array"
	"testing"
)

func TestRemoveElement(t *testing.T) {
	input := []int{3, 2, 2, 3}
	if array.RemoveElement(input, 3) == 2 {
		t.Log("通过")
	} else {
		t.Error("不通过")
	}

	input = []int{0, 1, 2, 2, 3, 0, 4, 2}
	if array.RemoveElement(input, 2) == 5 {
		t.Log("通过")
	} else {
		t.Error("不通过")
	}
}
