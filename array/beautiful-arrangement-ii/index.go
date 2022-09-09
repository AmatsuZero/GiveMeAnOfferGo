package beautiful_arrangement_ii

// ConstructArray https://leetcode.cn/problems/beautiful-arrangement-ii/
func ConstructArray(n int, k int) []int {
	arr := make([]int, 0, n)
	for i := 1; i < n-k; i++ {
		arr = append(arr, i)
	}

	for i, j := n-k, n; i <= j; i++ {
		arr = append(arr, i)
		if i != j {
			arr = append(arr, j)
		}
		j--
	}

	return arr
}
