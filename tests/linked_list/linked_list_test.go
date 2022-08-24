package linked_list

import (
	"GiveMeAnOffer/linked_list/copy_random_list"
	"testing"
)

func TestLinkedList(t *testing.T) {
	// [[7,null],[13,0],[11,4],[10,2],[1,0]]
	arr := []*copy_random_list.Node{
		&copy_random_list.Node{Val: 7},
		&copy_random_list.Node{Val: 13},
		&copy_random_list.Node{Val: 11},
		&copy_random_list.Node{Val: 10},
		&copy_random_list.Node{Val: 1},
	}
	var head *copy_random_list.Node
	cur := head
	for i, node := range arr {
		if head == nil {
			head = node
			cur = head
		} else {
			cur.Next = node
			cur = node
		}
		switch i {
		case 1:
			node.Random = arr[0]
		case 2:
			node.Random = arr[4]
		case 3:
			node.Random = arr[2]
		case 4:
			node.Random = arr[0]
		}
	}
	ans := copy_random_list.CopyRandomList(head)
	if ans != nil {
		t.Fail()
	}
}
