package stack

import (
	"fmt"
	"github.com/AmatsuZero/GiveMeAnOfferGo/Collections"
	"github.com/AmatsuZero/GiveMeAnOfferGo/Objects"
	"testing"
)

func TestPop(t *testing.T) {
	s := new(Collections.Stack)
	for i := 1; i < 5; i++ {
		s.Push(getInt(i))
	}

	fmt.Println(s)
	poppedElement := s.Pop()
	fmt.Printf("Popped: %v\n", poppedElement)
}

func TestIsBalancedParentheses(t *testing.T) {
	str := "h((e))llo(world)()"
	obj := &Objects.StringObject{GoString: &str}
	record := new(Collections.Stack)

	start := "("
	left := getString(&start)
	end := ")"
	right := getString(&end)

	for _, s := range obj.ToObjectSlice() {
		if s.IsEqualTo(left) {
			record.Push(s)
		} else if s.IsEqualTo(right) {
			record.Pop()
		}
	}

	if !record.IsEmpty() {
		t.Fail()
	}

	str2 := "(hello world"
	obj = &Objects.StringObject{GoString: &str2}
	record.RemoveAll()
	for _, s := range obj.ToObjectSlice() {
		if s.IsEqualTo(left) {
			record.Push(s)
		} else if s.IsEqualTo(right) {
			record.Pop()
		}
	}
	if record.IsEmpty() {
		t.Fail()
	}
}

func TestCopy(t *testing.T) {
	s1 := new(Collections.Stack)
	s2 := s1.Copy()
	s2.Push(getInt(1))
	fmt.Printf("s1 : %p\n s2 : %p\n", s1, s2)
}

func getInt(i int) *Objects.NumberObject {
	return Objects.NewNumberWithInt(i)
}

func getString(s *string) *Objects.StringObject {
	return &Objects.StringObject{GoString: s}
}
