package array

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/array"
	"testing"
)

func TestGetPermutation(t *testing.T) {
	if array.GetPermutation(3, 3) != "213" {
		t.Fail()
	}

	if array.GetPermutation(4, 9) != "2314" {
		t.Fail()
	}
}
