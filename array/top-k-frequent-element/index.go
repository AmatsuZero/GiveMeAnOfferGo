package top_k_frequent_element

import (
	"container/heap"
)

type hp struct {
	arr [][]int
}

func topKFrequent(nums []int, k int) []int {
	dict := map[int]int{}
	for _, num := range nums {
		dict[num] += 1
	}
	h := &hp{}
	for key, value := range dict {
		heap.Push(h, []int{key, value})
		if h.Len() > k {
			heap.Pop(h)
		}
	}
	return h.Keys()
}

func (h *hp) Less(i, j int) bool {
	return h.arr[i][1] < h.arr[j][1]
}

func (h *hp) Push(v interface{}) {
	h.arr = append(h.arr, v.([]int))
}

func (h *hp) Pop() interface{} {
	a := h.arr
	v := a[len(a)-1]
	h.arr = a[:len(a)-1]
	return v
}

func (h *hp) Len() int {
	return len(h.arr)
}

func (h *hp) Swap(i, j int) {
	h.arr[i], h.arr[j] = h.arr[j], h.arr[i]
}

func (h *hp) Keys() (ans []int) {
	for _, a := range h.arr {
		ans = append(ans, a[0])
	}
	return
}
