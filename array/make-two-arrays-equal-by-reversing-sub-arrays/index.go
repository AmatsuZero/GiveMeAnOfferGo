package make_two_arrays_equal_by_reversing_sub_arrays

import "sort"

func CanBeEqual(target []int, arr []int) bool {
	if len(target) != len(arr) {
		return false
	}
	isEqual := func(lhs, rhs []int) bool {
		for i, v := range lhs {
			if rhs[i] != v {
				return false
			}
		}
		return true
	}

	if isEqual(target, arr) {
		return true
	}
	sort.Ints(target)
	sort.Ints(arr)
	return isEqual(target, arr)
}
