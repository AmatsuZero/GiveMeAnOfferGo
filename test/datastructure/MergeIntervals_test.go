package datastructure

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/array"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMergeIntervals(t *testing.T) {
	assert.Equal(t, [][]int{{1, 6}, {8, 10}, {15, 18}}, array.Merge([][]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}}))
}
