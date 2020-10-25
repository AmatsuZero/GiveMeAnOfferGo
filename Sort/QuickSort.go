package Sort

import (
	"math/rand"
	"time"
)

func partitionDutchFlag(a []int, low, high, pivotIndex int) (smaller, larger int) {
	pivot := a[pivotIndex]
	smaller, larger = low, high
	equal := low
	for equal <= larger {
		if a[equal] < pivot {
			a[smaller], a[equal] = a[equal], a[smaller]
			smaller++
			equal++
		} else if a[equal] == pivot {
			equal++
		} else {
			a[equal], a[larger] = a[larger], a[equal]
			larger -= 1
		}
	}
	return
}

func QuickSortDutchFlag(slice []int, low, high int) {
	if low >= high {
		return
	}
	rand.Seed(time.Now().Unix())
	pivotIndex := rand.Intn(high-low) + low
	p, q := partitionDutchFlag(slice, low, high, pivotIndex)
	QuickSortDutchFlag(slice, low, p-1)
	QuickSortDutchFlag(slice, q+1, high)
}
