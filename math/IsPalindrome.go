// https://leetcode-cn.com/problems/palindrome-number/description/
package math

func IsPalindrome(x int) bool {
	if x < 0 {
		return false
	} else if x == 0 {
		return true
	}
	last := x % 10
	if last == 0 {
		return false
	}
	original := x
	stack := make([]int, 0)
	for x != 0 {
		stack = append(stack, last)
		x /= 10
		last = x % 10
	}

	length := len(stack)
	reversed := 0
	for i, num := range stack {
		pos := 1
		for index := 1; index < length-i; index++ {
			pos *= 10
		}
		reversed += num * pos
	}
	return reversed == original
}
