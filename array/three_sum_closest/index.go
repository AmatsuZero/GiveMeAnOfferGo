package three_sum_closest

import (
	"math"
	"sort"
)

// ThreeSumClosest https://leetcode.cn/problems/3sum-closest/
func ThreeSumClosest(nums []int, target int) int {
	sort.Ints(nums)
	n, ans := len(nums), math.MaxInt32
	update := func(cur int) {
		if abs(cur-target) < abs(ans-target) {
			ans = cur
		}
	}
	for i := 0; i < n; i++ {
		if i > 0 && nums[i] == nums[i-1] {
			continue
		}
		j, k := i+1, n-1
		for j < k {
			sum := nums[i] + nums[j] + nums[k]
			if sum == target {
				return target
			}
			update(sum)
			if sum > target {
				k0 := k - 1
				for j < k0 && nums[k0] == nums[k] {
					k0 -= 1
				}
				k = k0
			} else {
				j0 := j + 1
				for j0 < k && nums[j0] == nums[j] {
					j0 += 1
				}
				j = j0
			}
		}
	}
	return ans
}

func abs(x int) int {
	if x < 0 {
		x = -1 * x
	}
	return x
}
