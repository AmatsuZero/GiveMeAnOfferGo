package 剑指Offer

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/剑指Offer"
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func TestCQueue(t *testing.T) {
	queue := &剑指Offer.CQueue{}
	queue.AppendTail(10)
	queue.AppendTail(12)
	queue.AppendTail(13)

	assert.Equal(t, queue.DeleteHead(), 10)

	queue.AppendTail(14)
	assert.Equal(t, queue.DeleteHead(), 12)
}

func TestCStack(t *testing.T) {
	stack := &剑指Offer.CStack{}
	stack.Push(10)
	stack.Push(11)
	stack.Push(13)
	assert.Equal(t, stack.Pop(), 13)

	stack.Push(14)
	assert.Equal(t, stack.Pop(), 14)
}

func TestQuickSort(t *testing.T) {
	source := 剑指Offer.RandomIntArray(10)
	剑指Offer.QuickSort(&source)
	assert.True(t, sort.IntsAreSorted(source))
}
