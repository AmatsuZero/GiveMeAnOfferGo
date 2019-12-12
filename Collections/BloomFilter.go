package Collections

import (
	"github.com/spaolacci/murmur3"
	"hash"
	"hash/fnv"
)

type BloomFilter struct {
	bitSet   []bool
	k        uint
	n        uint
	m        uint
	hasFuncs []hash.Hash64
}

func NewBloomFilter(size uint) *BloomFilter {
	return &BloomFilter{
		bitSet:   make([]bool, size),
		k:        3,
		m:        size,
		n:        0,
		hasFuncs: []hash.Hash64{murmur3.New64(), fnv.New64(), fnv.New64a()},
	}
}

func (bf *BloomFilter) Add(item []byte) {
	hashes := bf.hashValues(item)
	i := uint(0)
	for {
		if i >= bf.k {
			break
		}

		position := uint(hashes[i]) % bf.m
		bf.bitSet[position] = true

		i += 1
	}

	bf.n += 1
}

func (bf *BloomFilter) hashValues(item []byte) []uint64 {
	var result []uint64
	for _, hashFunc := range bf.hasFuncs {
		if _, err := hashFunc.Write(item); err != nil {
			return result
		}
		result = append(result, hashFunc.Sum64())
		hashFunc.Reset()
	}

	return result
}

func (bf *BloomFilter) Test(item []byte) (exists bool) {
	hashes := bf.hashValues(item)
	i := uint(0)
	exists = true

	for {
		if i >= bf.k {
			break
		}

		position := uint(hashes[i]) % bf.m
		if !bf.bitSet[position] {
			exists = false
			break
		}
		i += 1
	}
	return
}
