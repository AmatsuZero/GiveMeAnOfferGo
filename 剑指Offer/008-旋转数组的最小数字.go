package 剑指Offer

/*
 题目：把一个数组最开始的若干个元素搬到数组的末尾，我们称之为数组的旋转。
	输入一个递增排序的数组的一个旋转，输出旋转数组的最小元素。
	例如数组{3,4,5,1,2}为{1,2,3,4,5}的一个旋转，该数组的最小值为1。
*/

func MinNumberInRotatedArray(numbers []int) int {
	if len(numbers) == 0 {
		panic("Invalid input")
	}
	index1, index2 := 0, len(numbers)-1
	indexMid := index1
	for numbers[index1] >= numbers[index2] { // 由于是旋转数组，左边值必然是大于等于右边的
		if index2-index1 == 1 { // 增序列和降序列相差1，表明找到了（最小数组的下一个就是递增数字）
			indexMid = index2
			break
		}
		indexMid = (index1 + index2) / 2
		// 如果下标为index1、index2 和 indexMid 指向的三个数字相等
		// 则只能顺序查找
		if numbers[index1] == numbers[index2] && numbers[indexMid] == numbers[index1] {
			return minInOrder(numbers, index1, index2)
		}
		// 二分查找
		if numbers[indexMid] >= numbers[index1] {
			index1 = indexMid
		} else {
			index2 = indexMid
		}
	}
	return numbers[indexMid]
}

func minInOrder(numbers []int, index1, index2 int) int {
	result := numbers[index1]
	for i := index1 + 1; i <= index2; i++ {
		if result > numbers[i] {
			result = numbers[i]
		}
	}
	return result
}
