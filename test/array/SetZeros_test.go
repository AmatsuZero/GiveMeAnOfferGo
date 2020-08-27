package array

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/array"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetZeros(t *testing.T) {
	input := [][]int{
		{1, 1, 1},
		{1, 0, 1},
		{1, 1, 1},
	}
	array.SetZeros(input)
	assert.Equal(t, [][]int{
		{1, 0, 1},
		{0, 0, 0},
		{1, 0, 1},
	}, input)

	input = [][]int{
		{0, 1, 2, 0},
		{3, 4, 5, 2},
		{1, 3, 1, 5},
	}
	array.SetZeros(input)
	assert.Equal(t, [][]int{
		{0, 0, 0, 0},
		{0, 4, 5, 0},
		{0, 3, 1, 0},
	}, input)
}
