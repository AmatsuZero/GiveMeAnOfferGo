package 剑指Offer

func (node *BinaryTreeNode) HasSubtree(root *BinaryTreeNode) bool {
	return HasSubtree(node, root)
}

/*
题目：输入两棵二叉树A和B，判断B是不是A的子结构。
*/
func HasSubtree(root1, root2 *BinaryTreeNode) (result bool) {
	if root1 != nil && root2 != nil {
		if root1.Value == root2.Value {
			result = doesTree1HaveTree2(root1, root2)
		}
		if !result {
			result = HasSubtree(root1.Left, root2)
		}
		if !result {
			result = HasSubtree(root1.Right, root2)
		}
	}
	return
}

func doesTree1HaveTree2(root1, root2 *BinaryTreeNode) bool {
	if root2 == nil {
		return true
	}
	if root1 == nil {
		return false
	}
	if root1.Value != root2.Value {
		return false
	}
	return doesTree1HaveTree2(root1.Left, root2.Left) &&
		doesTree1HaveTree2(root1.Right, root2.Right)
}
