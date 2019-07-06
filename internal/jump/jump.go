package jump

import (
	"hash"
	"hash/crc64"
	"io"
	"sync"
)

var hasher = sync.Pool{
	New: func() interface{} {
		return crc64.New(crc64.MakeTable(crc64.ECMA))
	},
}

// Hash consistently chooses a hash bucket number in the range [0, numBuckets) for the given key.
// numBuckets must be >= 1.
func Hash(key uint64, numBuckets int) int32 {
	var b int64 = -1
	var j int64

	for j < int64(numBuckets) {
		b = j
		key = key*2862933555777941757 + 1
		j = int64(float64(b+1) * (float64(int64(1)<<31) / float64((key>>33)+1)))
	}

	return int32(b)
}

// HashString works like hash but accept string as an input.
func HashString(key string, numBuckets int) int32 {
	// jump.Hash returns values from 0.
	k := Hash(Sum64(key), numBuckets)

	return k
}

// Sum64 ...
func Sum64(key string) uint64 {
	h := hasher.Get().(hash.Hash64)
	if _, err := io.WriteString(h, key); err != nil {
		panic(err)
	}
	r := h.Sum64()
	h.Reset()
	hasher.Put(h)
	return r
}
