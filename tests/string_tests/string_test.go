package string_tests

import (
	"GiveMeAnOffer/leetcode"
	"GiveMeAnOffer/string/solve_the_equation"
	"testing"
)

func TestSolveEquation(t *testing.T) {
	ans := solve_the_equation.SolveEquation("x+5-3+x=6+x-2")
	if ans != "x=2" {
		t.Fail()
	}
	ans = solve_the_equation.SolveEquation("x=x")
	if ans != "Infinite solutions" {
		t.Fail()
	}
	ans = solve_the_equation.SolveEquation("2x=x")
	if ans != "x=0" {
		t.Fail()
	}
}

func TestReplaceWords(t *testing.T) {
	ans := leetcode.ReplaceWords([]string{"cat", "bat", "rat"}, "the cattle was rattled by the battery")
	if ans != "the cattle was rattled by the battery" {
		t.Fail()
	}
}

func TestMinLengthEncoding(t *testing.T) {
	ans := leetcode.MiniMumLengthEncoding([]string{"time", "me", "bell"})
	if ans != 10 {
		t.Fail()
	}
}
