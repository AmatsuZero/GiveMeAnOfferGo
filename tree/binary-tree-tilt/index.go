package binary_tree_tilt

import (
	"GiveMeAnOffer/defines"
	"math"
)

// FindTilt https://leetcode-cn.com/problems/binary-tree-tilt/
func FindTilt(root *defines.TreeNode) int {
	if root == nil {
		return 0
	}
	_, sum := findTiltDFS(root, 0)
	return sum
}

func findTiltDFS(root *defines.TreeNode, lastSum int) (valueSum, tiltSum int) {
	if root == nil {
		return 0, lastSum
	}
	left, tiltSum := findTiltDFS(root.Left, lastSum)
	right, tiltSum := findTiltDFS(root.Right, tiltSum)
	tiltSum += int(math.Abs(float64(left) - float64(right)))
	valueSum = root.Val + left + right
	return
}
