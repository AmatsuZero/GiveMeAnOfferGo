package Sort

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/Sort"
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func TestDutchFlag(t *testing.T) {
	input := []int{10, 0, 3, 9, 2, 14, 8, 27, 1, 5, 8, -1, 26}
	Sort.QuickSortDutchFlag(input, 0, len(input)-1)
	assert.True(t, sort.IntsAreSorted(input))
}
