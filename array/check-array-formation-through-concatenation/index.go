package check_array_formation_through_concatenation

// CanFormArray https://leetcode.cn/problems/check-array-formation-through-concatenation/
func CanFormArray(arr []int, pieces [][]int) bool {
	index := make(map[int]int, len(pieces))

	for i, piece := range pieces {
		index[piece[0]] = i
	}

	for i := 0; i < len(arr); {
		j, ok := index[arr[i]]
		if !ok {
			return false
		}

		for _, x := range pieces[j] {
			if arr[i] != x {
				return false
			}
			i++
		}
	}

	return true
}
