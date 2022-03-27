package houserobber

// https://leetcode-cn.com/problems/house-robber/
func Rob(nums []int) int {
	n := len(nums)
	if n == 0 {
		return 0
	}
	curMax, preMax := 0, 0
	for i := 0; i < n; i++ {
		tmp := curMax
		if nums[i]+preMax > curMax {
			curMax = nums[i] + preMax
		}
		preMax = tmp
	}
	return curMax
}
