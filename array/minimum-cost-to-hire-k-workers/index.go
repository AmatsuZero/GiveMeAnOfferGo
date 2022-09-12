package minimum_cost_to_hire_k_workers

import (
	"container/heap"
	"math"
	"sort"
)

type hp struct {
	sort.IntSlice
}

func (h *hp) Less(i, j int) bool {
	return h.IntSlice[i] > h.IntSlice[j]
}

func (h *hp) Push(v interface{}) {
	h.IntSlice = append(h.IntSlice, v.(int))
}

func (h *hp) Pop() interface{} {
	a := h.IntSlice
	v := a[len(a)-1]
	h.IntSlice = a[:len(a)-1]
	return v
}

// MinCostToHireWorkers https://leetcode.cn/problems/minimum-cost-to-hire-k-workers/
func MinCostToHireWorkers(quality, wage []int, k int) float64 {
	n := len(quality)
	h := make([]int, n)
	for i := range h {
		h[i] = i
	}
	sort.Slice(h, func(i, j int) bool {
		a, b := h[i], h[j]
		return quality[a]*wage[b] > quality[b]*wage[a]
	})

	totalQ := 0
	q := hp{}
	for i := 0; i < k-1; i++ {
		totalQ += quality[h[i]]
		heap.Push(&q, quality[h[i]])
	}

	ans := 1e9
	for i := k - 1; i < n; i++ {
		idx := h[i]
		totalQ += quality[idx]
		heap.Push(&q, quality[idx])
		ans = math.Min(ans, float64(wage[idx])/float64(quality[idx])*float64(totalQ))
		totalQ -= heap.Pop(&q).(int)
	}

	return ans
}
