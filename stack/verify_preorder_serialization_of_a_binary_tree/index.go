package verifypreorderserializationofabinarytree

import "strings"

// https://leetcode-cn.com/problems/verify-preorder-serialization-of-a-binary-tree/
func IsValidSerialization(preorder string) bool {
	nodes, diff := strings.Split(preorder, ","), 1
	for _, node := range nodes {
		diff -= 1
		if diff < 0 {
			return false
		}
		if node != "#" {
			diff += 2
		}
	}
	return diff == 0
}
