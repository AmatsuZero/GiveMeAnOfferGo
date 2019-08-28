package goString

import (
	"testing"
)

func TestIsValidParenthesisPair(t *testing.T) {
	if IsValidParenthesisPair("()") {
		t.Log("通过")
	} else {
		t.Error("失败")
	}

	if IsValidParenthesisPair("()[]{}") {
		t.Log("通过")
	} else {
		t.Error("失败")
	}

	if !IsValidParenthesisPair("(]") {
		t.Log("通过")
	} else {
		t.Error("失败")
	}

	if !IsValidParenthesisPair("([)]") {
		t.Log("通过")
	} else {
		t.Error("失败")
	}

	if IsValidParenthesisPair("{[]}") {
		t.Log("通过")
	} else {
		t.Error("失败")
	}
}
