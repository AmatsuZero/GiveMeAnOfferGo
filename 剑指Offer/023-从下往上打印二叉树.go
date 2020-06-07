package 剑指Offer

/*
题目：从上往下打印出二叉树的每个结点，同一层的结点按照从左到右的顺序打印
*/
func (node *BinaryTreeNode) EnumerateFromTopToBottom(block func(value int)) {
	queue := make([]*BinaryTreeNode, 0)
	queue = append(queue, node)
	for len(queue) > 0 {
		var pNode *BinaryTreeNode
		pNode, queue = queue[0], queue[1:]
		block(pNode.Value)
		if pNode.Left != nil {
			queue = append(queue, pNode.Left)
		}
		if pNode.Right != nil {
			queue = append(queue, pNode.Right)
		}
	}
}
