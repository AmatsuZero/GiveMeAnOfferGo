package maximum_ascending_subarray_sum

// MaxAscendingSum https://leetcode.cn/problems/maximum-ascending-subarray-sum/
func MaxAscendingSum(nums []int) int {
	sum, ans := nums[0], 0
	for i := 1; i < len(nums); i++ {
		if nums[i-1] < nums[i] {
			sum += nums[i]
		} else {
			if sum > ans {
				ans = sum
			}
			sum = nums[i]
		}
	}

	if sum > ans {
		ans = sum
	}
	return ans
}
