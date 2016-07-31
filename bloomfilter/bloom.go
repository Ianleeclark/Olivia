package olilib

import (
	"log"
	"bytes"
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
	Filter        []int
	HashCache     *olilib_lru.LRUCacheInt64Array
}

// New Returns a pointer to a newly allocated `BloomFilter` object
func New(maxSize uint, hashFuns uint) *BloomFilter {
	return &BloomFilter{
		maxSize,
		hashFuns,
		make([]int, maxSize),
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
func (bf *BloomFilter) AddKey(key []byte) (bool, []uint64) {
	hasKey, hashIndexes := bf.HasKey(key)
	if !hasKey {
		hashIndexes = bf.hashKey(key)
	}

	for _, element := range hashIndexes {
		bf.Filter[element] += 1
	}

	return true, hashIndexes
}

// HasKey verifies if a key is or isn't in the bloom filter.
func (bf *BloomFilter) HasKey(key []byte) (bool, []uint64) {
	hashIndexes := bf.hashKey(key)
	for _, element := range hashIndexes {
		if bf.Filter[element] > 0 {
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
	var buffer bytes.Buffer

	for i := range bf.Filter {
		buffer.WriteString(fmt.Sprintf("%v", bf.Filter[i]))
	}

	return Encode(buffer.String())
}

// ConvertStringToBF Decodes the RLE'd bloom filter and then converts it to
// an actual bloom filter in-memory.
func ConvertStringtoBF(inputString string) (*BloomFilter, error) {
	// TODO(ian): Remove this magic number.
	bf := NewByFailRate(1000, 0.01)

	decodedString := Decode(inputString)

	index := 0
	for i, _ := range decodedString {
		number, err := strconv.Atoi(string(decodedString[i]))
		if err != nil {
			return nil, err
		}

		bf.Filter[index] = int(number)
		index++
	}

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
func calculateHash(key []byte, offSet int) uint64 {
	hasher := fnv.New64()
	hasher.Write(key)
	hasher.Write([]byte(strconv.Itoa(offSet)))
	return hasher.Sum64()
}

// hashKey Takes a string in as an argument and hashes it several times to
// create usable indexes for the bloom filter.
func (bf *BloomFilter) hashKey(key []byte) []uint64 {
	if hashes, ok := bf.HashCache.Get(string(key)); ok {
		return hashes
	}

	hashes := make([]uint64, bf.HashFunctions)

	for index, _ := range hashes {
		hashes[index] = calculateHash(key, index) % uint64(bf.MaxSize)
	}

	bf.HashCache.Add(string(key), hashes)
	return hashes
}
