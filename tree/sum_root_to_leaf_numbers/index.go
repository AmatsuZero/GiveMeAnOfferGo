package sum_root_to_leaf_numbers

import (
	"GiveMeAnOffer/defines"
	"strconv"
)

// SumNumbers https://leetcode-cn.com/problems/sum-root-to-leaf-numbers/
func SumNumbers(root *defines.TreeNode) int {
	res, nums := 0, binaryTreeNums(root)
	for _, n := range nums {
		num, _ := strconv.Atoi(n)
		res += num
	}
	return res
}

func binaryTreeNums(root *defines.TreeNode) (res []string) {
	if root == nil {
		return []string{}
	}
	current := strconv.Itoa(root.Val)
	if root.Left == nil && root.Right == nil {
		return []string{current}
	}
	tmpLeft := binaryTreeNums(root.Left)
	for i := 0; i < len(tmpLeft); i++ {
		res = append(res, current+tmpLeft[i])
	}
	tmpRight := binaryTreeNums(root.Right)
	for i := 0; i < len(tmpRight); i++ {
		res = append(res, current+tmpRight[i])
	}
	return
}
