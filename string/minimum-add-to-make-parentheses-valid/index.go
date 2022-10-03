package minimum_add_to_make_parentheses_valid

// MinAddToMakeValid https://leetcode.cn/problems/minimum-add-to-make-parentheses-valid/
func MinAddToMakeValid(s string) (ans int) {
	cnt := 0
	for _, c := range s {
		if c == '(' {
			cnt++
		} else if cnt > 0 {
			cnt--
		} else {
			ans++
		}
	}

	return ans + cnt
}
