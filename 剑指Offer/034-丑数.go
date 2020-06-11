package 剑指Offer

/*
题目：我们把只包含因子2、3和5的数称作丑数（Ugly Number）。求按从小到大的顺序的第1500个丑数。例如6、8都是丑数，但14不是，因为它包含因子7。习惯上我们把1当做第一个丑数。
*/
func GetUglyNumber(index int) int {
	if index <= 0 {
		return 0
	}
	uglyNumbers := make([]int, index, index)
	uglyNumbers[0] = 1
	nextUglyIndex := 0

	multiply2, multiply3, multiply5 := 0, 0, 0
	for nextUglyIndex < index {
		minNum := min(uglyNumbers[multiply2]*2, uglyNumbers[multiply3]*3, uglyNumbers[multiply5]*5)
		uglyNumbers[nextUglyIndex] = minNum
		for uglyNumbers[multiply2]*2 <= uglyNumbers[nextUglyIndex] {
			multiply2++
		}
		for uglyNumbers[multiply3]*3 <= uglyNumbers[nextUglyIndex] {
			multiply3++
		}
		for uglyNumbers[multiply5]*5 <= uglyNumbers[nextUglyIndex] {
			multiply5++
		}
		nextUglyIndex++
	}
	return uglyNumbers[nextUglyIndex-1]
}

func min(num ...int) (result int) {
	if len(num) == 0 {
		return
	}
	result = num[0]
	if len(num) == 1 {
		return
	}
	for i := 1; i < len(num); i++ {
		if num[i] < result {
			result = num[i]
		}
	}
	return
}
