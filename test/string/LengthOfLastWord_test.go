package string

import (
<<<<<<< Updated upstream
	string2 "GiveMeAnOfferGo/string"
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
=======
	str "GiveMeAnOfferGo/string"
	"testing"
)

func TestLengthOfLastWord(t *testing.T)  {
	if str.LengthOfLastWord("Hello World") == 5 {
		t.Log("Pass")
	} else {
		t.Error("失败")
>>>>>>> Stashed changes
	}
}
