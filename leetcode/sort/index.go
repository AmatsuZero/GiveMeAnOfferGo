package sort

import (
	"sort"
)

func Merge(intervals [][]int) (ans [][]int) {
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})

	for i := 0; i < len(intervals); i += 1 {
		temp := []int{intervals[i][0], intervals[i][1]}
		j := i + 1
		for j < len(intervals) && intervals[j][0] <= temp[1] {
			if intervals[j][1] > temp[1] {
				temp[1] = intervals[j][1]
			}
			j += 1
		}
		ans = append(ans, temp)
		i = j
	}
	return
}
