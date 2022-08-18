package maximum_equal_frequency

// MaxEqualFreq https://leetcode.cn/problems/maximum-equal-frequency/
func MaxEqualFreq(nums []int) (ans int) {
	freq := map[int]int{}
	count := map[int]int{}
	maxFreq := 0
	for i, num := range nums {
		if count[num] > 0 {
			freq[count[num]] -= 1
		}
		count[num] += 1
		maxFreq = max(maxFreq, count[num])
		freq[count[num]] += 1
		if maxFreq == 1 ||
			freq[maxFreq]*maxFreq+freq[maxFreq-1]*(maxFreq-1) == i+1 && freq[maxFreq] == 1 ||
			freq[maxFreq]*maxFreq+1 == i+1 && freq[1] == 1 {
			ans = max(ans, i+1)
		}
	}
	return
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
