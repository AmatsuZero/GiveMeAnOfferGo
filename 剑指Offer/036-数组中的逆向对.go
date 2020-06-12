package 剑指Offer

/*
题目：在数组中的两个数字如果前面一个数字大于后面的数字，则这两个数字组成一个逆序对。输入一个数组，求出这个数组中的逆序对的总数
*/
func InversePairs(nums []int) int {
	return mergeOuter(nums, 0, len(nums)-1, 0)
}

func mergeOuter(nums []int, left, right, res int) int {
	if left >= right {
		return res
	}
	mid := (left + right) / 2
	res = mergeOuter(nums, left, mid, res)
	res = mergeOuter(nums, mid+1, right, res)
	return merge(nums, left, mid+1, right, res)
}

func merge(nums []int, left, mid, right, res int) int {
	leftArray := make([]int, mid-left)
	rightArray := make([]int, right-mid+1)

	for i := left; i < mid; i++ {
		leftArray[i-left] = nums[i]
	}

	for i := mid; i <= right; i++ {
		rightArray[i-mid] = nums[i]
	}

	k := left
	l, r := 0, 0

	for l < len(leftArray) && r < len(rightArray) {
		if leftArray[l] > rightArray[r] {
			nums[k] = rightArray[r]
			k++
			r++
			res = (res + mid - left - l) % 1000000007
		} else {
			nums[k] = leftArray[l]
			k++
			l++
		}
	}

	for l < len(leftArray) {
		nums[k] = leftArray[l]
		k++
		l++
	}

	for r < len(rightArray) {
		nums[k] = rightArray[r]
		k++
		r++
	}
	return res
}
