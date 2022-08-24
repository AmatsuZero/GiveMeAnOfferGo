package backtrace

func Subsets(nums []int) (ans [][]int) {
	if len(nums) == 0 {
		return
	}
	var subset []int
	var helper func(index int)
	helper = func(index int) {
		if index == len(nums) {
			dst := make([]int, len(subset))
			copy(dst, subset)
			ans = append(ans, dst)
		} else if index < len(nums) {
			helper(index + 1)
			subset = append(subset, nums[index])
			helper(index + 1)
			subset = subset[:len(subset)-1]
		}
	}
	helper(0)
	return
}

func Combine(n, k int) (result [][]int) {
	var helper func(i int)
	var combination []int
	helper = func(i int) {
		if len(combination) == k {
			dst := make([]int, k)
			copy(dst, combination)
			result = append(result, dst)
		} else if i <= n {
			helper(i + 1)
			combination = append(combination, i)
			helper(i + 1)
			combination = combination[:len(combination)-1]
		}
	}
	helper(1)
	return
}

func addDigits(num int) int {
	if num < 10 {
		return num
	}
	sum := 0
	for num/10 != 0 {
		sum += num % 10
		num /= 10
	}
	if num != 0 {
		sum += num
	}
	return addDigits(sum)
}
