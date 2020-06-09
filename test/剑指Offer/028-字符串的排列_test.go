package 剑指Offer

import (
	"fmt"
	"github.com/AmatsuZero/GiveMeAnOfferGo/剑指Offer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStringPermutation(t *testing.T) {
	剑指Offer.Permutation("abc")
}

func TestCubeVertex(t *testing.T) {
	input := [8]int{1, 2, 3, 1, 2, 3, 2, 2}
	assert.True(t, 剑指Offer.CubVertex(input))
	input = [8]int{1, 2, 3, 1, 8, 3, 2, 2}
	assert.False(t, 剑指Offer.CubVertex(input))
}

func TestQueenQuestion(t *testing.T) {
	b := 剑指Offer.Queen()
	fmt.Println(b)
}
