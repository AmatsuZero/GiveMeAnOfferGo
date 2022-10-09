package score_of_parentheses

// ScoreOfParentheses https://leetcode.cn/problems/score-of-parentheses/
func ScoreOfParentheses(s string) int {
	st := []int{0}
	for _, c := range s {
		if c == '(' {
			st = append(st, 0)
		} else {
			v := st[len(st)-1]
			st = st[:len(st)-1]
			if v == 0 { // A为空串
				st[len(st)-1] = 1
			} else {
				st[len(st)-1] = 2 * v
			}
		}
	}

	return st[0]
}
