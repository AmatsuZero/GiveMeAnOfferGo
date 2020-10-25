package LeetCode解题

/*
给你一个由 '1'（陆地）和 '0'（水）组成的的二维网格，请你计算网格中岛屿的数量。

岛屿总是被水包围，并且每座岛屿只能由水平方向或竖直方向上相邻的陆地连接形成。

此外，你可以假设该网格的四条边均被水包围。



示例 1:

输入:
[
['1','1','1','1','0'],
['1','1','0','1','0'],
['1','1','0','0','0'],
['0','0','0','0','0']
]
输出: 1
示例 2:

输入:
[
['1','1','0','0','0'],
['1','1','0','0','0'],
['0','0','1','0','0'],
['0','0','0','1','1']
]
输出: 3
解释: 每座岛屿只能由水平和/或竖直方向上相邻的陆地连接而成。

作者：力扣 (LeetCode)
链接：https://leetcode-cn.com/leetbook/read/queue-stack/kbcqv/
来源：力扣（LeetCode）
著作权归作者所有。商业转载请联系作者获得授权，非商业转载请注明出处。
*/
func NumIslands(grid [][]byte) int {
	if len(grid) == 0 {
		return 0
	}
	if len(grid) > 10 {
		return numIslandsBFS(grid)
	}
	count := 0
	for x := 0; x < len(grid); x++ {
		for y := 0; y < len(grid[x]); y++ {
			if grid[x][y] == '1' {
				count++
				markAsWater(x, y, grid)
			}
		}
	}
	return count
}

func numIslandsBFS(grid [][]byte) (count int) {
	if len(grid) == 0 {
		return
	}
	height := len(grid)        // 网格的高
	queue := make([][2]int, 0) // 临时存放待探索的网格坐标
	for rowIndex, rows := range grid {
		width := len(rows) // 网格的宽
		for columnIndex, box := range rows {
			if box == '0' { // 是水，返回v
				continue
			}
			count += 1
			grid[rowIndex][columnIndex] = '0'                    // 改为海洋，跳过探索
			queue = append(queue, [2]int{rowIndex, columnIndex}) // 将陆地坐标加入队列
			//以当前队列中的陆地为出发点，对周围大陆进行广度优先遍历
			for len(queue) > 0 {
				var curr [2]int
				curr, queue = queue[0], queue[1:]                      // 当前网格
				if curr[0]-1 >= 0 && grid[curr[0]-1][curr[1]] == '1' { //探索上面网格
					grid[curr[0]-1][curr[1]] = '0'
					queue = append(queue, [2]int{curr[0] - 1, curr[1]})
				}
				if curr[0]+1 < height && grid[curr[0]+1][curr[1]] == '1' { //探索下面网格
					grid[curr[0]+1][curr[1]] = '0'
					queue = append(queue, [2]int{curr[0] + 1, curr[1]})
				}
				if curr[1]-1 >= 0 && grid[curr[0]][curr[1]-1] == '1' { //探索左边网格
					grid[curr[0]][curr[1]-1] = '0'
					queue = append(queue, [2]int{curr[0], curr[1] - 1})
				}
				if curr[1]+1 < width && grid[curr[0]][curr[1]+1] == '1' { //探索右边网格
					grid[curr[0]][curr[1]+1] = '0'
					queue = append(queue, [2]int{curr[0], curr[1] + 1})
				}
			}
		}
	}
	return
}

func markAsWater(x, y int, grid [][]byte) {
	if x < 0 || y < 0 || x >= len(grid) || y >= len(grid[x]) || grid[x][y] == '0' {
		return
	}
	grid[x][y] = '0'
	markAsWater(x-1, y, grid)
	markAsWater(x+1, y, grid)
	markAsWater(x, y-1, grid)
	markAsWater(x, y+1, grid)
}
