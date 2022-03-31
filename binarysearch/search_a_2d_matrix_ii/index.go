package searcha2dmatrixii

// https://leetcode-cn.com/problems/search-a-2d-matrix-ii/
func SearchMatrix(matrix [][]int, target int) bool {
	if len(matrix) == 0 {
		return false
	}
	row, col := 0, len(matrix[0])-1
	for col >= 0 && row <= len(matrix)-1 {
		if target == matrix[row][col] {
			return true
		} else if target > matrix[row][col] {
			row++
		} else {
			col--
		}

	}
	return false
}
