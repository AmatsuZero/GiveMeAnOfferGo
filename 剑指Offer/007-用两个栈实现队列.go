package 剑指Offer

/*
题目：用两个栈实现一个队列。

队列的声明如下，请实现它的两个函数appendTail和deleteHead，分别完成在队列尾部插入结点和在队列头部删除结点的功能。
*/
type CQueue struct {
	stackL []int
	stackR []int
}

func (queue *CQueue) AppendTail(node int) {
	queue.stackL = append(queue.stackL, node)
}

func (queue *CQueue) DeleteHead() (poppedValue int) {
	if len(queue.stackR) == 0 {
		for len(queue.stackL) > 0 {
			var top int
			top, queue.stackL = queue.stackL[len(queue.stackL)-1], queue.stackL[:len(queue.stackL)-1]
			queue.stackR = append(queue.stackR, top)
		}
	}
	if len(queue.stackR) == 0 {
		panic("Queue is Empty")
	}
	poppedValue, queue.stackR = queue.stackR[len(queue.stackR)-1], queue.stackR[:len(queue.stackR)-1]
	return
}

/*
 相关问题：

 用两个队列实现一个栈。
*/
type CStack struct {
	queueL []int // 负责入栈
	queueR []int // 负责出栈
}

func (stack *CStack) Push(value int) {
	stack.queueL = append(stack.queueL, value)
}

func (stack *CStack) Pop() (shifted int) {
	for len(stack.queueL) > 0 { // 由于出栈时，后面入队的要先出去，所以要将入栈队列全部清空，元素全部添加到出栈队列
		var front int
		front, stack.queueL = stack.queueL[0], stack.queueL[1:]
		stack.queueR = append([]int{front}, stack.queueR...)
	}
	if len(stack.queueR) == 0 {
		panic("Stack is Empty")
	}
	shifted, stack.queueR = stack.queueR[0], stack.queueR[1:]
	return
}
