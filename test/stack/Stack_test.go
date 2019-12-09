package stack

import (
	"GiveMeAnOfferGo/Objects"
	stack2 "GiveMeAnOfferGo/stack"
	"fmt"
	"testing"
)

func TestPop(t *testing.T) {
	stack := new(stack2.Stack)
	for i := 1; i < 5; i++ {
		stack.Push(getInt(i))
	}

	fmt.Println(stack)
	poppedElement := stack.Pop()
	fmt.Printf("Popped: %v\n", poppedElement)
}

func TestIsBalancedParentheses(t *testing.T) {
	str := Objects.StringObject{GoString: "h((e))llo(world)()"}.ToObjectSlice()
	record := new(stack2.Stack)
	for _, s := range str {
		if s.IsEqualTo(getString("(")) {
			record.Push(s)
		} else if s.IsEqualTo(getString(")")) && !record.IsEmpty() {
			record.Pop()
		}
	}
}

func getInt(i int) *Objects.NumberObject {
	return Objects.NewNumberWithInt(i)
}

func getString(s string) *Objects.StringObject {
	return &Objects.StringObject{GoString: s}
}
