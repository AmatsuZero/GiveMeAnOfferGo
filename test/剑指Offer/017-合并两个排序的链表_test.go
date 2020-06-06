package 剑指Offer

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/剑指Offer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMergeSortedList(t *testing.T) {
	lhs := 剑指Offer.NewListWithArray([]int{1, 3, 5, 7, 9})
	rhs := 剑指Offer.NewListWithArray([]int{2, 4, 6, 8})
	list := lhs.Merge(rhs)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, list.IntArray())
}
