package 剑指Offer

/*
输入某二叉树的前序遍历和中序遍历的结果，请重建出该二叉树。

假设输入的前序遍历和中序遍历的结果中都不含重复的数字。例如输入前序遍历序列{1,2,4,7,3,5,6,8}和中序遍历序列{4,7,2,1,5,3,8,6}，则重建出图2.6所示的二叉树并输出它的头结点
*/
func Construct(preorder []int, inorder []int) (root *BinaryTreeNode) {
	if len(preorder) == 0 || len(inorder) == 0 {
		return
	}
	// 前序遍历序列的第一个数字是根节点的值
	rootValue := preorder[0]
	root = &BinaryTreeNode{Value: rootValue}

	startPreorder, endPreorder := 0, len(preorder)-1
	startInorder, endInorder := 0, len(inorder)-1

	// 由于值不重复，如果相等，直接返回
	if startPreorder == endInorder {
		if preorder[startPreorder] == inorder[startInorder] {
			return root
		} else {
			panic("Invalid input")
		}
	}

	// 在中序遍历中找到根结点的值
	rootInOrder := startInorder
	for rootInOrder <= endInorder && inorder[rootInOrder] != rootValue {
		rootInOrder++
	}

	if rootInOrder == endInorder && inorder[rootInOrder] != rootValue {
		panic("Invalid input")
	}

	leftLength := rootInOrder - startInorder
	leftPreorderEnd := startPreorder + leftLength

	if leftLength > 0 { // 构建左子树
		root.Left = Construct(preorder[1:leftPreorderEnd], inorder[startInorder:rootInOrder])
	}
	if leftLength < endPreorder-startPreorder { // 构建右子树
		root.Right = Construct(preorder[leftPreorderEnd+1:endPreorder], inorder[rootInOrder+1:endInorder+1])
	}

	return
}
