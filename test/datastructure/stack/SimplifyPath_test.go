package stack

import (
	"GiveMeAnOfferGo/datastructure/stack"
	"testing"
)

func TestSimplifyPath(t *testing.T) {
	input := "/home/"
	if stack.SimplifyPath(input) != "/home" {
		t.Fail()
	}

	input = "/../"
	if stack.SimplifyPath(input) != "/" {
		t.Fail()
	}

	input = "/a/./b/../../c/"
	if stack.SimplifyPath(input) != "/c" {
		t.Fail()
	}

	input = "/a/../../b/../c//.//"
	if stack.SimplifyPath(input) != "/c" {
		t.Fail()
	}

	input = "/a//b////c/d//././/.."
	if stack.SimplifyPath(input) != "/a/b/c" {
		t.Fail()
	}
}
