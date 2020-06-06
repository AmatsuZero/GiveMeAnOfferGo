package 剑指Offer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindKthToTail(t *testing.T) {
	list := RandomIntList(10)
	assert.Equal(t, list.NodeAt(8), list.FindKthToTail(2))
	assert.Equal(t, list.NodeAt(0), list.FindKthToTail(10))
	assert.Equal(t, list.NodeAt(9), list.FindKthToTail(1))
	assert.Equal(t, list.NodeAt(-1), list.NodeAt(9))
	assert.Nil(t, list.FindKthToTail(0))
	assert.Nil(t, list.FindKthToTail(11))
}

func TestFindMiddleNode(t *testing.T) {
	list := RandomIntList(5) // 奇数长度
	assert.Equal(t, list.NodeAt(2), list.MidNode())

	list = RandomIntList(4) // 偶数长度
	assert.Equal(t, list.NodeAt(2), list.MidNode())

	list = RandomIntList(1)
	assert.Equal(t, list.NodeAt(0), list.MidNode())
}

func TestIsCircleList(t *testing.T) {
	list := RandomIntList(1)
	insertion := RandomIntList(3)
	list.Next = insertion
	insertion.Next = list
	assert.True(t, list.IsCircleList())

	list = RandomIntList(4)
	assert.False(t, list.IsCircleList())
}
