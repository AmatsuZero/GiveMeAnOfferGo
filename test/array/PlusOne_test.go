package array

import (
	"GiveMeAnOfferGo/array"
	"reflect"
	"testing"
)

func TestPlusOne(t *testing.T)  {
	if reflect.DeepEqual([]int{1,2,4}, array.PlusOne([]int{1,2,3}))  {
		t.Log("Pass")
	} else {
		t.Error("Fail")
	}

	if reflect.DeepEqual([]int{4,3,2,1}, array.PlusOne([]int{4,3,2,2}))  {
		t.Log("Pass")
	} else {
		t.Error("Fail")
	}

	if reflect.DeepEqual([]int{9,9}, array.PlusOne([]int{1,0,0}))  {
		t.Log("Pass")
	} else {
		t.Error("Fail")
	}
}
