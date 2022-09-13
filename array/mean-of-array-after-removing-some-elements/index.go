package mean_of_array_after_removing_some_elements

import "sort"

// TrimMean https://leetcode.cn/problems/mean-of-array-after-removing-some-elements/
func TrimMean(arr []int) float64 {
	sort.Ints(arr)
	// 移除头部和尾部5%的元素
	n := int(float64(len(arr)) * 0.05)
	arr = arr[n : len(arr)-n]
	var sum float64
	for _, num := range arr {
		sum += float64(num)
	}
	return sum / float64(len(arr))
}
