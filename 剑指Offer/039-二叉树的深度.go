package 剑指Offer

/*
输入一棵二叉树的根结点，求该树的深度。从根结点到叶结点依次经过的结点（含根、叶结点）形成树的一条路径，最长路径的长度为树的深度
*/
func (node *BinaryTreeNode) Depth() int {
	// 左子树深度
	leftDepth := 0
	if node.Left != nil {
		leftDepth = node.Left.Depth()
	}
	// 右子树深度
	rightDepth := 0
	if node.Right != nil {
		rightDepth = node.Right.Depth()
	}
	// 取左右子树最大深度的较大者，再加上自身
	if leftDepth > rightDepth {
		return leftDepth + 1
	} else {
		return rightDepth + 1
	}
}

/*
题目二：输入一棵二叉树的根结点，判断该树是不是平衡二叉树。如果某二叉树中任意结点的左右子树的深度相差不超过1，那么它就是一棵平衡二叉树
*/
func (node *BinaryTreeNode) IsBalanced() bool {
	depth := 0
	return node.isBalanced(&depth)
}

func (node *BinaryTreeNode) isBalanced(depth *int) bool {
	// 通过后序遍历的方式，求得左右子树深度
	isLeftBalanced := true
	leftDepth := 0
	if node.Left != nil {
		isLeftBalanced = node.Left.isBalanced(&leftDepth)
	}
	isRightBalanced := true
	rightDepth := 0
	if node.Right != nil {
		isRightBalanced = node.Right.isBalanced(&rightDepth)
	}
	if !isLeftBalanced || !isRightBalanced {
		return false
	}
	diff := leftDepth - rightDepth
	if diff <= 1 && diff >= -1 {
		if leftDepth > rightDepth {
			*depth += 1 + leftDepth
		} else {
			*depth += 1 + rightDepth
		}
		return true
	}
	return false
}
