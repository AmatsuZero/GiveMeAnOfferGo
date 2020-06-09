package 剑指Offer

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/剑指Offer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPrintListReverselyIteratively(t *testing.T) {
	input := 剑指Offer.RandomIntList(10)
	lhs := input.IntArray()
	for left, right := 0, len(lhs)-1; left < right; left, right = left+1, right-1 {
		lhs[left], lhs[right] = lhs[right], lhs[left]
	}
	rhs := make([]int, 0)
	input.TraverseReversely(func(val int) {
		rhs = append(rhs, val)
	})
	assert.Equal(t, lhs, rhs)

	lhs, rhs = lhs[:0], rhs[:0]
	剑指Offer.TraverseReversely(nil, func(val int) {
		rhs = append(rhs, val)
	})
	assert.Equal(t, lhs, rhs)
}
