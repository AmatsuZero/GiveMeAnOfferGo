package depth_first_search

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/Collections"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMyHashMap(t *testing.T) {
	hmap := Collections.NewHashMap()
	hmap.Put(12, 12)
	hmap.Put(9, 34)
	assert.Equal(t, 12, hmap.Get(12))
	assert.Equal(t, 34, hmap.Get(9))
	hmap.Remove(9)
	assert.Equal(t, -1, hmap.Get(9))
}
