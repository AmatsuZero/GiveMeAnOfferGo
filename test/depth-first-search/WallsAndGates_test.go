package depth_first_search

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/LeetCode解题"
	"github.com/stretchr/testify/assert"
	"testing"
)

const INF = 2147483647

func TestWallsAndGates(t *testing.T) {
	input := [][]int{
		{INF, -1, 0, INF},
		{INF, INF, INF, -1},
		{INF, -1, INF, -1},
		{0, -1, INF, INF},
	}
	LeetCode解题.WallsAndGates(input)
	assert.Equal(t, [][]int{
		{3, -1, 0, 1},
		{2, 2, 1, -1},
		{1, -1, 2, -1},
		{0, -1, 3, 4},
	}, input)
}
