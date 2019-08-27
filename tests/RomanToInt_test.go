package tests

import (
	"GiveMeAnOfferGo/goString"
	"testing"
)

func TestRomanToInt(t *testing.T) {
	if goString.RomanToInt("III") == 3 {
		t.Log("通过")
	} else {
		t.Error("不通过")
	}

	if goString.RomanToInt("IV") == 4 {
		t.Log("通过")
	} else {
		t.Error("不通过")
	}

	if goString.RomanToInt("LVIII") == 58 {
		t.Log("通过")
	} else {
		t.Error("不通过")
	}

	if goString.RomanToInt("MCMXCIV") == 1994 {
		t.Log("通过")
	} else {
		t.Error("不通过")
	}
}
