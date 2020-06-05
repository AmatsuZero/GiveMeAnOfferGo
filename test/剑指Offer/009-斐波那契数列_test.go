package 剑指Offer

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/剑指Offer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFibonacci(t *testing.T) {
	assert.Equal(t, int64(3), 剑指Offer.Fibonacci(4))
}

func TestStairs(t *testing.T) {
	assert.Equal(t, int64(2), 剑指Offer.FrogSteps(2))
}
