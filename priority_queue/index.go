package priority_queue

import (
	"container/heap"
	"sort"
)

var factors = []int{3, 5, 7}

type hp struct {
	sort.IntSlice
}

func (h *hp) Push(v interface{}) {
	h.IntSlice = append(h.IntSlice, v.(int))
}

func (h *hp) Pop() interface{} {
	n := len(h.IntSlice) - 1
	last := h.IntSlice[n]
	h.IntSlice = h.IntSlice[:n]
	return last
}

// GetKthMagicNumber https://leetcode.cn/problems/get-kth-magic-number-lcci/
func GetKthMagicNumber(k int) int {
	h := &hp{sort.IntSlice{1}}
	seen := map[int]struct{}{1: {}}
	for i := 1; ; i++ {
		x := heap.Pop(h).(int)
		if i == k {
			return x
		}
		for _, f := range factors {
			next := x * f
			if _, ok := seen[next]; !ok {
				heap.Push(h, next)
				seen[next] = struct{}{}
			}
		}
	}
}
