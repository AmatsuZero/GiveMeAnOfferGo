package array

import (
	"GiveMeAnOffer/leetcode/sort"
	"reflect"
	"testing"
)

func TestMerge(t *testing.T) {
	ans := sort.Merge([][]int{
		{1, 3}, {4, 5}, {8, 10},
		{2, 6}, {9, 12}, {15, 18},
	})
	if !reflect.DeepEqual(ans, [][]int{
		{1, 6}, {8, 12}, {15, 18},
	}) {
		t.Failed()
	}
}
