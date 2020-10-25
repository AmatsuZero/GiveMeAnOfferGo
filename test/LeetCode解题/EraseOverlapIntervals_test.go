package LeetCode解题

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/LeetCode解题"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEraseOverlapIntervals(t *testing.T) {
	assert.Equal(t, 1, LeetCode解题.EraseOverlapIntervals([][]int{
		{1, 2}, {2, 3}, {3, 4}, {1, 3},
	}))

	assert.Equal(t, 2, LeetCode解题.EraseOverlapIntervals([][]int{
		{1, 2}, {1, 2}, {1, 2},
	}))

	assert.Equal(t, 0, LeetCode解题.EraseOverlapIntervals([][]int{
		{1, 2}, {2, 3},
	}))
}
