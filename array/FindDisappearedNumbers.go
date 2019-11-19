package array

/*
找到所有数组中消失的数字
Category	Difficulty	Likes	Dislikes
algorithms	Easy (51.79%)	229	-
Tags
Companies
给定一个范围在  1 ≤ a[i] ≤ n ( n = 数组大小 ) 的 整型数组，数组中的元素一些出现了两次，另一些只出现一次。

找到所有在 [1, n] 范围之间没有出现在数组中的数字。

您能在不使用额外空间且时间复杂度为O(n)的情况下完成这个任务吗? 你可以假定返回的数组不算在额外空间内。

示例:

输入:
[4,3,2,7,8,2,3,1]

输出:
[5,6]
*/
func FindDisappearedNumbers(nums []int) (ans []int) {
	for i := 0; i < len(nums); i++ {
		num := nums[i]
		if num < 0 {
			num = -num
		}
		index := num - 1
		if nums[index] > 0 {
			nums[index] = -nums[index]
		}
	}

	for i := 0; i < len(nums); i++ {
		if nums[i] > 0 {
			ans = append(ans, i+1)
		}
	}
	return
}
