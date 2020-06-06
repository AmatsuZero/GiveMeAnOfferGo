package 剑指Offer

/*
 定义栈的数据结构，请在该类型中实现一个能够得到栈的最小元素的min函数。在该栈中，调用min、push及pop的时间复杂度都是O（1)
*/

type StackWithMin struct {
	dataStack []int
	minStack  []int
}

func (stack *StackWithMin) Push(value int) {
	stack.dataStack = append(stack.dataStack, value)
	if len(stack.minStack) == 0 || stack.minStack[len(stack.minStack)-1] > value {
		stack.minStack = append(stack.minStack, value)
	} else {
		stack.minStack = append(stack.minStack, stack.minStack[len(stack.minStack)-1])
	}
}

func (stack *StackWithMin) Pop() (popped int) {
	if len(stack.minStack) == 0 || len(stack.dataStack) == 0 {
		panic("Stack Is Empty")
	}
	popped, stack.dataStack = stack.dataStack[len(stack.dataStack)-1], stack.dataStack[:len(stack.dataStack)-1]
	_, stack.minStack = stack.minStack[len(stack.minStack)-1], stack.minStack[:len(stack.minStack)-1]
	return
}

func (stack *StackWithMin) Min() (min int) {
	if len(stack.minStack) == 0 || len(stack.dataStack) == 0 {
		panic("Stack Is Empty")
	}
	return stack.minStack[len(stack.minStack)-1]
}
