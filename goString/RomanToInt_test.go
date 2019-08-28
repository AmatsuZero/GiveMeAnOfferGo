package goString

import (
	"testing"
)

func TestRomanToInt(t *testing.T) {
	if RomanToInt("III") == 3 {
		t.Log("通过")
	} else {
		t.Error("不通过")
	}

	if RomanToInt("IV") == 4 {
		t.Log("通过")
	} else {
		t.Error("不通过")
	}

	if RomanToInt("LVIII") == 58 {
		t.Log("通过")
	} else {
		t.Error("不通过")
	}

	if RomanToInt("MCMXCIV") == 1994 {
		t.Log("通过")
	} else {
		t.Error("不通过")
	}
}
