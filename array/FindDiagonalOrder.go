package array

/*

给定一个含有 M x N 个元素的矩阵（M 行，N 列），请以对角线遍历的顺序返回这个矩阵中的所有元素，对角线遍历如下图所示。



示例:

输入:
[
 [ 1, 2, 3 ],
 [ 4, 5, 6 ],
 [ 7, 8, 9 ]
]

输出:  [1,2,4,7,5,3,6,8,9]

作者：力扣 (LeetCode)
链接：https://leetcode-cn.com/leetbook/read/array-and-string/cuxq3/
来源：力扣（LeetCode）
著作权归作者所有。商业转载请联系作者获得授权，非商业转载请注明出处。
*/
func FindDiagonalOrder(matrix [][]int) []int {
	if len(matrix) == 0 {
		return []int{}
	}
	// 计算矩阵尺寸
	n, m := len(matrix), len(matrix[0])
	result := make([]int, n*m)
	k := 0
	intermediate := make([]int, 0)
	// 遍历第一行和最后一列中的所有元素，来涵盖所有可能的对角线
	for d := 0; d < n+m-1; d++ {
		// 每次处理一条对角线，就清空
		intermediate = intermediate[:0]
		// 找到本次遍历对角线的起点，
		r := d - m + 1
		if d < m {
			r = 0
		}
		c := m - 1
		if d < m {
			c = d
		}
		for r < n && c > -1 {
			intermediate = append(intermediate, matrix[r][c])
			r++
			c--
		}
		if d%2 == 0 { // 如果是偶数，就将中间数组反转
			for left, right := 0, len(intermediate)-1; left < right; left, right = left+1, right-1 {
				intermediate[left], intermediate[right] = intermediate[right], intermediate[left]
			}
		}
		for i := 0; i < len(intermediate); i++ {
			result[k] = intermediate[i]
			k++
		}
	}
	return result
}
