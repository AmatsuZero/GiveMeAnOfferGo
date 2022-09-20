package partition_to_k_equal_sum_subsets

import "sort"

// CanPartitionKSubsets https://leetcode.cn/problems/partition-to-k-equal-sum-subsets/
func CanPartitionKSubsets(nums []int, k int) bool {
	all, n := 0, len(nums)
	for _, num := range nums {
		all += num
	}

	if all%k > 0 {
		return false
	}

	sort.Slice(nums, func(i, j int) bool {
		return nums[i] < nums[j]
	})

	cnt, avg := make([]int, k), all/k
	var dfs func(u int) bool
	dfs = func(u int) bool {
		if u == n {
			for _, item := range cnt {
				if item != avg {
					return false
				}
			}
			return true
		}

		for i := 0; i < k; i++ {
			if cnt[i]+nums[u] > avg { // 第 i 组无法添加
				continue
			}
			if i > 0 && cnt[i] == cnt[i-1] { // 第 i 组和第 i-1 组相同，且第 i-1 组已经枚举过
				continue
			}
			cnt[i] += nums[u]
			if dfs(u + 1) {
				return true
			}
			cnt[i] -= nums[u]
		}

		return false
	}

	return dfs(0)
}
