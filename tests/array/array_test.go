package array

import (
	"GiveMeAnOffer/array/count_battleships"
	crawler_log_folder "GiveMeAnOffer/array/crawler-log-folder"
	"GiveMeAnOffer/array/find_peak_element"
	"GiveMeAnOffer/array/make-two-arrays-equal-by-reversing-sub-arrays"
	"GiveMeAnOffer/array/set_zeros"
	validate_stack_sequences "GiveMeAnOffer/array/validate-stack-sequences"
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
	ans := validate_stack_sequences.ValidateStackSequences([]int{1, 2, 3, 4, 5}, []int{4, 5, 3, 2, 1})
	if !ans {
		t.Fail()
	}

	ans = validate_stack_sequences.ValidateStackSequences([]int{1, 2, 3, 4, 5}, []int{4, 3, 5, 1, 2})
	if ans {
		t.Fail()
	}

	ans = validate_stack_sequences.ValidateStackSequences([]int{2, 1, 0}, []int{1, 2, 0})
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
