package 剑指Offer

/*
请完成一个函数，输入一个二叉树，该函数输出它的镜像
*/
func (node *BinaryTreeNode) MirrorRecursively() {
	if node.Left == nil && node.Right == nil {
		return
	}
	node.Left, node.Right = node.Right, node.Left
	if node.Left != nil {
		node.Left.MirrorRecursively()
	}
	if node.Right != nil {
		node.Right.MirrorRecursively()
	}
}

/*
镜像二叉树的循环实现
由于递归的本质是编译器生成了一个函数调用的栈，因此用循环来完成同样任务时，最简单的办法就是用一个辅助栈来模拟递归。首先把树的头结点放入栈中。
在循环中，只要栈不为空，弹出栈的栈顶结点，交换它的左右子树。如果它有左子树，把它的左子树压入栈中；如果它有右子树，把它的右子树压入栈中。
这样在下次循环中就能交换它儿子结点的左右子树了。
*/
func (node *BinaryTreeNode) MirrorIteratively() {
	stack := make([]*BinaryTreeNode, 0)
	stack = append(stack, node)
	for len(stack) > 0 {
		var pNode *BinaryTreeNode
		pNode, stack = stack[len(stack)-1], stack[:len(stack)-1]
		pNode.Left, pNode.Right = pNode.Right, pNode.Left
		if pNode.Left != nil {
			stack = append(stack, pNode.Left)
		}
		if pNode.Right != nil {
			stack = append(stack, pNode.Right)
		}
	}
}
