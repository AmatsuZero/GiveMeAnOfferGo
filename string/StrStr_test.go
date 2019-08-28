package string

import (
	"testing"
)

func TestStrStr(t *testing.T) {
	if StrStr("hello", "ll") == 2 {
		t.Log("Pass")
	} else {
		t.Error("Fail")
	}

	if StrStr("aaaaa", "bba") == -1 {
		t.Log("Pass")
	} else {
		t.Error("Fail")
	}

	if StrStr("a", "a") == 0 {
		t.Log("Pass")
	} else {
		t.Error("Fail")
	}
}
