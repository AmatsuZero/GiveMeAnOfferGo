package maximum_swap

import "strconv"

// MaximumSwap https://leetcode.cn/problems/maximum-swap/
func MaximumSwap(num int) int {
	index := make([]int, 10)
	n := []byte(strconv.Itoa(num))

	for i, c := range n {
		index[c-'0'] = i
	}

	for i, c := range n {
		for j := 9; j > int(c-'0'); j-- { // 数字比当前数字大
			if index[j] > i {
				n[i], n[index[j]] = n[index[j]], n[i]
				ans, _ := strconv.Atoi(string(n))
				return ans
			}
		}
	}

	return num
}
