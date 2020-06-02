package 剑指Offer

/*
“题目：在一个二维数组中，每一行都按照从左到右递增的顺序排序，每一列都按照从上到下递增的顺序排序。请完成一个函数，输入这样的一个二维数组和一个整数，判断数组中是否含有该整数。”
*/
func Find(array [][]int, target int) (found bool) {
	rows := len(array)
	if rows == 0 {
		return false
	}
	cols := len(array[0])
	if cols == 0 {
		return false
	}
	col, row := cols-1, 0
	for row < rows && col >= 0 {
		if array[row][col] == target {
			found = true
			break
		} else if array[row][col] > target { // 大于要查找的数字，向左查找
			col--
		} else { // 小于要查找的数字向下查找
			row++
		}
	}
	return found
}
