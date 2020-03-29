package math

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/math"
	"testing"
)

func TestMyPow(t *testing.T) {
	if math.MyPow(2.00000, 10) != 1024 {
		t.Fail()
	}

	if math.MyPow(2.10000, 3) != 9.26100 {
		t.Fail()
	}

	if math.MyPow(2.00000, -2) != 0.25 {
		t.Fail()
	}
}
