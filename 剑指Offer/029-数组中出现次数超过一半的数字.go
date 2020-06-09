package 剑指Offer

/*
题目：数组中有一个数字出现的次数超过数组长度的一半，请找出这个数字。
例如输入一个长度为9的数组{1,2,3,2,2,2,5,4,2}。由于数字2在数组中出现了5次，超过数组长度的一半，因此输出2。
*/

func MoreThanHalfNum(numbers []int) (result int) {
	if len(numbers) == 0 {
		panic("invalid input")
	}
	times := 0
	result = numbers[0]
	for i := 1; i < len(numbers); i++ {
		if times == 0 {
			result = numbers[i]
			times++
		} else if numbers[i] == result {
			times++
		}
	}
	return
}
