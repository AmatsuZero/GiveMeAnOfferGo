package count_unique_characters_of_all_substrings_of_a_given_string

// UniqueLetterString https://leetcode.cn/problems/count-unique-characters-of-all-substrings-of-a-given-string/
func UniqueLetterString(s string) (ans int) {
	idx := map[rune][]int{}
	for i, c := range s {
		idx[c] = append(idx[c], i)
	}

	for _, arr := range idx {
		arr = append(append([]int{-1}, arr...), len(s))
		for i := 1; i < len(arr)-1; i++ {
			ans += (arr[i] - arr[i-1]) * (arr[i+1] - arr[i])
		}
	}
	return
}
