package convert_sorted_array_to_binary_search_tree

import "GiveMeAnOffer/defines"

func SortedArrayToBST(nums []int) *defines.TreeNode {
	if len(nums) == 0 {
		return nil
	}
	halfLen := len(nums) / 2
	return &defines.TreeNode{
		Val:   nums[halfLen],
		Left:  SortedArrayToBST(nums[:halfLen]),
		Right: SortedArrayToBST(nums[halfLen+1:]),
	}
}
