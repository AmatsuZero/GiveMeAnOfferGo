package 剑指Offer

import "strconv"

/*
题目：一个整型数组里除了两个数字之外，其他的数字都出现了两次。请写程序找出这两个只出现一次的数字。要求时间复杂度是O（n），空间复杂度是O（1）
*/
func FindNumsAppearOnce(data []int) (num1, num2 int) {
	if len(data) < 2 {
		return
	}
	resultExclusiveOR := 0
	for _, num := range data {
		resultExclusiveOR ^= num
	}
	indexOf1 := findFirstBitsOne(resultExclusiveOR)
	for j := 0; j < len(data); j++ {
		if isBit(data[j], indexOf1) {
			num1 ^= data[j]
		} else {
			num2 ^= data[j]
		}
	}
	return
}

// 找到整数二进制表示最右边是1的位
func findFirstBitsOne(num int) (indexBit uint) {
	for ((num & 1) == 0) && (indexBit < 8*strconv.IntSize) {
		num = num >> 1
		indexBit++
	}
	return
}

func isBit(num int, indexBit uint) bool {
	num = num >> indexBit
	return (num & 1) != 0
}
