package array

import (
	"GiveMeAnOffer/leetcode/backtrace"
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

func TestSubsets(t *testing.T) {
	ans := backtrace.Subsets([]int{1, 2})
	if len(ans) != 4 {
		t.Failed()
	}
}

func TestCombine(t *testing.T) {
	ans := backtrace.Combine(3, 2)
	if len(ans) != 3 {
		t.Failed()
	}
}
