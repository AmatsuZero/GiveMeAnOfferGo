package 剑指Offer

func ReplaceBlank(str *[]rune) {
	if str == nil || len(*str) == 0 {
		return
	}
	target := " "
	replacement := "%20"

	originalLength := len(*str) // 原始长度
	numberOfBlank := 0

	for i := 0; i < originalLength; i++ {
		if string((*str)[i]) == target {
			numberOfBlank++
		}
	}

	newLength := originalLength + (len(replacement)-len(target))*numberOfBlank // 转换后的长度
	for i := originalLength; i < newLength; i++ {                              // 向后追加字符
		*str = append(*str, ' ')
	}
	indexOfOriginal, indexOfNew := originalLength-1, newLength-1

	for indexOfOriginal >= 0 && indexOfNew > indexOfOriginal {
		if string((*str)[indexOfOriginal]) == target {
			for i := range replacement {
				(*str)[indexOfNew] = rune(replacement[len(replacement)-i-1])
				indexOfNew--
			}
		} else {
			(*str)[indexOfNew] = (*str)[indexOfOriginal]
			indexOfNew--
		}
		indexOfOriginal--
	}
	return
}
