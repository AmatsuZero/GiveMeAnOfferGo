package 剑指Offer

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/剑指Offer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStackWithMin(t *testing.T) {
	stack := &剑指Offer.StackWithMin{}
	stack.Push(2)
	// 押入大数
	stack.Push(5)
	assert.Equal(t, 2, stack.Min())
	// 押入小数
	stack.Push(1)
	assert.Equal(t, 1, stack.Min())
	// 弹出栈的数字不是最小的元素
	stack.Push(4)
	assert.NotEqual(t, stack.Pop(), stack.Min())
	// 弹出的数字是最小的元素
	assert.Equal(t, stack.Min(), stack.Pop())
}
