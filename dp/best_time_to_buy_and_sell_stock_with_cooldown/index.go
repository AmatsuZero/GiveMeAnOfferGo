package besttimetobuyandsellstockwithcooldown

import "math"

// https://leetcode-cn.com/problems/best-time-to-buy-and-sell-stock-with-cooldown/
func MaxProfit(prices []int) int {
	if len(prices) <= 1 {
		return 0
	}
	buy := []int{-prices[0], max(-prices[0], -prices[1]), math.MinInt32}
	sell := []int{0, max(0, -prices[0]+prices[1]), 0}
	for i := 2; i < len(prices); i++ {
		sell[i%3] = max(sell[(i-1)%3], buy[(i-1)%3]+prices[i])
		buy[i%3] = max(buy[(i-1)%3], sell[(i-2)%3]-prices[i])
	}
	return sell[(len(prices)-1)%3]
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
