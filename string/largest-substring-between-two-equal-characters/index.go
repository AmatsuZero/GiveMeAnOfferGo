package largest_substring_between_two_equal_characters

func MaxLengthBetweenEqualCharacters(s string) int {
	dict := map[rune]int{}
	ans := -1
	for i, c := range s {
		if pre, ok := dict[c]; ok {
			if i-pre > ans {
				ans = i - pre
			}
		} else {
			dict[c] = i
		}
	}
	if ans != -1 {
		ans -= 1
	}
	return ans
}
