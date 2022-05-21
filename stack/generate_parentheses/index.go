package generate_parentheses

// GenerateParenthesis https://leetcode.cn/problems/generate-parentheses/
func GenerateParenthesis(n int) (ans []string) {
	str := ""
	backtrack(&ans, &str, 0, 0, n)
	return
}

func backtrack(ans *[]string, cur *string, open, close, max int) {
	if len(*cur) == max*2 {
		*ans = append(*ans, *cur)
		return
	}

	if open < max {
		*cur += "("
		backtrack(ans, cur, open+1, close, max)
		deleteLastChar(cur)
	}

	if close < open {
		*cur += ")"
		backtrack(ans, cur, open, close+1, max)
		deleteLastChar(cur)
	}
}

func deleteLastChar(cur *string) {
	str := *cur
	str = str[:len(str)-1]
	*cur = str
}
