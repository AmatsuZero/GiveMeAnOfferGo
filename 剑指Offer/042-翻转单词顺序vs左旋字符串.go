package 剑指Offer

func reverseString(source string, start, end int) string {
	if len(source) <= 1 || start >= end {
		return source
	}
	chars := []rune(source)
	for start < end {
		chars[start], chars[end] = chars[end], chars[start]
		start++
		end--
	}
	return string(chars)
}

/*
输入一个英文句子，翻转句子中单词的顺序，但单词内字符的顺序不变。为简单起见，标点符号和普通字母一样处理。例如输入字符串"I am a student. "，则输出"student. a am I
*/
func ReverseSentence(source string) string {
	if len(source) == 0 {
		return source
	}
	start, end := 0, len(source)-1
	// 翻转整个句子
	source = reverseString(source, start, end)
	// 翻转句子中的每个单词
	end = 0
	for start < len(source)-1 {
		if string(source[start]) == " " {
			start++
			end++
		} else if string(source[end]) == " " || end == len(source)-1 {
			end--
			source = reverseString(source, start, end)
			end++
			start = end
		} else {
			end++
		}
	}
	return source
}

/*
题目二：字符串的左旋转操作是把字符串前面的若干个字符转移到字符串的尾部。
请定义一个函数实现字符串左旋转操作的功能。比如输入字符串"abcdefg"和数字2，该函数将返回左旋转2位得到的结果"cdefgab”
*/
func LeftRotateString(source string, start int) string {
	if len(source) == 0 || start > len(source) || start == 0 {
		return source
	}

	// 原生切片解法：
	//chars := []rune(source)
	//chars = append(chars[start:], chars[:start]...)
	//return string(chars)

	// 翻转解法
	firstStart, firstEnd := 0, start-1
	secondStart, secondEnd := start, len(source)-1
	// 翻转字符串的前面n个字符
	source = reverseString(source, firstStart, firstEnd)
	// 翻转字符串的后面部分
	source = reverseString(source, secondStart, secondEnd)
	// 翻转整个字符串
	source = reverseString(source, firstStart, secondEnd)
	return source
}
