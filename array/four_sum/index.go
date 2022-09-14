package four_sum

import "sort"

// FourSum https://leetcode.cn/problems/4sum/
func FourSum(nums []int, target int) (ans [][]int) {
	sort.Ints(nums)

	var dfs func(low, sum int)
	var subArr []int
	numSize := len(nums)

	dfs = func(low, sum int) {
		if sum == target && len(subArr) == 4 {
			ans = append(ans, append([]int{}, subArr...))
			return
		}

		for i := low; i < numSize; i++ {
			if numSize-i < (4 - len(subArr)) { // 剪枝
				return
			}
			if i > low && nums[i] == nums[i-1] { // 去重
				continue
			}
			if i < numSize-1 && (sum+nums[i]+(3-len(subArr)))*nums[i+1] > target { //剪枝
				return
			}
			if i < numSize-1 && (sum+nums[i]+(3-len(subArr)))*nums[numSize-1] < target { //剪枝
				continue
			}

			subArr = append(subArr, nums[i])
			dfs(i+1, nums[i]+sum)
			subArr = subArr[:len(subArr)-1]
		}
	}
	if numSize < 4 {
		return nil
	}
	dfs(0, 0)
	return
}
