package global_and_local_inversions

// IsIdealPermutation https://leetcode.cn/problems/global-and-local-inversions/
func IsIdealPermutation(nums []int) bool {
	n := len(nums)
	minSuf := nums[n-1]
	for i := n - 2; i > 0; i-- {
		if nums[i-1] > minSuf {
			return false
		}
		minSuf = min(minSuf, nums[i])
	}
	return true
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}
