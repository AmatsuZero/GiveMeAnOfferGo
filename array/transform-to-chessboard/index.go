package transform_to_chessboard

// MovesToChessboard https://leetcode.cn/problems/transform-to-chessboard/
func MovesToChessboard(board [][]int) (ans int) {
	isSameRow := func(a, b int) bool {
		for i, lhs := range board[a] {
			if lhs != board[b][i] {
				return false
			}
		}
		return true
	}

	isOppoRow := func(a, b int) bool {
		for i, lhs := range board[a] {
			if lhs == board[b][i] {
				return false
			}
		}
		return true
	}

	n := len(board)
	// 一样的列计数
	sameColCount := 0
	// 位置错误的列计数
	errorColCount := 0

	for i, num := range board[0] {
		if num == board[0][0] {
			sameColCount++
			// 与第一列相同，奇数下标的列
			if i%2 == 1 {
				errorColCount++
			}
		} else {
			if i%2 == 0 { // 与第一列不同，偶数下标的列
				errorColCount++
			}
		}
	}

	// n为偶数，same应当是一半
	if n%2 == 0 {
		if sameColCount != n/2 {
			return -1
		}
	} else { // n为奇数，same而应当为一半或者一半+1
		if sameColCount != n/2 && sameColCount != n/2+1 {
			return -1
		}
	}

	//一样的行计数
	sameRowCount := 0
	//位置错误的行计数
	errorRowCount := 0
	for i, _ := range board {
		if isSameRow(0, i) {
			sameRowCount++
			//与第一行相同，奇数下标的行
			if i%2 == 1 {
				errorRowCount++
			}
		} else if isOppoRow(0, i) {
			//与第一行不同，偶数下标的行
			if i%2 == 0 {
				errorRowCount++
			}
		} else {
			return -1
		}
	}

	//n为偶数，same应当是一半
	if n%2 == 0 {
		if sameRowCount != n/2 {
			return -1
		}
	} else {
		//n为奇数，same应当是一半或一半+1
		if sameRowCount != n/2 && sameRowCount != n/2+1 {
			return -1
		}
	}

	if n%2 == 0 {
		//偶数取小
		if n-errorColCount < errorColCount {
			errorColCount = n - errorColCount
		}
		if n-errorRowCount < errorRowCount {
			errorRowCount = n - errorRowCount
		}
	} else {
		//奇数看看列头和行头正不正确
		//第一列应该放在第二列，错的列其实都是对的
		if sameColCount*2 < n {
			errorColCount = n - errorColCount
		}
		//第一行应该放在第二行，错的行其实都是对的
		if sameRowCount*2 < n {
			errorRowCount = n - errorRowCount
		}
	}

	return (errorColCount + errorRowCount) / 2
}
