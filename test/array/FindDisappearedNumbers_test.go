package array

import (
	"GiveMeAnOfferGo/array"
	"reflect"
	"testing"
)

func TestFindDisappearedNumbers(t *testing.T) {
	if !reflect.DeepEqual(array.FindDisappearedNumbers([]int{4, 3, 2, 7, 8, 2, 3, 1}), []int{5, 6}) {
		t.Fail()
	}
}
