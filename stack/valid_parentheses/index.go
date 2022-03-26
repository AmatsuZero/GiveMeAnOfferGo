package validparentheses

// https://leetcode-cn.com/problems/valid-parentheses/
func IsValid(s string) bool {
	paris := map[rune]rune{
		')': '(',
		']': '[',
		'}': '{',
	}

	var record []rune
	for _, str := range s {
		switch str {
		case '(', '{', '[':
			record = append(record, str)
		case ')', '}', ']':
			if len(record) == 0 || paris[str] != record[len(record)-1] {
				return false
			}
			record = record[:len(record)-1]
		default:
			return false
		}
	}
	return len(record) == 0
}
