package LeetCode解题

import (
	"strconv"
	"strings"
)

type PascalTriangle [][]int

/*
Given numRows, generate the first numRows of Pascal's triangle.

For example, given numRows = 5, Return

[
     [1],
    [1,1],
   [1,2,1],
  [1,3,3,1],
 [1,4,6,4,1]
]
*/
func NewPascalTriangle(numRows int) (result PascalTriangle) {
	if numRows <= 0 {
		return make([][]int, 0)
	}
	result = make([][]int, numRows)
	for i := 0; i < numRows; i++ {
		result[i] = make([]int, i+1)
		result[i][0] = 1
		result[i][len(result[i])-1] = 1
		for k := 1; k < len(result[i])-1; k++ {
			result[i][k] = result[i-1][k-1] + result[i-1][k]
		}
	}
	return
}

func (t PascalTriangle) String() string {
	str := "\n[\n"
	for i, elm := range t {
		space := len(t) - i
		str += strings.Repeat(" ", space)
		str += "["
		for k, n := range elm {
			str += strconv.Itoa(n)
			if k != len(elm)-1 {
				str += ","
			}
		}
		str += "]\n"
	}
	return str + "]"
}

/*
Given an index k, return the kth row of the Pascal's triangle.

For example, given k = 3, Return [1,3,3,1].
*/
func (t PascalTriangle) GetRow(k uint) (result []int) {
	if t != nil && len(t) > 0 && int(k) < len(t) {
		return t[k]
	}
	result = make([]int, 1, k+1)
	result[0] = 1
	if k == 0 {
		return
	}
	for i := 0; i < int(k); i++ {
		result = append(result, 1)
		for j := len(result) - 2; j > 0; j-- { // 每一次循环，就得到一层的数据
			result[j] += result[j-1]
		}
	}
	return
}
