package array

/*
编写一种算法，若M × N矩阵中某个元素为0，则将其所在的行与列清零。



示例 1：

输入：
[
  [1,1,1],
  [1,0,1],
  [1,1,1]
]
输出：
[
  [1,0,1],
  [0,0,0],
  [1,0,1]
]
示例 2：

输入：
[
  [0,1,2,0],
  [3,4,5,2],
  [1,3,1,5]
]
输出：
[
  [0,0,0,0],
  [0,4,5,0],
  [0,3,1,0]
]

作者：力扣 (LeetCode)
链接：https://leetcode-cn.com/leetbook/read/array-and-string/ciekh/
来源：力扣（LeetCode）
著作权归作者所有。商业转载请联系作者获得授权，非商业转载请注明出处。
*/
func SetZeros(matrix [][]int) {
	rLen, cLen := len(matrix), len(matrix[0])
	hMap, lMap := map[int]bool{}, map[int]bool{}
	for i := 0; i < rLen; i++ {
		for j := 0; j < cLen; j++ {
			if matrix[i][j] == 0 {
				hMap[i] = true
				lMap[j] = true
			}
		}
	}
	for i := 0; i < rLen; i++ {
		for j := 0; j < cLen; j++ {
			if hMap[i] || lMap[j] {
				matrix[i][j] = 0
			}
		}
	}
}
