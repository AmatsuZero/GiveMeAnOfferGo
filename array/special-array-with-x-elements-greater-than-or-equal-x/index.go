package special_array_with_x_elements_greater_than_or_equal_x

// SpecialArray https://leetcode.cn/problems/special-array-with-x-elements-greater-than-or-equal-x/
func SpecialArray(nums []int) int {
	n := len(nums)
	ans := n
	for ; ans >= 0; ans-- {
		valid := 0
		for _, num := range nums {
			if num >= ans {
				valid += 1
			}
		}
		if valid == ans {
			return ans
		}
	}
	return -1
}
