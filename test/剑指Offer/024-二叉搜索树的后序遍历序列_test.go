package 剑指Offer

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/剑指Offer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVerifyPostSequenceOfBST(t *testing.T) {
	assert.True(t, 剑指Offer.VerifyPostSequenceOfBST([]int{5, 7, 6, 9, 11, 10, 8}))
	assert.False(t, 剑指Offer.VerifyPostSequenceOfBST([]int{7, 4, 6, 5}))
	assert.False(t, 剑指Offer.VerifyPostSequenceOfBST([]int{8, 6, 5, 7, 10, 9, 11}))
}

func TestTestVerifyPreSequenceOfBST(t *testing.T) {
	assert.True(t, 剑指Offer.VerifyPreSequenceOfBST([]int{8, 6, 5, 7, 10, 9, 11}))
	assert.False(t, 剑指Offer.VerifyPreSequenceOfBST([]int{5, 7, 6, 9, 11, 10, 8}))
}
