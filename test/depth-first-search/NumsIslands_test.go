package depth_first_search

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/LeetCode解题"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNumIslands(t *testing.T) {
	input := [][]byte{
		{'1', '1', '1', '1', '0'},
		{'1', '1', '0', '1', '0'},
		{'1', '1', '0', '0', '0'},
		{'0', '0', '0', '0', '0'},
	}
	assert.Equal(t, 1, LeetCode解题.NumIslands(input))

	input = [][]byte{
		{'1', '1', '0', '0', '0'},
		{'1', '1', '0', '0', '0'},
		{'0', '0', '1', '0', '0'},
		{'0', '0', '0', '1', '1'},
	}
	assert.Equal(t, 3, LeetCode解题.NumIslands(input))
}
