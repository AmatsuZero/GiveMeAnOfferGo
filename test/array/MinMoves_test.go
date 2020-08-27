package array

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/array"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMinMoves(t *testing.T) {
	assert.Equal(t, 3, array.MinMoves([]int{1, 2, 3}))
}
