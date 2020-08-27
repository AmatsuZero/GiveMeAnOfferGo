package array

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/array"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindDiagonalOrder(t *testing.T) {
	input := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	result := array.FindDiagonalOrder(input)
	assert.Equal(t, []int{1, 2, 4, 7, 5, 3, 6, 8, 9}, result)
}
