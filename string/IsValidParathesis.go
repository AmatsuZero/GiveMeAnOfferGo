package string

import "strings"

func IsValidParenthesisPair(s string) bool {
	length := len(s)
	if length == 0 {
		return true
	} else if length%2 != 0 { // 奇数个，肯定不是成对的
		return false
	}
	stack := make([]string, 0)
	for _, symbol := range strings.Split(s, "") {
		if len(stack) == 0 {
			stack = append(stack, symbol)
			continue
		}
		last := stack[len(stack)-1]
		switch {
		case
			symbol == ")" && last == "(",
			symbol == "}" && last == "{",
			symbol == "]" && last == "[":
			stack = stack[0 : len(stack)-1]
		default:
			stack = append(stack, symbol)
		}
	}
	return len(stack) == 0
}
