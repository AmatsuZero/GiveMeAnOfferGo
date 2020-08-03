package LeetCode解题

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/LeetCode解题"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMergeSortedArray(t *testing.T) {
	lhs := append([]int{1, 2, 4}, make([]int, 3)...)
	rhs := []int{2, 5, 6}
	LeetCode解题.MergeSortedArray(lhs, rhs, 3, 3)
	assert.Equal(t, []int{1, 2, 2, 4, 5, 6}, lhs)
}
