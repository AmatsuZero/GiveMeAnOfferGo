package 剑指Offer

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/剑指Offer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReverseList(t *testing.T) {
	list := 剑指Offer.RandomIntList(10)
	lhs := list.IntArray()
	for left, right := 0, len(lhs)-1; left < right; left, right = left+1, right-1 {
		lhs[left], lhs[right] = lhs[right], lhs[left]
	}
	list = list.Reverse()
	rhs := list.IntArray()
	assert.Equal(t, lhs, rhs)

	list = 剑指Offer.RandomIntList(1)
	assert.NotPanics(t, func() {
		list = list.Reverse()
		t.Log(list)
	})
}
