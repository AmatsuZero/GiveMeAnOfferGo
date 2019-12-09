package Collections

import (
	"GiveMeAnOfferGo/Objects"
	"fmt"
	"strings"
)

type Stack struct {
	storage []Objects.Comparable
}

func NewStackWithSlice(elements []Objects.Comparable) *Stack {
	return &Stack{storage: elements}
}

func NewStack(element ...Objects.Comparable) *Stack {
	return &Stack{storage: element}
}

func (stack *Stack) String() string {
	topDivider := "----top----\n"
	bottomDivider := "\n-----------"
	s := stack.storage
	var des []string
	for i := range s {
		des = append(des, fmt.Sprint(s[len(s)-1-i]))
	}
	return topDivider + strings.Join(des, "\n") + bottomDivider
}

func (stack *Stack) Push(val ...Objects.Comparable) {
	stack.storage = append(stack.storage, val...)
}

func (stack *Stack) Pop() (popped Objects.Comparable) {
	if stack.IsEmpty() {
		return nil
	}
	popped, stack.storage = stack.storage[len(stack.storage)-1], stack.storage[:len(stack.storage)-1]
	return
}

func (stack *Stack) Peek() Objects.Comparable {
	if len(stack.storage) == 0 {
		return nil
	}
	return stack.storage[len(stack.storage)-1]
}

func (stack *Stack) IsEmpty() bool {
	return len(stack.storage) == 0
}

func (stack *Stack) AddFromStack(s *Stack) {
	stack.Push(s.storage...)
}

func (stack *Stack) RemoveAll() {
	stack.storage = make([]Objects.Comparable, 0)
}

func (stack *Stack) Copy() *Stack {
	return &Stack{
		storage: append([]Objects.Comparable(nil), stack.storage...),
	}
}

func (stack *Stack) Filter(filter func(index int, element Objects.Comparable) bool) {
	if filter == nil {
		return
	}
	n := 0
	for i, x := range stack.storage {
		if filter(i, x) {
			stack.storage[n] = x
			n++
		}
	}
	stack.storage = stack.storage[:n]
}

func (stack *Stack) Length() int {
	return len(stack.storage)
}

func (stack *Stack) SubRange(from int, to int) *Stack {
	if from > stack.Length() || to > stack.Length() {
		return nil
	}
	return &Stack{
		storage: append(
			stack.storage[:from],
			stack.storage[to:]...),
	}
}

func (stack *Stack) ForEach(traverse func(index int, element Objects.Comparable)) {
	if traverse == nil {
		return
	}
	for i, e := range stack.storage {
		traverse(i, e)
	}
}
