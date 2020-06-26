package LeetCode解题

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/LeetCode解题"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTwoSum(t *testing.T) {
	index1, index2 := LeetCode解题.TwoSum([]int{2, 7, 11, 5}, 9)
	assert.Equal(t, 1, index1)
	assert.Equal(t, 2, index2)

	index1, index2 = LeetCode解题.TwoSum([]int{1, 3, 3, 5}, 6)
	assert.Equal(t, 2, index1)
	assert.Equal(t, 3, index2)

	index1, index2 = LeetCode解题.TwoSum([]int{1, 7, 0, 5}, 7)
	assert.Equal(t, 2, index1)
	assert.Equal(t, 3, index2)
}

func TestTreeSum(t *testing.T) {
	input := []int{9, -3, -5, 1, 2}
	lhs := [][]int{{-3, 1, 2}}
	rhs := LeetCode解题.ThreeSum(input)
	assert.Equal(t, lhs, rhs)
}
