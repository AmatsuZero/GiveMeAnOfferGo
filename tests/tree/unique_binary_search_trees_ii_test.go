package tree

import (
	"GiveMeAnOffer/defines"
	"GiveMeAnOffer/tree/unique_binary_search_trees_ii"
	"reflect"
	"strconv"
	"testing"
)

func TestUniqueBinarySearchTreesII(t *testing.T) {
	input := 3
	trees := unique_binary_search_trees_ii.GenerateTrees(input)
	var output [][]string
	for _, tree := range trees {
		output = append(output, helper(tree))
	}
	if reflect.DeepEqual(output, [][]string{
		{"1", "null", "2", "null", "3"},
		{"1", "null", "3", "2"},
		{"2", "1", "3"},
		{"3", "1", "null", "null", "2"},
		{"3", "2", "null", "1"},
	}) {
		t.Logf("输入 %d 正确", input)
	} else {
		t.Logf("输入 %d 错误", input)
	}
}

func helper(n *defines.TreeNode) (output []string) {
	if n == nil {
		output = append(output, "null")
		return
	}
	output = append(output, strconv.Itoa(n.Val))
	if n.Left != nil || n.Right != nil { // 只有左子树或者右子树存在，再打印
		output = append(output, helper(n.Left)...)
		output = append(output, helper(n.Right)...)
	}
	return
}
