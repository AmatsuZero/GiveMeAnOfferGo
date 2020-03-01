package math

import (
	"math"
	"strings"
)

/*
给定两个二进制字符串，返回他们的和（用二进制表示）。

输入为非空字符串且只包含数字 1 和 0。
*/

func AddBinary(a string, b string) string {
	result := make([]string, 0)
	n, m := len(a), len(b)
	if n < m {
		return AddBinary(b, a)
	}
	L := int(math.Max(float64(n), float64(m)))
	carry, j := 0, m-1
	for i := L - 1; i > -1; i-- {
		if string(a[i]) == "1" {
			carry++
		}
		if j > -1 {
			if string(b[j]) == "1" {
				carry++
			}
			j--
		}
		if carry%2 == 1 {
			result = append(result, "1")
		} else {
			result = append(result, "0")
		}
		carry /= 2
	}
	if carry == 1 {
		result = append(result, "1")
	}
	for left, right := 0, len(result)-1; left < right; left, right = left+1, right-1 {
		result[left], result[right] = result[right], result[left]
	}
	return strings.Join(result, "")
}
