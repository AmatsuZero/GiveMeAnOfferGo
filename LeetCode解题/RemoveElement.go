package LeetCode解题

/*
Given an array and a value, remove all instances of that > value in place and return the new length.

The order of elements can be changed. It doesn't matter what you leave beyond the new length.
*/
func RemoveElement(arr []int, elem int) (result []int) {
	result = make([]int, len(arr), cap(arr))
	length := copy(result, arr)
	j := 0
	for i := 0; i < length; i++ {
		if result[i] == elem {
			continue
		}
		result[j] = result[i]
		j++
	}
	return result[:j+1]
}

/*
Given a sorted array, remove the duplicates in place such that > each element appear only once and return the new length.

Do not allocate extra space for another array, you must do this in place with constant memory.

For example, Given input array A = [1,1,2],

Your function should return length = 2, and A is now [1,2].
*/
func RemoveDuplicates(arr []int) (result []int) {
	result = make([]int, len(arr), cap(arr))
	if len(arr) == 0 {
		return
	}
	j, n := 0, copy(result, arr)
	for i := 1; i < n; i++ {
		if result[j] != result[i] {
			j++
			result[j] = result[i]
		}
	}
	return result[:j+1]
}

/*
Follow up for "Remove Duplicates": What if duplicates are allowed at most twice?

For example, Given sorted array A = [1,1,1,2,2,3],

Your function should return length = 5, and A is now [1,1,2,2,3].
*/
func RemoveSomeDuplicates(arr []int, repeatCount int) (result []int) {
	result = make([]int, len(arr), cap(arr))
	if len(arr) == 0 {
		return
	}
	j, n, count := 0, copy(result, arr), 0
	for i := 1; i < n; i++ {
		if result[j] != result[i] {
			j++
			result[j] = result[i]
			count = 0
		} else {
			count++
			if count < repeatCount {
				j++
				result[j] = result[i]
			}
		}
	}
	return result[:j+1]
}
