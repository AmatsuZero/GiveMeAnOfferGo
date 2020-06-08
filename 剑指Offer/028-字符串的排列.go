package 剑指Offer

import "fmt"

/*
题目：输入一个字符串，打印出该字符串中字符的所有排列。
例如输入字符串abc，则打印出由字符a、b、c所能排列出来的所有字符串abc、acb、bac、bca、cab和cba。
*/

func Permutation(str string) {
	chars := []rune(str)
	permutation(0, &chars)
}

/*
start 指向排列字符串的第一个
*/
func permutation(start int, arranged *[]rune) {
	if len(*arranged) == start {
		fmt.Println(string(*arranged))
		return
	}
	for i := 0; i < len(*arranged); i++ {
		if i != start {
			(*arranged)[start], (*arranged)[i] = (*arranged)[i], (*arranged)[start]
		}
		permutation(start+1, arranged)
		if i != start {
			(*arranged)[start], (*arranged)[i] = (*arranged)[i], (*arranged)[start]
		}
	}
}
