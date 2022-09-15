package task_scheduler

import (
	"container/heap"
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

// LeastInterval https://leetcode.cn/problems/task-scheduler/
func LeastInterval(tasks []byte, n int) int {
	ans := 0
	mp := map[byte]int{}
	for _, task := range tasks {
		mp[task]++
	}

	pq := &hp{}
	for _, i := range mp { //转化成大顶堆
		heap.Push(pq, i)
	}

	for pq.Len() > 0 {
		var arr []int //临时存储执行完，但剩余数量不为0的任务
		for i := 0; i <= n; i++ {
			if pq.Len() > 0 {
				top := heap.Pop(pq).(int) //挑数量最多的任务先执行
				if top > 1 {
					top--
					arr = append(arr, top)
				}
			} else {
				if len(arr) == 0 { //任务正好都执行完了，不用执行冷却了
					return ans
				}
			}
			ans++
		}
		for _, ele := range arr {
			heap.Push(pq, ele)
		}
	}

	return ans
}
