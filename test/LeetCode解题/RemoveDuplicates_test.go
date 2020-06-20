package LeetCode解题

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/LeetCode解题"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRemoveElement(t *testing.T) {
	lhs := []int{1}
	rhs := LeetCode解题.RemoveElement(lhs, 1)
	assert.Equal(t, lhs, rhs)

	lhs = []int{1, 1, 2, 2}
	rhs = LeetCode解题.RemoveElement(lhs, 3)
	assert.Equal(t, lhs, rhs)

	lhs = []int{1, 1, 2, 2}
	rhs = LeetCode解题.RemoveElement(lhs, 2)
	assert.Equal(t, []int{1, 1}, rhs)

	lhs = []int{}
	rhs = LeetCode解题.RemoveElement(lhs, 3)
	assert.Equal(t, lhs, rhs)
}

func TestRemoveDuplicates(t *testing.T) {
	lhs := []int{1}
	rhs := LeetCode解题.RemoveDuplicates(lhs)
	assert.Equal(t, lhs, rhs)

	lhs = []int{1, 1, 2, 2, 3}
	rhs = LeetCode解题.RemoveDuplicates(lhs)
	assert.Equal(t, []int{1, 2, 3}, rhs)
}

func TestRemoveDuplicatesWithRepeat(t *testing.T) {
	lhs := []int{1}
	rhs := LeetCode解题.RemoveSomeDuplicates(lhs, 0)
	assert.Equal(t, lhs, rhs)

	lhs = []int{1, 1, 2, 2, 3}
	rhs = LeetCode解题.RemoveSomeDuplicates(lhs, 2)
	assert.Equal(t, []int{1, 1, 2, 2, 3}, rhs)
}
