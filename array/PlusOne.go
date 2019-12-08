package array
/*
给定一个由整数组成的非空数组所表示的非负整数，在该数的基础上加一。

最高位数字存放在数组的首位， 数组中每个元素只存储单个数字。

你可以假设除了整数 0 之外，这个整数不会以零开头。
*/
func PlusOne(digits []int) []int {
	carry := 1

	for i := len(digits) - 1; i >= 0; i-- {
		if carry == 0 {
			return digits
		}
		tmp := digits[i] + carry
		carry = tmp / 10
		digits[i] = tmp % 10
	}

	if carry != 0 {
		result := make([]int, len(digits) + 1)
		result[0] = 1
		return result
	}

	return digits
}