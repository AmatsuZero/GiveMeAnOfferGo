package Collections

import (
	"GiveMeAnOfferGo/Objects"
	"fmt"
	"strings"
)

type Stack struct {
	storage []Objects.ObjectProtocol
}

func NewStackWithSlice(elements []Objects.ObjectProtocol) *Stack {
	return &Stack{storage: elements}
}

func NewStack(elements ...Objects.ObjectProtocol) *Stack {
	return NewStackWithSlice(elements)
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

func (stack *Stack) Push(val ...Objects.ObjectProtocol) {
	stack.storage = append(stack.storage, val...)
}

func (stack *Stack) Pop() (popped Objects.ObjectProtocol) {
	if stack.IsEmpty() {
		return nil
	}
	popped, stack.storage = stack.storage[len(stack.storage)-1], stack.storage[:len(stack.storage)-1]
	return
}

func (stack *Stack) Peek() Objects.ObjectProtocol {
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
	stack.storage = make([]Objects.ObjectProtocol, 0)
}

func (stack *Stack) Copy() *Stack {
	return &Stack{
		storage: append([]Objects.ObjectProtocol(nil), stack.storage...),
	}
}

func (stack *Stack) Filter(filter func(index int, element Objects.ObjectProtocol) bool) {
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

func (stack *Stack) ForEach(traverse func(index int, element Objects.ObjectProtocol)) {
	if traverse == nil {
		return
	}
	for i, e := range stack.storage {
		traverse(i, e)
	}
}
