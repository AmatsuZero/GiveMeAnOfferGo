package 剑指Offer

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/剑指Offer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRollDiceProbabilities(t *testing.T) {
	t.Log(剑指Offer.ProbabilityOfDiceNumbers(2))
	assert.Nil(t, 剑指Offer.ProbabilityOfDiceNumbers(0))
}
