package 剑指Offer

type CQueue struct {
	stackL []int
	stackR []int
}

func NewCQueue() *CQueue {
	return &CQueue{
		stackL: make([]int, 0),
		stackR: make([]int, 0),
	}
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
		panic("Stack is Empty")
	}
	poppedValue, queue.stackR = queue.stackR[len(queue.stackR)-1], queue.stackR[:len(queue.stackR)-1]
	return
}
