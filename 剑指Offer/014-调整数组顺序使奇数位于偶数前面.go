package 剑指Offer

/*
输入一个整数数组，实现一个函数来调整该数组中数字的顺序，使得所有奇数位于数组的前半部分，所有偶数位于数组的后半部分。
*/
func ReorderOddEven(nums *[]int) {
	Reorder(nums, func(val int) bool {
		return val&0x1 == 0
	})
}

func Reorder(nums *[]int, block func(val int) bool) {
	if nums == nil || len(*nums) == 0 {
		return
	}
	start, end := 0, len(*nums)-1
	for start < end {
		// 向后移动start，直到它不满足条件
		for start < end && !block((*nums)[start]) {
			start++
		}
		// 向前移动end，直到它满足条件
		for start < end && block((*nums)[end]) {
			end--
		}
		if start < end { // start还在end前面，说明可以交换
			(*nums)[start], (*nums)[end] = (*nums)[end], (*nums)[start]
		}
	}
}
