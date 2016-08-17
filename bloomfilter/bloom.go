package olilib

import (
	"fmt"
	"github.com/GrappigPanda/Olivia/lru_cache"
	"hash/fnv"
	"math"
	"strconv"
)

type BloomFilter struct {
	// The maximum size for the bloom filter
	MaxSize uint
	// Total number of hashing functions
	HashFunctions uint
	Filter        *Bitset
	HashCache     *olilib_lru.LRUCacheInt64Array
}

// New Returns a pointer to a newly allocated `BloomFilter` object
func New(maxSize uint, hashFuns uint) *BloomFilter {
	return &BloomFilter{
		maxSize,
		hashFuns,
		NewBitset(maxSize),
		olilib_lru.NewInt64Array(int((float64(maxSize) * float64(0.1)))),
	}
}

// NewByFailRate allows generation of a bloom filter with a pre-conceived
// amount of items and a false-positive failure rate. We calculate our bloom
// filter bounds and generate the new bloom filter this way.
func NewByFailRate(items uint, probability float64) *BloomFilter {
	m, k := estimateBounds(items, probability)
	return New(m, k)
}

// AddKey Adds a new key to the bloom filter
func (bf *BloomFilter) AddKey(key []byte) (bool, []uint) {
	hasKey, hashIndexes := bf.HasKey(key)
	if !hasKey {
		hashIndexes = bf.hashKey(key)
	}

	for _, index := range hashIndexes {
		bf.Filter.Add(index)
	}

	return true, hashIndexes
}

// HasKey verifies if a key is or isn't in the bloom filter.
func (bf *BloomFilter) HasKey(key []byte) (bool, []uint) {
	hashIndexes := bf.hashKey(key)

	for _, element := range hashIndexes {
		if bf.Filter.Contains(element) {
			continue
		} else {
			return false, nil
		}
	}

	return true, hashIndexes
}

// ConvertToString handles conversion of a bloom filter to a string. Moreover,
// it enforces RLE encoding, so that fewer bytes are transferred per request.
func (bf *BloomFilter) ConvertToString() string {
	return Encode(bf.Filter.ToString())
}

// ConvertStringToBF Decodes the RLE'd bloom filter and then converts it to
// an actual bloom filter in-memory.
func ConvertStringtoBF(inputString string) (*BloomFilter, error) {
	// TODO(ian): Remove this magic number.
	bf := NewByFailRate(1000, 0.01)

	sz := fmt.Sprintf("\"%s=\"", Decode(inputString))
	bf.Filter.FromString(sz)

	return bf, nil
}

// estimateBounds Generates the bounds for total hash function calls and for
// the total bloom filter size
func estimateBounds(items uint, probability float64) (uint, uint) {
	// https://en.wikipedia.org/wiki/Bloom_filter#Counting_filters
	// See "Optimal number of hash functions section"
	n := items
	m := (-1 * float64(n) * math.Log(probability)) / (math.Pow(math.Log(2), 2))
	k := uint((m / float64(n)) * math.Log(2))

	return uint(m), k
}

// calculateHash Takes in a string and calculates the 64bit hash value.
func calculateHash(key []byte, offSet int) uint {
	hasher := fnv.New32()
	hasher.Write(key)
	hasher.Write([]byte(strconv.Itoa(offSet)))
	return uint(hasher.Sum32())
}

// hashKey Takes a string in as an argument and hashes it several times to
// create usable indexes for the bloom filter.
func (bf *BloomFilter) hashKey(key []byte) []uint {
	hashes := make([]uint, bf.HashFunctions)

	for index, _ := range hashes {
		hashes[index] = calculateHash(key, index) % uint(bf.MaxSize)
	}

	return hashes
}
