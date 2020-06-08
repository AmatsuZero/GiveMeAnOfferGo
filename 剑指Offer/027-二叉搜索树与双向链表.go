package 剑指Offer

/*
题目：输入一棵二叉搜索树，将该二叉搜索树转换成一个排序的双向链表。要求不能创建任何新的结点，只能调整树中结点指针的指向
*/
func (node *BinaryTreeNode) Convert() (pHeadOfList *BinaryTreeNode) {
	var lastNodeIntList *BinaryTreeNode
	node.convert(&lastNodeIntList)
	// lastNodeIntList 指向双向链表的尾节点，我们需要返回头节点
	pHeadOfList = lastNodeIntList
	for pHeadOfList != nil && pHeadOfList.Left != nil {
		pHeadOfList = pHeadOfList.Left
	}
	return
}

func (node *BinaryTreeNode) convert(lastNodeIntList **BinaryTreeNode) {
	pCurrent := node
	if pCurrent.Left != nil {
		pCurrent.Left.convert(lastNodeIntList)
	}
	pCurrent.Left = *lastNodeIntList
	if *lastNodeIntList != nil {
		(*lastNodeIntList).Right = pCurrent
	}
	*lastNodeIntList = pCurrent
	if pCurrent.Right != nil {
		pCurrent.Right.convert(lastNodeIntList)
	}
}
