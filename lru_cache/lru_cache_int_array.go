package olilib_lru

import (
	"sync"
	"time"
	binheap "github.com/GrappigPanda/Olivia/shared"
)

// LRUCacheInt64Array is a simple implementation of an LRU cache which will be used in
// the cache based whenever we want to cache values that we don't care too much
// if they're frequently thrown away, so long as the most frequently sought
// keys are preserved within the datastructure.
type LRUCacheInt64Array struct {
	KeyCount    int
	Keys        map[string][]uint64
	KeyTimeouts *binheap.Heap
	Mutex       *sync.Mutex
}

// New simply allocates a new instance of an LRU cache with `maxEntries` total
// slots.
func NewInt64Array(maxEntries int) *LRUCacheInt64Array {
	return &LRUCacheInt64Array{
		KeyCount:    maxEntries,
		Keys:        make(map[string][]uint64),
		KeyTimeouts: binheap.NewHeap(maxEntries),
		Mutex:       &sync.Mutex{},
	}
}

// Add handles adding keys to the cache and verifying that any values already
// existing in the map are prioritized higher. If too many (max amount) of keys
// are already in the LRU Cache, we will remove the least high prioritized to
// make room for a new key.
// If the return value for the `bool` is false, that means the key was added.
// If the return value for the `bool` is false, that means the key already
// existed in the LRU cache.
func (l *LRUCacheInt64Array) Add(key string, value []uint64) ([]uint64, bool) {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()

	foundValue, keyExists := l.Keys[key]
	if keyExists {
		l.KeyTimeouts.UpdateNodeTimeout(key)
		return foundValue, keyExists
	}

	if len(l.Keys) == l.KeyCount {
		l.RemoveLeastUsed()
	}

	l.Keys[key] = value
	l.KeyTimeouts.Insert(binheap.NewNode(key, time.Now().UTC()))

	return value, false
}

// Get Retrieves a key from the LRU cache and increases its priority.
func (l *LRUCacheInt64Array) Get(key string) ([]uint64, bool) {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()

	value, keyExists := l.Keys[key]
	if keyExists {
		l.KeyTimeouts.UpdateNodeTimeout(key)
	}

	return value, keyExists
}

// RemoveLeastUsed removes the least high prioritized key in the LRU cache.
// Because we use an underlying map of string : uint64 (unix timestamp), we also
// remove any keys from that map, as well.
func (l *LRUCacheInt64Array) RemoveLeastUsed() {
	deletedNode := l.KeyTimeouts.EvictMinNode()
	delete(l.Keys, deletedNode.Key)
}
