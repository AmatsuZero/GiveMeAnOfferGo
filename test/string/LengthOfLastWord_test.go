package string

import (
	string2 "github.com/AmatsuZero/GiveMeAnOfferGo/string"
	"testing"
)

func TestLengthOfLastWord(t *testing.T) {

	if string2.LengthOfLastWord("Hello World") == 5 {
		t.Log("Pass")
	} else {
		t.Error("Fail")
	}

	if string2.LengthOfLastWord("a ") == 1 {
		t.Log("Pass")
	} else {
		t.Error("Fail")
	}
}
