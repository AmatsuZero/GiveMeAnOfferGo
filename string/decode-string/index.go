package decode_string

import (
	"strings"
	"unicode"
)

// DecodeString https://leetcode.cn/problems/decode-string/?favorite=2cktkvj
func DecodeString(s string) string {
	ptr := 0

	var getString func() string

	getDigits := func() int {
		ret := 0
		for ; unicode.IsDigit(rune(s[ptr])); ptr++ {
			ret = ret*10 + int(s[ptr]-'0')
		}
		return ret
	}

	getString = func() string {
		if ptr == len(s) || s[ptr] == ']' {
			return ""
		}

		cur := rune(s[ptr])
		repeated := 1
		ret := strings.Builder{}

		if unicode.IsDigit(cur) {
			repeated = getDigits()
			ptr++
			str := getString()
			ptr++
			ret.WriteString(strings.Repeat(str, repeated))
		} else if unicode.IsLetter(cur) {
			ret.WriteRune(cur)
			ptr++
		}
		ret.WriteString(getString())
		return ret.String()
	}

	return getString()
}
