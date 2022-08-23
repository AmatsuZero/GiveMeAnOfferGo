package set_zeros

type location struct {
	row, column int
}

// SetZeroes https://leetcode.cn/problems/set-matrix-zeroes/
func SetZeroes(matrix [][]int) {
	m, n := len(matrix), len(matrix[0])
	replaced := map[location]bool{}

	for i := 0; i < m; i++ {
		row := matrix[i]
		needReplace := false
		for j := 0; j < n; j++ {
			l := location{i, j}
			if !replaced[l] && row[j] == 0 {
				needReplace = true
				for c := 0; c < m; c++ {
					l = location{c, j}
					// 如果原来就是0，则不标记
					if matrix[c][j] != 0 {
						replaced[l] = true
					}
					matrix[c][j] = 0
				}
			}
		}
		if needReplace {
			matrix[i] = make([]int, n)
		}
	}
}
