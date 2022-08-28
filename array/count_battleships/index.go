package count_battleships

func CountBattleships(board [][]byte) (ans int) {
	m, n := len(board), len(board[0])

	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if board[i][j] == '.' {
				continue
			}

			originJ := j // 记录原始横向扫描的位置
			ans += 1

			// 水平扫描
			for ; j < n-1 && board[i][j] == 'X'; j++ {
				board[i][j] = '.'
			}

			// 垂直扫描
			for start := i + 1; start < m && board[start][originJ] == 'X'; start++ {
				board[start][originJ] = '.'
			}
		}
	}
	return
}
