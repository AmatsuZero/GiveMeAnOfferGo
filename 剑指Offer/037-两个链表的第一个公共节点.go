package 剑指Offer

func (node *ListNode) Len() (length int) {
	for head := node; head != nil; head = head.Next {
		length++
	}
	return length
}

/*
题目：输入两个链表，找出它们的第一个公共结点
*/
func FindFirstCommonNode(head1, head2 *ListNode) *ListNode {
	nLength1, nLength2 := head1.Len(), head2.Len()
	nLengthDif := nLength1 - nLength2
	headLong, headShort := head1, head2
	if nLength2 > nLength1 {
		headLong, headShort = head2, head1
		nLengthDif = nLength2 - nLength1
	}
	// 现在长链表上走几步，再同时在两个链表上遍历
	for i := 0; i < nLengthDif; i++ {
		headLong = headLong.Next
	}
	for headLong != nil &&
		headShort != nil &&
		headLong != headShort {
		headShort = headShort.Next
		headLong = headLong.Next
	}

	return headLong
}
