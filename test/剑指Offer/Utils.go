package 剑指Offer

import (
	. "github.com/AmatsuZero/GiveMeAnOfferGo/剑指Offer"
	"math/rand"
)

func RandomIntArray(len int) (output []int) {
	if len <= 0 {
		return
	}
	output = make([]int, len)
	for i := 0; i <= len-1; i++ {
		output[i] = rand.Intn(len)
	}
	return output
}

func RandomIntList(length int) (pHead *ListNode) {
	source := RandomIntArray(length)
	if len(source) == 0 {
		return
	}
	pHead = &ListNode{Val: source[0]}
	if len(source) == 1 {
		return
	}
	for i := 1; i < len(source); i++ {
		pHead.AddToTail(source[i])
	}
	return
}

func RandomInRange(min int, max int) int {
	if min >= max {
		panic("Check your Input!!!")
	}
	return min + rand.Intn(max-min+1) // a ≤ n ≤ b
}

func Partition(data *[]int, start, end int) (small int) {
	if len(*data) == 0 || start < 0 || end >= len(*data) {
		panic("Invalid parameters")
	}

	index := RandomInRange(start, end)
	small = start - 1
	(*data)[index], (*data)[end] = (*data)[end], (*data)[index]
	for index = start; index < end; index++ {
		if (*data)[index] < (*data)[end] {
			small++
			if small != index {
				(*data)[index], (*data)[small] = (*data)[small], (*data)[index]
			}
		}
	}
	small++
	(*data)[small], (*data)[end] = (*data)[end], (*data)[small]
	return
}

func QuickSort(data *[]int) {
	quickSort(data, 0, len(*data)-1)
}

func quickSort(data *[]int, start, end int) {
	if start == end {
		return
	}
	index := Partition(data, start, end)
	if index > start {
		quickSort(data, start, index-1)
	}
	if index < end {
		quickSort(data, index+1, end)
	}
}
