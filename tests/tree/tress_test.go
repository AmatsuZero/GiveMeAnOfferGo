package tree

import (
	"GiveMeAnOffer/defines"
	maximum_binary_tree "GiveMeAnOffer/tree/maximum-binary-tree"
	"GiveMeAnOffer/tree/same_tree"
	"testing"
)

func TestConstructMaximumBinaryTree(t *testing.T) {
	my := maximum_binary_tree.ConstructMaximumBinaryTree([]int{3, 2, 1, 6, 0, 5})
	ans := &defines.TreeNode{
		Val: 6,
		Left: &defines.TreeNode{
			Val: 3,
			Right: &defines.TreeNode{
				Val:   2,
				Right: &defines.TreeNode{Val: 1},
			},
		},
		Right: &defines.TreeNode{
			Val:  5,
			Left: &defines.TreeNode{Val: 0},
		},
	}

	if !same_tree.IsSameTree(my, ans) {
		t.Fail()
	}

	my = maximum_binary_tree.ConstructMaximumBinaryTree([]int{3, 2, 1})
	ans = &defines.TreeNode{
		Val: 3,
		Right: &defines.TreeNode{
			Val:   2,
			Right: &defines.TreeNode{Val: 1},
		},
	}
	if !same_tree.IsSameTree(my, ans) {
		t.Fail()
	}
}
