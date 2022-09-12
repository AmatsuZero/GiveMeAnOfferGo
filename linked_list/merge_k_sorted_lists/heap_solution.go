package merge_k_sorted_lists

import (
	"GiveMeAnOffer/defines"
	"container/heap"
)

type ListNode = defines.ListNode

type hp struct {
	arr []*ListNode
}

func solution(lists []*ListNode) *ListNode {
	ans := &ListNode{}
	p := ans
	h := &hp{}
	for _, l := range lists {
		if l != nil {
			heap.Push(h, l)
		}
	}

	for h.Len() > 0 {
		n := heap.Pop(h).(*ListNode)
		p.Next = n
		p = p.Next
		if n.Next != nil {
			heap.Push(h, n.Next)
		}
	}

	return ans.Next
}

func (h *hp) Len() int {
	return len(h.arr)
}

func (h *hp) Swap(i, j int) {
	h.arr[i], h.arr[j] = h.arr[j], h.arr[i]
}

func (h *hp) Less(i, j int) bool {
	return h.arr[i].Val < h.arr[j].Val
}

func (h *hp) Push(x interface{}) {
	h.arr = append(h.arr, x.(*ListNode))
}

func (h *hp) Pop() (v interface{}) {
	idx := len(h.arr) - 1
	v, h.arr = h.arr[idx], h.arr[:idx]
	return
}
