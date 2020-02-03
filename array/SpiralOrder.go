package array

func SpiralOrder(matrix [][]int) []int {
	result := make([]int, 0)
	if len(matrix) == 0 {
		return result
	}
	r1, r2 := 0, len(matrix)-1    // 规定当前层的上下边界
	c1, c2 := 0, len(matrix[0])-1 // 规定当前层的左右边界
	for r1 <= r2 && c1 <= c2 {
		for c := c1; c <= c2; c++ {
			result = append(result, matrix[r1][c])
		}
		for r := r1 + 1; r <= r2; r++ {
			result = append(result, matrix[r][c2])
		}
		if r1 < r2 && c1 < c2 {
			for c := c2 - 1; c > c1; c-- {
				result = append(result, matrix[r2][c])
			}
			for r := r2; r > r1; r-- {
				result = append(result, matrix[r][c1])
			}
		}
		// 往内部进一层
		r1++
		r2--
		c1++
		c2--
	}
	return result
}
