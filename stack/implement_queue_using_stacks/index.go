package implementqueueusingstacks

// https://leetcode-cn.com/problems/implement-queue-using-stacks/
type MyQueue struct {
	in, out []int
}

func Constructor() MyQueue {
	return MyQueue{}
}

func (m *MyQueue) Push(x int) {
	m.in = append(m.in, x)
}

func (m *MyQueue) in2out() {
	for len(m.in) > 0 {
		m.out = append(m.out, m.in[len(m.in)-1])
		m.in = m.in[:len(m.in)-1]
	}
}

func (m *MyQueue) Pop() int {
	if len(m.out) == 0 {
		m.in2out()
	}
	x := m.out[len(m.out)-1]
	m.out = m.out[:len(m.out)-1]
	return x
}

func (m *MyQueue) Peek() int {
	if len(m.out) == 0 {
		m.in2out()
	}
	return m.out[len(m.out)-1]
}

func (m *MyQueue) Empty() bool {
	return len(m.in) == 0 && len(m.out) == 0
}
