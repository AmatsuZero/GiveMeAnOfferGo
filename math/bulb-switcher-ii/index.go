package bulb_switcher_ii

// FlipLights https://leetcode.cn/problems/bulb-switcher-ii/
func FlipLights(n int, presses int) int {
	if presses == 0 {
		return 1
	}
	if n == 1 {
		return 2
	}
	if n+presses == 3 {
		return 3
	}
	if n == 2 || presses == 1 {
		return 4
	}
	if presses == 2 {
		return 7
	}
	return 8
}