package 剑指Offer

/*
题目：输入n个整数，找出其中最小的k个数。例如输入4、5、1、6、2、7、3、8这8个数字，则最小的4个数字是1、2、3、4
*/
func GetLeastNumbers(k int, numbers []int) (output []int) {
	if k > len(numbers) || k <= 0 || len(numbers) == 0 {
		return
	}
	start, end := 0, len(numbers)-1
	index := Partition(&numbers, start, end)
	for index != k-1 {
		if index > k-1 {
			end = index - 1
			index = Partition(&numbers, start, end)
		} else {
			start = index + 1
			index = Partition(&numbers, start, end)
		}
	}
	output = make([]int, k, k)
	for i := 0; i < k; i++ {
		output[i] = numbers[i]
	}
	return
}
