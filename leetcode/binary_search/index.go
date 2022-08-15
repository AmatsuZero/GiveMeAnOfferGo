package binary_search

func SearchInsert(nums []int, target int) int {
	left, right := 0, len(nums)-1
	for left <= right {
		mid := (left + right) / 2
		if nums[mid] >= target {
			if mid == 0 || nums[mid-1] < target {
				return mid
			}
			return mid - 1
		} else {
			left = mid + 1
		}
	}
	return len(nums)
}

func PeakIndexInMountainArray(nums []int) int {
	left, right := 1, len(nums)-2
	for left <= right {
		mid := (left + right) / 2
		if nums[mid] > nums[mid+1] && nums[mid] > nums[mid-1] {
			return mid
		}
		if nums[mid] > nums[mid-1] {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return -1
}
