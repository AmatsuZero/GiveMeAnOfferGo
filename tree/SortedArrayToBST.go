package tree

/*
将一个按照升序排列的有序数组，转换为一棵高度平衡二叉搜索树。

本题中，一个高度平衡二叉树是指一个二叉树每个节点 的左右两个子树的高度差的绝对值不超过 1。

示例:

给定有序数组: [-10,-3,0,5,9],

一个可能的答案是：[0,-3,9,-10,null,5]，它可以表示下面这个高度平衡二叉搜索树：

      0
     / \
   -3   9
   /   /
 -10  5
*/

func SortedArrayToBST(nums []int) (root *TreeNode) {
	return sortedArrayToBST(nums, 0, len(nums)-1)
}

func sortedArrayToBST(arr []int, start int, end int) (node *TreeNode) {
	if start > end {
		return
	}

	mid := start + (end-start)/2
	node = &TreeNode{Val: arr[mid]}
	node.Left = sortedArrayToBST(arr, start, mid-1)
	node.Right = sortedArrayToBST(arr, mid+1, end)
	return
}
