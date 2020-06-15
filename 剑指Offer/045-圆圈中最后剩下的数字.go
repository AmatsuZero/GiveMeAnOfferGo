package 剑指Offer

/*
0,1,…,n－1这n个数字排成一个圆圈，从数字0开始每次从这个圆圈里删除第m个数字。求出这个圆圈里剩下的最后一个数字
*/
func LastRemain(n, m uint) int {
	if n < 1 || m < 1 {
		return -1
	}
	head := &ListNode{Val: 0}
	for i := 0; i < int(n); i++ {
		head.AddToTail(i)
	}
	curHead := head
	for curHead != nil {
		for i := 0; i < int(m); i++ {
			curHead = curHead.Next
			if curHead == nil {
				curHead = head
			}
		}
		toBDeleted := curHead
		curHead = toBDeleted.Next
		head.DeleteNode(toBDeleted)
	}
	return head.Val
}
