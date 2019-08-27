package goString

import "strings"

func LongestCommonPrefix(strs []string) (prefix string) {
	// 待检查的数组为空，不检查
	if len(strs) == 0 {
		return
	}
	// 找到最短的一个字符串, 作为最长前缀
	prefix = strs[0]
	for _, str := range strs {
		if len(str) < len(prefix) {
			prefix = str
		}
	}
	// 不断从后向前缩短前缀，直至所有字符串都有此前缀
	for len(prefix) > 0 {
		checkCount := len(strs) // 待检查的数量
		for _, str := range strs {
			if !strings.HasPrefix(str, prefix) {
				prefix = prefix[:len(prefix)-1]
				break
			}
			checkCount -= 1
		}
		if checkCount == 0 { // 等于0，说明此前缀经过了全部检查
			return
		}
	}
	return
}
