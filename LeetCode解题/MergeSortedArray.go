package LeetCode解题

// https://leetcode-cn.com/problems/merge-sorted-array/

func MergeSortedArray(lhs, rhs []int, m, n int) {
	k := m + n - 1
	i, j := m-1, n-1
	for i >= 0 && j >= 0 {
		if lhs[i] > rhs[j] {
			lhs[k] = lhs[i]
			k--
			i--
		} else {
			lhs[k] = rhs[j]
			k--
			j--
		}
	}
	for i >= 0 {
		lhs[k] = lhs[i]
		k--
		i--
	}
	for j >= 0 {
		lhs[k] = rhs[j]
		k--
		j--
	}
}
