package minstack

// MinStack https://leetcode-cn.com/problems/min-stack/
type MinStack struct {
	elements, min []int
	l             int
}

func Constructor() MinStack {
	return MinStack{}
}

func (m *MinStack) Push(val int) {
	m.elements = append(m.elements, val)
	if m.l == 0 {
		m.min = append(m.min, val)
	} else {
		min := m.GetMin()
		if val < min {
			m.min = append(m.min, val)
		} else {
			m.min = append(m.min, min)
		}
	}
	m.l += 1
}

func (m *MinStack) Pop() {
	m.l -= 1
	m.min = m.min[:m.l]
	m.elements = m.elements[:m.l]
}

func (m *MinStack) Top() int {
	return m.elements[m.l-1]
}

func (m *MinStack) GetMin() int {
	return m.min[m.l-1]
}
