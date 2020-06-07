package 剑指Offer

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/剑指Offer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindPath(t *testing.T) {
	tree := 剑指Offer.Construct([]int{10, 5, 4, 7, 12}, []int{4, 5, 7, 10, 12})
	tree.FindPath(22, func(path []int) {
		sum := 0
		for _, elm := range path {
			sum += elm
		}
		assert.Equal(t, 22, sum)
	})
}
