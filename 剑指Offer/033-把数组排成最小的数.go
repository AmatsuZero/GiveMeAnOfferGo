package 剑指Offer

import (
	"sort"
	"strconv"
	"strings"
)

/*
题目：输入一个正整数数组，把数组里所有数字拼接起来排成一个数，打印能拼接出的所有数字中最小的一个。
例如输入数组{3,32,321}，则打印出这3个数字能排成的最小数字321323。
*/

func GetMinNumber(numbers []int) string {
	if len(numbers) == 0 {
		return ""
	}
	str := make([]string, len(numbers), len(numbers))
	for i := 0; i < len(numbers); i++ {
		str[i] = strconv.Itoa(numbers[i])
	}
	sort.Slice(str, func(i, j int) bool {
		lhs, rhs := str[i], str[j]
		return strings.Compare(lhs+rhs, rhs+lhs) == -1
	})
	return strings.Join(str, "")
}
