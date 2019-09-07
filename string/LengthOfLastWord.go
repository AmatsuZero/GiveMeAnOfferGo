package string

import "strings"

func LengthOfLastWord(s string) int {
	if len(s) == 0 {
		return 0
	}
	str := strings.Split(s, "")
	length := len(str)
	res := 0
	idx := length - 1
	for str[idx] == " " && idx != 0 {
		idx--
	}
	for idx != -1 && str[idx] != " " {
		res++
		idx--
	}
	return res
}
