package implementstackusingqueues

// https://leetcode-cn.com/problems/implement-stack-using-queues/
type MyStack struct {
	enque, deque []int
}

func Constructor() MyStack {
	return MyStack{}
}

func (m *MyStack) Push(x int) {
	m.enque = append(m.enque, x)
}

func (m *MyStack) Pop() int {
	length := len(m.enque)
	for i := 0; i < length-1; i++ {
		m.deque = append(m.deque, m.enque[0])
		m.enque = m.enque[1:]
	}
	topElm := m.enque[0]
	m.enque = m.deque
	m.deque = nil
	return topElm
}

func (m *MyStack) Top() int {
	top := m.Pop()
	m.enque = append(m.enque, top)
	return top
}

func (m *MyStack) Empty() bool {
	return len(m.enque) == 0
}
