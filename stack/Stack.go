package stack

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

func (stack *Stack) Push(val Objects.Comparable) {
	stack.storage = append(stack.storage, val)
}

func (stack *Stack) Pop() (popped Objects.Comparable) {
	if len(stack.storage) == 0 {
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
