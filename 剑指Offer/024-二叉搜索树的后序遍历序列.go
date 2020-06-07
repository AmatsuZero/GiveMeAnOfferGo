package 剑指Offer

/*
题目：输入一个整数数组，判断该数组是不是某二叉搜索树的后序遍历的结果。
如果是则返回true，否则返回false。假设输入的数组的任意两个数字都互不相同。
*/
func VerifyPostSequenceOfBST(sequence []int) bool {
	if len(sequence) == 0 {
		return false
	}
	root := sequence[len(sequence)-1]
	// 在二叉搜索树中左子树的结点小于根结点
	var i int
	for i = 0; i < len(sequence)-1; i++ {
		if sequence[i] > root {
			break
		}
	}
	// 在二叉搜索树中右子树的结点大于根结点
	for j := i; j < len(sequence)-1; j++ {
		if sequence[j] < root {
			return false
		}
	}
	// 判读左子树是不是二叉搜索树
	left := true
	if i > 0 {
		left = VerifyPostSequenceOfBST(sequence[:i])
	}
	// 判断右子树是不是二叉搜索树
	right := true
	if i < len(sequence)-1 {
		right = VerifyPostSequenceOfBST(sequence[i : len(sequence)-1])
	}
	return left && right
}

/*
输入一个整数数组，判断该数组是不是某二叉搜索树的前序遍历的结果。
这和前面问题的后序遍历很类似，只是在前序遍历得到的序列中，第一个数字是根结点的值
*/
func VerifyPreSequenceOfBST(sequence []int) bool {
	if len(sequence) == 0 {
		return false
	}
	root := sequence[0]
	// 在二叉搜索树中左子树的结点小于根结点
	var i int
	for i = 1; i < len(sequence); i++ {
		if sequence[i] > root {
			break
		}
	}
	// 在二叉搜索树中右子树的结点大于根结点
	for j := i; j < len(sequence); j++ {
		if sequence[j] < root {
			return false
		}
	}
	// 判读左子树是不是二叉搜索树
	left := true
	if i > 1 {
		left = VerifyPreSequenceOfBST(sequence[1:i])
	}
	// 判断右子树是不是二叉搜索树
	right := true
	if i < len(sequence) {
		right = VerifyPreSequenceOfBST(sequence[i:])
	}
	return left && right
}
