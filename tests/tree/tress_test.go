package tree

import (
	"GiveMeAnOffer/defines"
	maximum_binary_tree "GiveMeAnOffer/tree/maximum-binary-tree"
	"GiveMeAnOffer/tree/print_binary_tree"
	"GiveMeAnOffer/tree/same_tree"
	"reflect"
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

func TestPrintOfTree(t *testing.T) {
	tree := &defines.TreeNode{
		Val:  1,
		Left: &defines.TreeNode{Val: 2},
	}
	ans := print_binary_tree.PrintTree(tree)
	if !reflect.DeepEqual(ans, [][]string{
		{"", "1", ""},
		{"2", "", ""},
	}) {
		t.Fail()
	}

	tree = &defines.TreeNode{
		Val: 1,
		Left: &defines.TreeNode{
			Val:   2,
			Right: &defines.TreeNode{Val: 4},
		},
		Right: &defines.TreeNode{Val: 3},
	}

	ans = print_binary_tree.PrintTree(tree)
	if !reflect.DeepEqual(ans, [][]string{
		{"", "", "", "1", "", "", ""},
		{"", "2", "", "", "", "3", ""},
		{"", "", "4", "", "", "", ""},
	}) {
		t.Fail()
	}
}
