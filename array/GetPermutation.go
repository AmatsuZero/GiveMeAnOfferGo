package array

import "strconv"

/*
给出集合 [1,2,3,…,n]，其所有元素共有 n! 种排列。

按大小顺序列出所有排列情况，并一一标记，当 n = 3 时, 所有排列如下：

"123"
"132"
"213"
"231"
"312"
"321"
给定 n 和 k，返回第 k 个排列。

说明：

给定 n 的范围是 [1, 9]。
给定 k 的范围是[1,  n!]。
示例 1:

输入: n = 3, k = 3
输出: "213"
示例 2:

输入: n = 4, k = 9
输出: "2314"
*/

func GetPermutation(n int, k int) string {
	visited := make([]bool, n)
	// 将 n! 种排列分为：n 组，每组有 (n - 1)! 种排列
	return recursive(n, factorial(n-1), k, visited)
}

func recursive(n int, f int, k int, visited []bool) string {
	offset := k % f     // 组内偏移量
	groupIndex := k / f // 第几组
	if offset > 0 {
		groupIndex += 1
	}
	i := 0
	for ; i < len(visited) && groupIndex > 0; i++ { // 在没有被访问的数字里找第groupIndex个数字
		if !visited[i] {
			groupIndex--
		}
	}
	visited[i-1] = true // 标记为已访问
	if n-1 > 0 {
		// offset = 0 时，则取第 i 组的第 f 个排列，否则取第 i 组的第 offset 个排列
		if offset == 0 {
			offset = f
		}
		return strconv.Itoa(i) + recursive(n-1, f/(n-1), offset, visited)
	} else {
		// 最后一位
		return strconv.Itoa(i)
	}
}

/**
 * 求 n!
 */
func factorial(n int) (res int) {
	res = 1
	for i := n; i > 1; i-- {
		res *= i
	}
	return
}
