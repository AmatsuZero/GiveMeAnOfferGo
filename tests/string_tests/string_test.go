package string_tests

import (
	"GiveMeAnOffer/string/solve_the_equation"
	"testing"
)

func TestSolveEquation(t *testing.T) {
	ans := solve_the_equation.SolveEquation("x+5-3+x=6+x-2")
	if ans == "x=2" {
		t.Logf("right")
	}
	ans = solve_the_equation.SolveEquation("x=x")
	if ans == "Infinite solutions" {
		t.Logf("right")
	}
	ans = solve_the_equation.SolveEquation("2x=x")
	if ans == "x=0" {
		t.Logf("right")
	}
}
