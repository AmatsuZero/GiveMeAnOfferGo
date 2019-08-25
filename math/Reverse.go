//https://leetcode-cn.com/problems/reverse-integer/description/

package math

import "math"

func Reverse(x int) (result int) {
	if (x >= 0 && x < 10) || (x < 0 && x > -10) {
		return x
	}
	stack := make([]int, 0)
	left := x % 10
	for x != 0 {
		stack = append(stack, left)
		x /= 10
		left = x % 10
	}
	length := len(stack)
	for i, num := range stack {
		pos := 1
		for index := 1; index < length-i; index++ {
			pos *= 10
		}
		result += num * pos
	}
	// 边界检查
	if result < math.MinInt32 || result > math.MaxInt32 {
		result = 0
	}
	return
}
