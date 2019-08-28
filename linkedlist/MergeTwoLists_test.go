package linkedlist

import (
	"testing"
)

func TestMergeTwoLists(t *testing.T) {
	l1 := BuildLinkedList([]int{1, 2, 4})
	l2 := BuildLinkedList([]int{1, 3, 4})
	result := MergeTwoLists(l1, l2)
	target := BuildLinkedList([]int{1, 1, 2, 3, 4, 4})
	if result.IsEqual(target) {
		t.Log("通过: " + result.String())
	} else {
		t.Error("不通过: \n" + "Target: " + target.String() + "\nResult: " + result.String())
	}
	t.Log(target.ToArray())
}
