package 剑指Offer

import (
	. "github.com/AmatsuZero/GiveMeAnOfferGo/剑指Offer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeleteNode(t *testing.T) {
	list := RandomIntList(10)
	lhs := list.IntArray()
	toBeDeleted := list.NodeAt(4)
	list.DeleteNode(toBeDeleted)
	rhs := list.IntArray()
	assert.Equal(t, len(lhs)-len(rhs), 1)

	// 删除尾节点
	lhs = rhs
	toBeDeleted = list.NodeAt(len(rhs) - 1)
	list.DeleteNode(toBeDeleted)
	rhs = list.IntArray()
	assert.Equal(t, len(lhs)-len(rhs), 1)

	// 只有一个节点
	list = &ListNode{Val: 10}
	toBeDeleted = list
	DeleteNode(list, toBeDeleted)
}
