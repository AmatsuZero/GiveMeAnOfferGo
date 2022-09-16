package combination_sum_ii

import "sort"

// CombinationSum2 https://leetcode.cn/problems/combination-sum-ii/
func CombinationSum2(candidates []int, target int) (ans [][]int) {
	sort.Ints(candidates)
	var freq [][2]int
	for _, candidate := range candidates {
		if freq == nil || candidate != freq[len(freq)-1][0] {
			freq = append(freq, [2]int{candidate, 1})
		} else {
			freq[len(freq)-1][1]++
		}
	}

	var seq []int
	var dfs func(pos, rest int)
	dfs = func(pos, rest int) {
		if rest == 0 {
			ans = append(ans, append([]int{}, seq...))
			return
		}

		if pos == len(freq) || rest < freq[pos][0] {
			return
		}

		dfs(pos+1, rest)
		most := min(rest/freq[pos][0], freq[pos][1])
		for i := 1; i <= most; i++ {
			seq = append(seq, freq[pos][0])
			dfs(pos+1, rest-i*freq[pos][0])
		}
		seq = seq[:len(seq)-most]
	}

	dfs(0, target)
	return
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
