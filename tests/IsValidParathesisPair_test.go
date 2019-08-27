package tests

import (
	"GiveMeAnOfferGo/goString"
	"testing"
)

func TestIsValidParenthesisPair(t *testing.T) {
	if goString.IsValidParenthesisPair("()") {
		t.Log("通过")
	} else {
		t.Error("失败")
	}

	if goString.IsValidParenthesisPair("()[]{}") {
		t.Log("通过")
	} else {
		t.Error("失败")
	}

	if !goString.IsValidParenthesisPair("(]") {
		t.Log("通过")
	} else {
		t.Error("失败")
	}

	if !goString.IsValidParenthesisPair("([)]") {
		t.Log("通过")
	} else {
		t.Error("失败")
	}

	if goString.IsValidParenthesisPair("{[]}") {
		t.Log("通过")
	} else {
		t.Error("失败")
	}
}
