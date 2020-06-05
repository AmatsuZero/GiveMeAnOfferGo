package 剑指Offer

import "fmt"

/*
题目：输入数字n，按顺序打印出从1最大的n位十进制数。
比如输入3，则打印出1、2、3一直到最大的3位数即999。
*/

func PrintOneToMaxDigits(n int) {
	if n <= 0 {
		return
	}
	number := make([]rune, n)
	for i := range number {
		number[i] = '0'
	}
	for i := 0; i < 10; i++ {
		number[0] = rune(i) + '0'
		printOneToMaxOfDigitsRecursively(number, 0)
	}
}

func printOneToMaxOfDigitsRecursively(num []rune, index int) {
	if index == len(num)-1 {
		printNumber(num)
		return
	}
	for i := 0; i < 10; i++ {
		num[index+1] = rune(i) + '0'
		printOneToMaxOfDigitsRecursively(num, index+1)
	}
}

func printNumber(num []rune) {
	isBeginning0 := true
	nLength := len(num)
	for i := 0; i < nLength; i++ {
		if isBeginning0 && num[i] != '0' { // 从不为0的地方开始打印
			isBeginning0 = false
		}
		if !isBeginning0 {
			fmt.Printf("%v", string(num[i]))
		}
	}
	fmt.Print("\t")
}
