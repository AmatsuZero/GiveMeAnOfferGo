package datastructure

import (
	"GiveMeAnOfferGo/datastructure"
	"testing"
)

func TestBloomFilter(t *testing.T) {
	bf := datastructure.NewBloomFilter(1024)
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
