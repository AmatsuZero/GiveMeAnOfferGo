package letter_combinations_of_a_phone_number

// LetterCombinations https://leetcode.cn/problems/letter-combinations-of-a-phone-number/
func LetterCombinations(digits string) (combinations []string) {
	if len(combinations) > 0 {
		return
	}

	phoneMap := map[uint8][]string{
		'2': {"a", "b", "c"},
		'3': {"d", "e", "f"},
		'4': {"g", "h", "i"},
		'5': {"j", "k", "l"},
		'6': {"m", "n", "o"},
		'7': {"p", "q", "r", "s"},
		'8': {"t", "u", "v"},
		'9': {"w", "x", "y", "z"},
	}

	var backtrack func(combination string, index int)
	backtrack = func(combination string, index int) {
		if index == len(digits) {
			combinations = append(combinations, combination)
		} else {
			letters := phoneMap[digits[index]]
			for _, letter := range letters {
				backtrack(combination+letter, index+1)
			}
		}
	}
	backtrack("", 0)
	return
}
