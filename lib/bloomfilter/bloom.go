package olilib

import (
        "hash/fnv"
        "strconv"
        "math"
)

type BloomFilter struct {
        // The maximum size for the bloom filter
        MaxSize uint
        // Total number of hashing functions
        HashFunctions uint
        Filter []int
}

// New Returns a pointer to a newly allocated `BloomFilter` object
func New(maxSize uint, hashFuns uint) *BloomFilter {
        return &BloomFilter{
                maxSize,
                hashFuns,
                make([]int, maxSize),
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
func (bf *BloomFilter) AddKey(key []byte) {
        hashIndexes := bf.hashKey(key)
        for _, element := range hashIndexes {
                bf.Filter[element] += 1
        }
}

// HasKey verifies if a key is or isn't in the bloom filter.
func (bf *BloomFilter) HasKey(key []byte) bool {
        hashIndexes := bf.hashKey(key)
        for _, element := range hashIndexes {
                if bf.Filter[element] > 0 {
                        continue
                } else {
                        return false
                }
        }

        return true
}

// estimateBounds Generates the bounds for total hash function calls and for
// the total bloom filter size
func estimateBounds(maxSize uint, probability float64) (uint, uint) {
        // TODO(ian): Do we need to actually calculate with expected supported
        // values?
        // https://en.wikipedia.org/wiki/Bloom_filter#Counting_filters
        // See "Optimal number of hash functions section"
        n := 1
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
// TODO(ian): Impelement a LRU cache to speed up hashing lookups.
func (bf *BloomFilter) hashKey(key []byte) []uint64 {
        hashes := make([]uint64, bf.HashFunctions)

        for index, _ := range hashes {
                hashes[index] = calculateHash(key, index) % uint64(bf.MaxSize)
        }

        return hashes
}


