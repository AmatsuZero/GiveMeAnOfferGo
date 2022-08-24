package is_palindrome

import (
	"strings"
	"unicode"
)

func IsPalindrome(s string) bool {
	s = strings.ToLower(s)
	i, j := 0, len(s)-1
	for i != j && i < len(s) && j >= 0 {
		lhs, rhs := s[i], s[j]
		for !unicode.IsLetter(rune(lhs)) && !unicode.IsDigit(rune(lhs)) {
			i += 1
			if i == len(s) {
				return true
			}
			lhs = s[i]
		}

		for !unicode.IsLetter(rune(rhs)) && !unicode.IsDigit(rune(rhs)) {
			j -= 1
			if j == -1 {
				return true
			}
			rhs = s[j]
		}

		if lhs != rhs {
			return false
		}

		i += 1
		j -= 1
	}

	return true
}
