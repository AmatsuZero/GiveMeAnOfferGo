package 剑指Offer

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/剑指Offer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCQueue(t *testing.T) {
	queue := 剑指Offer.NewCQueue()
	queue.AppendTail(10)
	queue.AppendTail(12)
	queue.AppendTail(13)

	assert.Equal(t, queue.DeleteHead(), 10)

	queue.AppendTail(14)
	assert.Equal(t, queue.DeleteHead(), 12)
}
