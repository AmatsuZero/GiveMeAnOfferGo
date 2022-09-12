package kth_largest_element_in_an_array

import (
	"container/heap"
	"sort"
)

// FindKthLargest https://leetcode.cn/problems/kth-largest-element-in-an-array/
func FindKthLargest(nums []int, k int) int {
	h := &hp{}
	for _, num := range nums {
		heap.Push(h, num)
		if h.Len() > k {
			heap.Pop(h)
		}
	}
	return h.IntSlice[0]
}

type hp struct {
	sort.IntSlice
}

func (h *hp) Less(i, j int) bool {
	return h.IntSlice[i] < h.IntSlice[j]
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
