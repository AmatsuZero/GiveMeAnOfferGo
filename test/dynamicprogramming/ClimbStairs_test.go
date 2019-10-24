package dynamicprogramming

import (
	"GiveMeAnOfferGo/dynamicprogramming"
	"testing"
)

func TestClimbStairs(t *testing.T) {
	if dynamicprogramming.ClimbStairs(2) != 2 {
		t.Error("Fail")
	}

	if dynamicprogramming.ClimbStairs(3) != 3 {
		t.Error("Fail")
	}
}
