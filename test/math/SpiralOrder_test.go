package math

import (
	"GiveMeAnOfferGo/array"
	"reflect"
	"testing"
)

func TestSpiralTest(t *testing.T) {
	sm := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	if !reflect.DeepEqual(array.SpiralOrder(sm), []int{1, 2, 3, 6, 9, 8, 7, 4, 5}) {
		t.Fail()
	}

	sm = [][]int{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}}
	if !reflect.DeepEqual(array.SpiralOrder(sm), []int{1, 2, 3, 4, 8, 12, 11, 10, 9, 5, 6, 7}) {
		t.Fail()
	}
}
