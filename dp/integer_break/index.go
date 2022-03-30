package integerbreak

// https://leetcode-cn.com/problems/integer-break/
func IntegerBreak(n int) int {
	dp := make([]int, n+1)
	dp[0], dp[1] = 1, 1
	for i := 1; i < n; i++ {
		for j := 0; j < i; j++ {
			dp[i] = max(dp[i], j*max(dp[i-j], i-j))
		}
	}
	return dp[n]
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
