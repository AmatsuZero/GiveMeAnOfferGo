package datastructure

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/Collections"
	"testing"
)

func TestBloomFilter(t *testing.T) {
	bf := Collections.NewBloomFilter(1024)
	bf.Add([]byte("hello"))
	bf.Add([]byte("world"))
	bf.Add([]byte("sir"))
	bf.Add([]byte("madam"))
	bf.Add([]byte("io"))

	if bf.Test([]byte("hello")) {

	}
	if bf.Test([]byte("world")) {

	}
	if !bf.Test([]byte("hi")) {

	}
}
