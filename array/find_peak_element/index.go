package find_peak_element

func FindPeakElement(nums []int) int {
	if len(nums) == 1 {
		return 0
	} else if len(nums) == 2 {
		ans := 0
		if nums[ans] < nums[1] {
			ans = 1
		}
		return ans
	}

	maxNum := 0
	for i := 1; i < len(nums)-1; i++ {
		if nums[maxNum] < nums[i] {
			maxNum = i
		}
		if nums[i] > nums[i-1] && nums[i] > nums[i+1] {
			return i
		}
	}

	if nums[len(nums)-1] > nums[maxNum] {
		maxNum = len(nums) - 1
	}

	return maxNum
}
