package can_measure_water

type state struct {
	remainX, remainY int
}

func (st *state) canFill(targetCapacity int) bool {
	return st.remainX == targetCapacity || st.remainY == targetCapacity || st.remainX+st.remainY == targetCapacity
}

func CanMeasureWater(x, y, targetCapacity int) bool {
	stack := []state{{0, 0}}
	seen := map[state]bool{}

	for len(stack) > 0 {
		var st state
		st, stack = stack[len(stack)-1], stack[:len(stack)-1]
		if seen[st] {
			continue
		}

		seen[st] = true
		if st.canFill(targetCapacity) {
			return true
		}
		// 把壶1灌满
		stack = append(stack, state{x, st.remainY})
		// 把壶2灌满
		stack = append(stack, state{st.remainX, y})
		// 把壶1倒空
		stack = append(stack, state{0, st.remainY})
		// 把壶2倒空
		stack = append(stack, state{st.remainX, 0})
		// 把壶1的水倒入壶2，直至灌满或倒空
		stack = append(stack, state{st.remainX - min(st.remainX, y-st.remainY), st.remainY + min(st.remainX, y-st.remainY)})
		// 把壶2的水倒入壶1，直至灌满或倒空
		stack = append(stack, state{st.remainX + min(st.remainY, x-st.remainX), st.remainY - min(st.remainY, x-st.remainX)})
	}

	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
