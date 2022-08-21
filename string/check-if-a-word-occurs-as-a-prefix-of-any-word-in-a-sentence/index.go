package check_if_a_word_occurs_as_a_prefix_of_any_word_in_a_sentence

// IsPrefixOfWord https://leetcode.cn/problems/check-if-a-word-occurs-as-a-prefix-of-any-word-in-a-sentence/
func IsPrefixOfWord(sentence string, searchWord string) int {
	wordCnt := 1
	i := 0
	n, m := len(sentence), len(searchWord)
out:
	for ; i < n; i++ {
		if sentence[i] == ' ' {
			wordCnt += 1
			continue
		}

		for j := 0; j < m; j++ {
			if searchWord[j] != sentence[i] {
				for i < n && sentence[i] != ' ' {
					i += 1
				}
				goto out
			} else {
				i += 1
			}
		}
		return wordCnt
	}
	return -1
}
