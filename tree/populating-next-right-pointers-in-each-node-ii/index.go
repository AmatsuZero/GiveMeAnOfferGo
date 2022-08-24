package populating_next_right_pointers_in_each_node_ii

type Node struct {
	Val   int
	Left  *Node
	Right *Node
	Next  *Node
}

func Connect(root *Node) *Node {
	if root == nil {
		return root
	}

	for cur := root; cur != nil; {
		dummy := &Node{} // 每层的假节点
		pre := dummy     // pre表示访下一层节点的前一个节点

		for ; cur != nil; cur = cur.Next { // 当前层
			if cur.Left != nil {
				pre.Next = cur.Left
				pre = pre.Next
			}

			if cur.Right != nil {
				pre.Next = cur.Right
				pre = pre.Next
			}
		}
		//把下一层串联成一个链表之后，让他赋值给cur，
		//后续继续循环，直到cur为空为止
		cur = dummy.Next
	}

	return root
}
