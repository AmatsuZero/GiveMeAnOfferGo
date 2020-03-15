package array

/*
给你一个长度为 n 的整数数组 nums，其中 n > 1，返回输出数组 output ，
其中  output[i] 等于 nums 中除 nums[i] 之外其余各元素的乘积。
*/
func ProductExceptSelf(nums []int) []int {
	length := len(nums)
	answer := make([]int, length)
	answer[0] = 1
	for i := 1; i < length; i++ {
		answer[i] = nums[i-1] * answer[i-1]
	}
	R := 1
	for i := length - 1; i >= 0; i-- {
		answer[i] = answer[i] * R
		R *= nums[i]
	}
	return answer
}
