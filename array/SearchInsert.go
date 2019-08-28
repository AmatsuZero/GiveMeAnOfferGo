// https://leetcode-cn.com/problems/search-insert-position/description/
package array

func SearchInsert(nums []int, target int) int {
	length := len(nums)
	if length == 0 {
		return 0
	}

	if nums[length-1] < target {
		return length
	}

	left := 0
	right := length - 1

	for left < right {
		mid := left + (right-left)/2
		if nums[mid] == target {
			return mid
		} else if nums[mid] < target {
			left = mid + 1
		} else {
			right = mid
		}
	}

	return right
}
