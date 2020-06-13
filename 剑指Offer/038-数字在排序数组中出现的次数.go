package 剑指Offer

/*
题目：统计一个数字在排序数组中出现的次数。例如输入排序数组{1,2,3,3,3,3,4,5}和数字3，由于3在这个数组中出现了4次，因此输出4
*/
func GetNumberOfK(data []int, k int) (number int) {
	if len(data) == 0 {
		return
	}
	first := getFirstK(data, k, 0, len(data)-1)
	last := getLastK(data, k, 0, len(data)-1)
	if first > -1 && last > -1 {
		number = last - first + 1
	}
	return number
}

// 找到第一个K
func getFirstK(data []int, k, start, end int) int {
	if start > end {
		return -1
	}
	mid := (end + start) / 2
	middleData := data[mid]
	if middleData == k { // 中位数正好等于 K， 判断是不是第一个K
		if (mid > 0 && data[mid-1] != k) || mid == 0 { // 如果前一个数字不是k，或者正好是第一个元素
			return mid
		} else {
			end = mid - 1
		}
	} else if middleData > k {
		end = mid - 1
	} else {
		start = mid + 1
	}
	return getFirstK(data, k, start, end)
}

// 找到最后一个
func getLastK(data []int, k, start, end int) int {
	if start > end {
		return -1
	}
	midIndex := (start + end) / 2
	middleData := data[midIndex]
	if middleData == k {
		if (midIndex < len(data)-1 && data[midIndex+1] != k) || midIndex == len(data)-1 { // 如果后一个数字不是k，或者正好是最后一个数字
			return midIndex
		} else {
			start = midIndex + 1
		}
	} else if middleData < k {
		start = midIndex + 1
	} else {
		end = midIndex - 1
	}
	return getLastK(data, k, start, end)
}
