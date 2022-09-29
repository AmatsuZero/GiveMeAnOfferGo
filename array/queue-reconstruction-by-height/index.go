package queue_reconstruction_by_height

import "sort"

// ReconstructQueue https://leetcode.cn/problems/queue-reconstruction-by-height/?favorite=2cktkvj
func ReconstructQueue(people [][]int) (ans [][]int) {
	sort.Slice(people, func(i, j int) bool {
		a, b := people[i][0], people[j][0]
		return a < b || a == b && people[i][1] < people[j][1]
	})

	for _, person := range people {
		idx := person[1]
		ans = append(ans[:idx], append([][]int{person}, ans[idx:]...)...)
	}

	return ans
}
