package string

import (
	str "github.com/AmatsuZero/GiveMeAnOfferGo/string"
	"testing"
)

func TestCountAndSay(t *testing.T) {
	if str.CountAndSay(1) == "1" {
		t.Log("通过")
	} else {
		t.Error("失败")
	}

	if str.CountAndSay(4) == "1211" {
		t.Error("通过")
	} else {
		t.Error("失败")
	}
}
