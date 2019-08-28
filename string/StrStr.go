package string

import "strings"

func kmpPreprocess(s []string) []int {
	n := len(s)
	match := make([]int, n)
	for i := 0; i < n; i++ {
		match[i] = -1
	}
	j := -1
	for i := 1; i < n; i++ {
		for j >= 0 && s[i] != s[j+1] {
			j = match[j]
		}
		if s[i] == s[j+1] {
			j++
		}
		match[i] = j
	}
	return match
}

func StrStr(haystack string, needle string) int {
	if len(needle) == 0 {
		return 0
	}
	if len(haystack) == 0 {
		return -1
	}
	hStack := strings.Split(haystack, "")
	nStack := strings.Split(needle, "")
	m := len(hStack)
	n := len(nStack)

	match := kmpPreprocess(nStack)
	j := -1
	for i := 0; i < m; i++ {
		for j >= 0 && hStack[i] != nStack[j+1] {
			j = match[j]
		}
		if hStack[i] == nStack[j+1] {
			j++
		}
		if j == n-1 {
			return i - n + 1
		}
	}
	return -1
}
