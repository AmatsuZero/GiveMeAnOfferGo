package 剑指Offer

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/剑指Offer"
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func TestConvertBSTToBiDirectionList(t *testing.T) {
	tree := 剑指Offer.Construct([]int{10, 6, 4, 8, 14, 12, 16}, []int{4, 6, 8, 10, 12, 14, 16})
	list := tree.Convert()
	arr := make([]int, 0)
	for node := list.Right; node != nil; node = node.Right {
		arr = append(arr, node.Value)
	}
	assert.True(t, sort.IntsAreSorted(arr))
}
