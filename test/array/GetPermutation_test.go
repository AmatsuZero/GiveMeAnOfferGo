package array

import (
	"GiveMeAnOfferGo/array"
	"testing"
)

func TestGetPermutation(t *testing.T) {
	if array.GetPermutation(3,3) != "213" {
		t.Fail()
	}
}
