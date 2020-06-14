package 剑指Offer

/*
题目：输入一棵二叉树和一个整数，打印出二叉树中结点值的和为输入整数的所有路径。
从树的根结点开始往下一直到叶结点所经过的结点形成一条路径。
*/

func (node *BinaryTreeNode) FindPath(expectedNum int, block func(path []int)) {
	if node != nil {
		return
	}
	currentSum := 0
	path := make([]int, 0)
	node.findPath(expectedNum, &path, &currentSum, block)
}

func (node *BinaryTreeNode) findPath(expectedNum int, path *[]int, currentSum *int, block func(path []int)) {
	if node == nil {
		return
	}
	*currentSum += node.Value
	*path = append(*path, node.Value)
	// 如果是叶子节点，并且路径上节点的和等于输出的值
	if *currentSum == expectedNum && node.IsLeaf() {
		block(*path)
	}
	// 如果不是叶子节点，遍历子节点
	node.Left.findPath(expectedNum, path, currentSum, block)
	node.Right.findPath(expectedNum, path, currentSum, block)
	// 返回到父节点之前，在路径上删除当前节点
	// 并在 currentSum 中减去当前节点的值
	*currentSum -= node.Value
	*path = (*path)[:len(*path)-1]
}

func (node *BinaryTreeNode) IsLeaf() bool {
	return node.Left == nil && node.Right == nil
}
