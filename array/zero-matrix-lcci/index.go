package zero_matrix_lcci

// SetZeroes https://leetcode.cn/problems/zero-matrix-lcci/
func SetZeroes(matrix [][]int) {
	m, n := len(matrix), len(matrix[0])
	type record struct {
		r, c int
	}

	visited := map[record]struct{}{}
	for i, nums := range matrix {
		for j, num := range nums {
			if num == 0 {
				visited[record{i, j}] = struct{}{}
			}
		}
	}

	for rec, _ := range visited {
		matrix[rec.r] = make([]int, n)
		for i := 0; i < m; i++ {
			matrix[i][rec.c] = 0
		}
	}
}
