package math

import (
	"testing"
)

func TestPositiveIntegerIsPalindrome(t *testing.T) {
	if IsPalindrome(121) {
		t.Log("通过")
	} else {
		t.Error("失败")
	}

	if !IsPalindrome(123) {
		t.Log("通过")
	} else {
		t.Error("失败")
	}
}

func TestNegativeIntegerIsPalindrome(t *testing.T) {
	if !IsPalindrome(-121) {
		t.Log("通过")
	} else {
		t.Error("失败")
	}
}

func TestIntegerWithZeroTailIsPalindrome(t *testing.T) {
	if !IsPalindrome(10) {
		t.Log("通过")
	} else {
		t.Error("未通过")
	}
}
