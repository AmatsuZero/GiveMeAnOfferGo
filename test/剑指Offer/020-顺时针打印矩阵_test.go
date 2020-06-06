package 剑指Offer

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/剑指Offer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPrintMatrixClockWisely(t *testing.T) {
	input := [][]int{
		{1, 2, 3, 4},
		{5, 6, 7, 8},
		{9, 10, 11, 12},
		{13, 14, 15, 16},
	}
	lhs := make([]int, 0)
	剑指Offer.EnumerateMatrixClockWisely(input, func(val int) {
		lhs = append(lhs, val)
	})
	assert.Equal(t, lhs, []int{1, 2, 3, 4, 8, 12, 16, 15, 14, 13, 9, 5, 6, 7, 11, 10})

	input = [][]int{{1, 2, 3, 4}}
	lhs = lhs[:0]
	剑指Offer.EnumerateMatrixClockWisely(input, func(val int) {
		lhs = append(lhs, val)
	})
	assert.Equal(t, lhs, []int{1, 2, 3, 4})

	input = [][]int{
		{1}, {5}, {9}, {13},
	}
	lhs = lhs[:0]
	剑指Offer.EnumerateMatrixClockWisely(input, func(val int) {
		lhs = append(lhs, val)
	})
	assert.Equal(t, lhs, []int{1, 5, 9, 13})
}
