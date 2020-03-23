package linkedlist

func GetIntersectionNode(headA, headB *ListNode) *ListNode {
	a := headA
	b := headB
	if a == nil || b == nil {
		return nil
	}
	for a != b {
		if a == nil {
			a = headB
		} else {
			a = a.Next
		}
		if b == nil {
			b = headA
		} else {
			b = b.Next
		}
	}
	return a
}
