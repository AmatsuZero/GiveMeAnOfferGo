// https://leetcode-cn.com/problems/two-sum/description/

package array

func TwoSum(nums []int, target int) []int {
	result := make([]int, 2)
	if len(nums) == 0 {
		return result
	}
	cachedResult := map[int]int{}
	for i, num := range nums {
		left := target - num
		if has, ok := cachedResult[left]; ok {
			result = []int{has, i}
		} else {
			cachedResult[num] = i
		}
	}
	return result
}
