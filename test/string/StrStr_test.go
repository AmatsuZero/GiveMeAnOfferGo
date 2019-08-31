package string

import (
	string2 "GiveMeAnOfferGo/string"
	"testing"
)

func TestStrStr(t *testing.T) {
	if string2.StrStr("hello", "ll") == 2 {
		t.Log("Pass")
	} else {
		t.Error("Fail")
	}

	if string2.StrStr("aaaaa", "bba") == -1 {
		t.Log("Pass")
	} else {
		t.Error("Fail")
	}

	if string2.StrStr("a", "a") == 0 {
		t.Log("Pass")
	} else {
		t.Error("Fail")
	}
}
