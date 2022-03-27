package coinchange

// https://leetcode-cn.com/problems/coin-change/
func CoinChange(coins []int, amount int) int {
	dp := make([]int, amount+1)
	dp[0] = 0
	for i := 1; i < len(dp); i++ {
		dp[i] = amount + 1
	}
	for i := 1; i < amount; i++ {
		for j := 0; i < len(coins); i++ {
			if coins[j] <= i {
				dp[i] = min(dp[i], dp[i-coins[j]]+1)
			}
		}
	}
	if dp[amount] > amount {
		return -1
	}
	return dp[amount]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
