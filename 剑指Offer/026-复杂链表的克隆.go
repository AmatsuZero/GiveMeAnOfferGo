package 剑指Offer

/*
题目：请实现函数ComplexListNode* Clone（ComplexListNode* pHead），复制一个复杂链表。
在复杂链表中，每个结点除了有一个m_pNext指针指向下一个结点外，还有一个m_pSibling 指向链
*/

func (node *ComplexListNode) Clone() *ComplexListNode {
	node.cloneNodes()
	node.connectSiblingNodes()
	return node.reconnectNodes()
}

// 复制原始链表的任意结点N并创建新结点N'，再把N'链接到N的后面
func (node *ComplexListNode) cloneNodes() {
	pNode := node
	for pNode.Next != nil {
		pCloned := &ComplexListNode{
			Val:  pNode.Val,
			Next: pNode.Next,
		}
		pNode.Next = pCloned
		pNode = pCloned.Next
	}
}

// 如果原始链表上的结点N的m_pSibling指向S，则它对应的复制结点N'的m_pSibling指向S的下一结点S'
func (node *ComplexListNode) connectSiblingNodes() {
	pNode := node
	for node != nil {
		pCloned := pNode.Next
		if pNode.Sibling != nil {
			pCloned.Sibling = pNode.Sibling.Next
		}
		pNode = pCloned.Next
	}
}

// 把第二步得到的链表拆分成两个链表，奇数位置上的结点组成原始链表，偶数位置上的结点组成复制出来的链表。
func (node *ComplexListNode) reconnectNodes() *ComplexListNode {
	var pClonedHead, pClonedNode *ComplexListNode
	pNode := node
	if pNode != nil {
		pClonedNode = pNode.Next
		pClonedHead = pClonedNode

		pNode.Next = pClonedNode.Next
		pNode = pNode.Next
	}
	for pNode != nil {
		if pClonedNode != nil {
			pClonedNode.Next = pNode.Next
			pClonedNode = pClonedNode.Next
			if pClonedNode.Next != nil {
				pNode.Next = pClonedNode.Next
			}
		}
		pNode = pNode.Next
	}
	return pClonedHead
}
