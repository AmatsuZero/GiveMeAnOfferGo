package 剑指Offer

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/剑指Offer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPrintFromTopToBottom(t *testing.T) {
	tree := 剑指Offer.Construct([]int{8, 6, 5, 7, 10, 9, 11}, []int{5, 6, 7, 8, 9, 10, 11})
	rhs := make([]int, 0)
	tree.EnumerateFromTopToBottom(func(value int) {
		rhs = append(rhs, value)
	})
	assert.Equal(t, []int{8, 6, 10, 5, 7, 9, 11}, rhs)
}
