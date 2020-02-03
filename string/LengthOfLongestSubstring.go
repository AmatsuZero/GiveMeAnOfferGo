package string

func LengthOfLongestSubstring(s string) int {
	n, ans := len(s), 0
	m := make(map[rune]int, 0)
	for j, i := 0, 0; j < n; j++ {
		r := rune(s[j])
		if t, ok := m[r]; ok && t > i {
			i = t
		}
		t := j - i + 1
		if t > ans {
			ans = t
		}
		m[r] = j + 1
	}
	return ans
}
