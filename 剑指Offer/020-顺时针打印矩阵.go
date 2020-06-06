package 剑指Offer

/*
题目：输入一个矩阵，按照从外向里以顺时针的顺序依次打印出每一个数字。例如：如果输入如下矩阵：

则依次打印出数字1、2、3、4、8、12、16、15、14、13、9、5、6、7、11、10。
*/
func EnumerateMatrixClockWisely(numbers [][]int, block func(val int)) {
	rows := len(numbers)
	if rows == 0 {
		return
	}
	columns := len(numbers[0])
	if columns == 0 {
		return
	}
	start := 0
	for columns > start*2 && rows > start*2 {
		enumerateMatrixInCircle(numbers, start, block)
		start++
	}
}

func enumerateMatrixInCircle(numbers [][]int, start int, block func(val int)) {
	rows := len(numbers)
	if rows == 0 {
		return
	}
	columns := len(numbers[0])
	if columns == 0 {
		return
	}
	endX, endY := columns-1-start, rows-1-start
	// 从左向右遍历
	for i := start; i <= endX; i++ {
		block(numbers[start][i])
	}
	// 从下到上遍历
	if start < endY {
		for i := start + 1; i <= endY; i++ {
			block(numbers[i][endX])
		}
	}
	// 从右到左遍历
	if start < endX && start < endY {
		for i := endX - 1; i >= start; i-- {
			block(numbers[endY][i])
		}
	}
	// 从下到上遍历
	if start < endX && start < endY-1 {
		for i := endY - 1; i >= start+1; i-- {
			block(numbers[i][start])
		}
	}
}
