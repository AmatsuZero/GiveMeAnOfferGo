package string

import (
	string2 "GiveMeAnOfferGo/string"
	"testing"
)

func TestIsValidParenthesisPair(t *testing.T) {
	if string2.IsValidParenthesisPair("()") {
		t.Log("通过")
	} else {
		t.Error("失败")
	}

	if string2.IsValidParenthesisPair("()[]{}") {
		t.Log("通过")
	} else {
		t.Error("失败")
	}

	if !string2.IsValidParenthesisPair("(]") {
		t.Log("通过")
	} else {
		t.Error("失败")
	}

	if !string2.IsValidParenthesisPair("([)]") {
		t.Log("通过")
	} else {
		t.Error("失败")
	}

	if string2.IsValidParenthesisPair("{[]}") {
		t.Log("通过")
	} else {
		t.Error("失败")
	}
}
