package string_rotation_lcci

// IsFlippedString IsFlipedString https://leetcode.cn/problems/string-rotation-lcci/
func IsFlippedString(s1 string, s2 string) bool {
	m, n := len(s1), len(s2)
	if m != n {
		return false
	}

	if s1 == s2 {
		return true
	}

	for i := 0; i < n; i++ {
<<<<<<< HEAD
		lhs, rhs := string(s2[:i]), string(s2[i:])
=======
		lhs, rhs := s2[:i], s2[i:]
>>>>>>> 7f45743 (ADD：每日一题)
		if lhs+rhs == s1 || rhs+lhs == s1 {
			return true
		}
	}
	return false
}
