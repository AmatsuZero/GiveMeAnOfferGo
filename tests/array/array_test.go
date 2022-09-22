package array

import (
	check_array_formation_through_concatenation "GiveMeAnOffer/array/check-array-formation-through-concatenation"
	"GiveMeAnOffer/array/count_battleships"
	crawler_log_folder "GiveMeAnOffer/array/crawler-log-folder"
	"GiveMeAnOffer/array/find_peak_element"
	"GiveMeAnOffer/array/make-two-arrays-equal-by-reversing-sub-arrays"
	partition_to_k_equal_sum_subsets "GiveMeAnOffer/array/partition-to-k-equal-sum-subsets"
	"GiveMeAnOffer/array/set_zeros"
	validatestacksequences "GiveMeAnOffer/array/validate-stack-sequences"
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

func TestSetZeros(t *testing.T) {
	matrix := [][]int{
		{1, 1, 1},
		{1, 0, 1},
		{1, 1, 1},
	}
	set_zeros.SetZeroes(matrix)
	if !reflect.DeepEqual(matrix, [][]int{
		{1, 0, 1},
		{0, 0, 0},
		{1, 0, 1},
	}) {
		t.Fail()
	}

	matrix = [][]int{
		{0, 1, 2, 0},
		{3, 4, 5, 2},
		{1, 3, 1, 5},
	}
	set_zeros.SetZeroes(matrix)
	if !reflect.DeepEqual(matrix, [][]int{
		{0, 0, 0, 0},
		{0, 4, 5, 0},
		{0, 3, 1, 0},
	}) {
		t.Fail()
	}

	matrix = [][]int{
		{0, 0, 0, 5},
		{4, 3, 1, 4},
		{0, 1, 1, 4},
		{1, 2, 1, 3},
		{0, 0, 1, 1},
	}
	set_zeros.SetZeroes(matrix)
	if !reflect.DeepEqual(matrix, [][]int{
		{0, 0, 0, 0},
		{0, 0, 0, 4},
		{0, 0, 0, 0},
		{0, 0, 0, 3},
		{0, 0, 0, 0},
	}) {
		t.Fail()
	}
}

func TestCanEqual(t *testing.T) {
	if !make_two_arrays_equal_by_reversing_sub_arrays.CanBeEqual([]int{1, 2, 3, 4}, []int{2, 4, 1, 3}) {
		t.Fail()
	}

	if !make_two_arrays_equal_by_reversing_sub_arrays.CanBeEqual([]int{7}, []int{7}) {
		t.Fail()
	}

	if make_two_arrays_equal_by_reversing_sub_arrays.CanBeEqual([]int{3, 7, 9}, []int{3, 7, 11}) {
		t.Fail()
	}
}

func TestFindPeakElement(t *testing.T) {
	ans := find_peak_element.FindPeakElement([]int{3, 2, 1})
	if ans != 0 {
		t.Fail()
	}
}

func TestCountBattleShips(t *testing.T) {
	ans := count_battleships.CountBattleships([][]byte{
		{'X', '.', '.', 'X'},
		{'.', '.', '.', 'X'},
		{'.', '.', '.', 'X'},
	})
	if ans != 2 {
		t.Fail()
	}

	ans = count_battleships.CountBattleships([][]byte{
		{'.'},
	})
	if ans != 0 {
		t.Fail()
	}

	ans = count_battleships.CountBattleships([][]byte{
		{'X', 'X', 'X'},
	})
	if ans != 1 {
		t.Fail()
	}
}

func TestValidStackSequence(t *testing.T) {
	ans := validatestacksequences.ValidateStackSequences([]int{1, 2, 3, 4, 5}, []int{4, 5, 3, 2, 1})
	if !ans {
		t.Fail()
	}

	ans = validatestacksequences.ValidateStackSequences([]int{1, 2, 3, 4, 5}, []int{4, 3, 5, 1, 2})
	if ans {
		t.Fail()
	}

	ans = validatestacksequences.ValidateStackSequences([]int{2, 1, 0}, []int{1, 2, 0})
	if !ans {
		t.Fail()
	}
}

func TestMinOperation(t *testing.T) {
	ans := crawler_log_folder.MinOperations([]string{"d1/", "d2/", "../", "d21/", "./"})
	if ans != 2 {
		t.Fail()
	}
}

func TestCanPartitionSubsets(t *testing.T) {
	if !partition_to_k_equal_sum_subsets.CanPartitionKSubsets([]int{4, 3, 2, 3, 5, 2, 1}, 4) {
		t.Fail()
	}

	if partition_to_k_equal_sum_subsets.CanPartitionKSubsets([]int{1, 2, 3, 4}, 3) {
		t.Fail()
	}
}

func TestCanFormArray(t *testing.T) {
	if !check_array_formation_through_concatenation.CanFormArray([]int{15, 88}, [][]int{
		{88}, {15},
	}) {
		t.Fail()
	}

	if check_array_formation_through_concatenation.CanFormArray([]int{49, 18, 16}, [][]int{
		{16, 18, 49},
	}) {
		t.Fail()
	}

	if !check_array_formation_through_concatenation.CanFormArray([]int{91, 4, 64, 78}, [][]int{
		{78}, {4, 64}, {91},
	}) {
		t.Fail()
	}
}
