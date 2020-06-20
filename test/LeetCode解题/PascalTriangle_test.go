package LeetCode解题

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/LeetCode解题"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGeneratePascalTriangle(t *testing.T) {
	assert.Equal(t, LeetCode解题.PascalTriangle{{1}}, LeetCode解题.NewPascalTriangle(1))
	assert.Equal(t, LeetCode解题.PascalTriangle{}, LeetCode解题.NewPascalTriangle(0))
	assert.Equal(t, LeetCode解题.PascalTriangle{
		{1},
		{1, 1},
		{1, 2, 1},
		{1, 3, 3, 1},
		{1, 4, 6, 4, 1},
	}, LeetCode解题.NewPascalTriangle(5))
}

func TestGetRows(t *testing.T) {
	triangle := LeetCode解题.NewPascalTriangle(0)
	assert.Equal(t, []int{1}, triangle.GetRow(0))
	assert.Equal(t, []int{1, 1}, triangle.GetRow(1))
	assert.Equal(t, []int{1, 3, 3, 1}, triangle.GetRow(3))
	assert.Equal(t, []int{1, 4, 6, 4, 1}, triangle.GetRow(4))
}
