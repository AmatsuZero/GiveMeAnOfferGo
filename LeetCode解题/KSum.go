package LeetCode解题

import (
	"math"
	"sort"
)

/*
Given an array of integers, find two numbers such that they add up to a specific target number. The function twoSum should return indices of the two numbers such that they add up to the target, where index1 must be less than index2 Please note that your returned answers (both index1 and index2) are not zero-based.

You may assume that each input would have exactly one solution.

Input: numbers={2, 7, 11, 15}, target=9 Output: index1=1, index2=2
*/
func TwoSum(numbers []int, target int) (index1, index2 int) {
	if len(numbers) <= 1 {
		return -1, -1
	}
	sums := make(map[int]int) // res - index
	for i, lhs := range numbers {
		if index, ok := sums[target-lhs]; ok {
			index1, index2 = i+1, index+1
			if index1 > index2 {
				index1, index2 = index2, index1
			}
			return
		}
		sums[lhs] = i
	}
	return -1, -1
}

/*
Given an array S of n integers, are there elements a, b, c in S such that a + b + c = 0? Find all unique triplets in the array which gives the sum of zero.

Note: Elements in a triplet (a,b,c) must be in non-descending order. (ie, a ≤ b ≤ c) The solution set must not contain duplicate triplets.
*/
func ThreeSum(nums []int) [][]int {
	result := make([][]int, 0)
	if len(nums) <= 2 { // 边际条件检查
		return result
	}
	// 先进行排序，按升序排列
	sort.Slice(nums, func(i, j int) bool {
		return nums[i] < nums[j]
	})
	for i := 0; i < len(nums)-2; i++ {
		j, k := i+1, len(nums)-1
		for j < k { // 创建一个栈的临时变量，来保存满足条件的每一组结果
			currStack := make([]int, 0)
			if nums[i]+nums[j]+nums[k] == 0 {
				currStack = append(currStack, nums[i])
				currStack = append(currStack, nums[j])
				currStack = append(currStack, nums[k])
				result = append(result, currStack)
				j++
				k--
				// 下面两个循环，用来去掉重复的解
				for j < k && nums[j-1] == nums[j] {
					j++
				}
				for j < k && nums[k] == nums[k+1] {
					k--
				}
			} else if nums[i]+nums[j]+nums[k] < 0 { // 三者之和小于期望值，指针向前
				j++
			} else {
				k--
			}
			for i < len(nums)-1 && nums[i] == nums[i+1] { // 此循环同样用来去掉重复解
				i++
			}
		}
	}
	return result
}

/*
Given an array S of n integers, find three integers in S such that the sum is closest to a given number, target.
Return the sum of the three integers.
You man assume that each input would have exactly one solution
*/
func ThreeSumCloset(nums []int, target int) (ret int) {
	if len(nums) <= 2 {
		return -1
	}
	distance := math.MaxInt64
	sort.Slice(nums, func(i, j int) bool {
		return nums[i] < nums[j]
	})
	for i := 0; i < len(nums)-2; i++ {
		j, k := i+1, len(nums)-1
		for j < k {
			tmpVal := nums[i] + nums[j] + nums[k]
			tmpDistance := 0
			if tmpVal < target {
				tmpDistance = target - tmpDistance
				if tmpDistance < distance {
					distance = tmpDistance
					ret = nums[i] + nums[j] + nums[k]
				}
				j++
			} else if tmpVal > target {
				tmpDistance = tmpVal
				if tmpDistance < distance {
					distance = tmpDistance
					ret = nums[i] + nums[j] + nums[k]
				}
				k--
			} else { // 在这里，三数之和等于0意味着距离最短
				return nums[i] + nums[j] + nums[k]
			}
		}
	}
	return
}
