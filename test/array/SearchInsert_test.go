package array

import (
	"GiveMeAnOfferGo/array"
	"testing"
)

func TestSearchInsert(t *testing.T) {
	input := []int{1, 3, 5, 6}

	if array.SearchInsert(input, 5) == 2 {
		t.Log("通过")
	} else {
		t.Error("失败")
	}

	if array.SearchInsert(input, 2) == 1 {
		t.Log("通过")
	} else {
		t.Error("失败")
	}

	if array.SearchInsert(input, 7) == 4 {
		t.Log("通过")
	} else {
		t.Error("失败")
	}

	if array.SearchInsert(input, 0) == 0 {
		t.Log("通过")
	} else {
		t.Error("失败")
	}
}
