package tests

import (
	"../string"
	"testing"
)

func TestRomanToInt(t *testing.T) {
	if string.RomanToInt("III") == 3 {
		t.Log("通过")
	} else {
		t.Error("不通过")
	}

	if string.RomanToInt("IV") == 4 {
		t.Log("通过")
	} else {
		t.Error("不通过")
	}

	if string.RomanToInt("LVIII") == 58 {
		t.Log("通过")
	} else {
		t.Error("不通过")
	}

	if string.RomanToInt("MCMXCIV") == 1994 {
		t.Log("通过")
	} else {
		t.Error("不通过")
	}
}
