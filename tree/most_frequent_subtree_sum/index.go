package most_frequent_subtree_sum

import "GiveMeAnOffer/defines"

// FindFrequentTreeSum https://leetcode-cn.com/problems/most-frequent-subtree-sum/
func FindFrequentTreeSum(root *defines.TreeNode) (res []int) {
	var memo map[int]int
	collectSum(root, memo)
	most := 0
	for key, val := range memo {
		if most == val {
			res = append(res, key)
		} else if most < val {
			most = val
			res = []int{}
		}
	}
	return
}

func collectSum(root *defines.TreeNode, memo map[int]int) int {
	if root == nil {
		return 0
	}
	sum := root.Val + collectSum(root.Left, memo) + collectSum(root.Right, memo)
	if v, ok := memo[sum]; ok {
		memo[sum] = v + 1
	} else {
		memo[sum] = 1
	}
	return sum
}
