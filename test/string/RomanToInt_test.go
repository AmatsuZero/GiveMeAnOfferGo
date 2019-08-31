package string

import (
	string2 "GiveMeAnOfferGo/string"
	"testing"
)

func TestRomanToInt(t *testing.T) {
	if string2.RomanToInt("III") == 3 {
		t.Log("通过")
	} else {
		t.Error("不通过")
	}

	if string2.RomanToInt("IV") == 4 {
		t.Log("通过")
	} else {
		t.Error("不通过")
	}

	if string2.RomanToInt("LVIII") == 58 {
		t.Log("通过")
	} else {
		t.Error("不通过")
	}

	if string2.RomanToInt("MCMXCIV") == 1994 {
		t.Log("通过")
	} else {
		t.Error("不通过")
	}
}
