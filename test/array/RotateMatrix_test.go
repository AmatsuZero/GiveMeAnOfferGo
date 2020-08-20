package array

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/array"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRotateMatrix(t *testing.T) {
	source := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	array.RotateMatrix(source)
	assert.Equal(t, [][]int{
		{7, 4, 1},
		{8, 5, 2},
		{9, 6, 3},
	}, source)
}
