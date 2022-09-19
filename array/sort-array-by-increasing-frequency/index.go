package sort_array_by_increasing_frequency

import "sort"

// FrequencySort https://leetcode.cn/problems/sort-array-by-increasing-frequency/
func FrequencySort(nums []int) []int {
	freq := map[int]int{}

	for _, num := range nums {
		freq[num]++
	}

	sort.Slice(nums, func(i, j int) bool {
		lhs, rhs := freq[nums[i]], freq[nums[j]]
		if lhs == rhs {
			return nums[i] < nums[j]
		}
		return lhs < rhs
	})

	return nums
}
