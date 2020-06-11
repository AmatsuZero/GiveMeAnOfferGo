package 剑指Offer

import "fmt"

/*
题目：在字符串中找出第一个只出现一次的字符。如输入"abaccdeff"，则输出'b
*/
func FirstNotRepeatingChar(str string) (result string) {
	if len(str) == 0 {
		return
	}
	table := map[rune]int{}
	for _, char := range str {
		table[char]++
	}
	for _, char := range str {
		if table[char] == 1 {
			result = string(char)
			break
		}
	}
	return
}

/*
定义一个函数，输入两个字符串，从第一个字符串中删除在第二个字符串中出现过的所有字符。例如从第一个字符串"We are students"中删除在第二个字符串"aeiou"中出现过的字符得到的结果是"W r Stdnts".
*/
func RemoveCharacters(source, removeSrc string) string {
	if len(source) == 0 {
		return source
	}
	if len(removeSrc) == 0 {
		return source
	}
	dict := map[rune]bool{}
	for _, char := range removeSrc {
		dict[char] = true
	}
	output := make([]rune, 0)
	for _, char := range source {
		if _, ok := dict[char]; !ok {
			output = append(output, char)
		}
	}
	return string(output)
}

func RemoveRepeatCharacters(str string) string {
	if len(str) <= 1 {
		return str
	}
	dict := map[rune]bool{}
	for _, c := range str {
		dict[c] = true
	}
	source := []rune(str)
	for i := 1; i < len(str); i++ {
		if _, ok := dict[source[i]]; ok {
			source = append(source[:i], source[i+1:]...)
			fmt.Println(source)
		}
	}
	return string(source)
}
