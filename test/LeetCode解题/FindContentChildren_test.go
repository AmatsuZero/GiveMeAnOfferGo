package LeetCode解题

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/LeetCode解题"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindContentChildren(t *testing.T) {
	assert.Equal(t, 1, LeetCode解题.FindContentChildren([]int{1, 2, 3}, []int{1, 1}))
	assert.Equal(t, 2, LeetCode解题.FindContentChildren([]int{1, 2}, []int{1, 2, 3}))
}
