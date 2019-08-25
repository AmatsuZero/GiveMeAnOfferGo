package tests

import (
	"../array"
	"reflect"
	"testing"
)

func TestTwoSum(t *testing.T) {
	input := []int{2, 7, 11, 15}
	if reflect.DeepEqual(array.TwoSum(input, 9), []int{0, 1}) {
		t.Log("相等")
	} else {
		t.Error("不相等")
	}
}
