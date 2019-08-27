package array

func RemoveElement(nums []int, val int) int {
	pre := 0
	for i := 0; i < len(nums); i++ {
		if nums[i] != val {
			nums[pre] = nums[i]
			pre += 1
		}
	}
	return pre
}
