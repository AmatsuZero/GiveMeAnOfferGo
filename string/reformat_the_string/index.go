package reformat_the_string

import (
	"strings"
	"unicode"
)

// Reformat https://leetcode.cn/problems/reformat-the-string/
func Reformat(s string) string {
	var letters []rune
	var digits []rune

	for _, ch := range s {
		if unicode.IsDigit(ch) {
			digits = append(digits, ch)
		} else {
			letters = append(letters, ch)
		}
	}
	m, n := len(letters), len(digits)
	if abs(m-n) > 1 {
		return ""
	}

	ans := strings.Builder{}
	i := 0
	l := min(m, n)
	for ; i < l; i++ {
		if m < n {
			ans.WriteRune(digits[i])
			ans.WriteRune(letters[i])
		} else {
			ans.WriteRune(letters[i])
			ans.WriteRune(digits[i])
		}
	}
	if i < m {
		ans.WriteRune(letters[m-1])
	} else if i < n {
		ans.WriteRune(digits[n-1])
	}

	return ans.String()
}

func abs(num int) int {
	if num > 0 {
		return num
	}
	return -num
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
