package partitionequalsubsetsum

// https://leetcode-cn.com/problems/partition-equal-subset-sum/
func CanPartition(nums []int) bool {
	sum := 0
	for _, v := range nums {
		sum += v
	}
	if sum%2 != 0 {
		return false
	}
	n, C, dp := len(nums), sum/2, make([]bool, sum/2+1)
	for i := 0; i < C; i++ {
		dp[i] = (nums[0] == i)
	}
	for i := 1; i < n; i++ {
		for j := C; j >= nums[i]; j-- {
			dp[j] = dp[j] || dp[j-nums[i]]
		}
	}
	return dp[C]
}
