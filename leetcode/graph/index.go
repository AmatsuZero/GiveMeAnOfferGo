package graph

func MaxAreaOfIslands(grid [][]int) (maxArea int) {
	var (
		rows    = len(grid)
		cols    = len(grid[0])
		visited = make([][]bool, rows)
	)
	for i := 0; i < rows; i++ {
		visited[i] = make([]bool, cols)
	}

	getArea := func(i, j int) (area int) {
		type position struct {
			x, y int
		}

		var queue []position
		queue = append(queue, position{i, j})
		visited[i][j] = true

		dirs := []position{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
		for pos := queue[len(queue)-1]; len(queue) > 0; queue = queue[:len(queue)-1] {
			area += 1
			for _, dir := range dirs {
				r := pos.x + dir.x
				c := pos.y + dir.y
				if r >= 0 && r < rows &&
					c >= 0 && c < cols &&
					grid[r][c] == 1 && !visited[r][c] {
					queue = append(queue, position{r, c})
					visited[r][c] = true
				}
			}
		}
		return
	}

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if grid[i][j] == 1 && !visited[i][j] {
				area := getArea(i, j)
				if area > maxArea {
					maxArea = area
				}
			}
		}
	}
	return maxArea
}
