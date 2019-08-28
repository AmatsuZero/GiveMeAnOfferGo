package math

import (
	"testing"
)

func TestReversePositiveInteger(t *testing.T) {
	if Reverse(123) == 321 {
		t.Log("相等")
	} else {
		t.Error("不相等")
	}
}

func TestReverseNegativeInteger(t *testing.T) {
	if Reverse(-123) == -321 {
		t.Log("相等")
	} else {
		t.Error("不相等")
	}
}

func TestReverseIntegerWithZeroTail(t *testing.T) {
	if Reverse(120) == 21 {
		t.Log("相等")
	} else {
		t.Error("不相等")
	}
}

func TestReverseTooLargeInteger(t *testing.T) {
	if Reverse(1534236469) == 0 {
		t.Log("相等")
	} else {
		t.Error("不相等")
	}
}
