package lru

import (
	binheap "github.com/GrappigPanda/Olivia/shared"
	"sync"
	"time"
)

// LRUCacheInt32Array is a simple implementation of an LRU cache which will be used in
// the cache based whenever we want to cache values that we don't care too much
// if they're frequently thrown away, so long as the most frequently sought
// keys are preserved within the datastructure.
type LRUCacheInt32Array struct {
	KeyCount    int
	Keys        map[string][]uint32
	KeyTimeouts binheap.BinHeap
	Mutex       *sync.Mutex
}

// New simply allocates a new instance of an LRU cache with `maxEntries` total
// slots.
func NewInt32Array(maxEntries int) *LRUCacheInt32Array {
	return &LRUCacheInt32Array{
		KeyCount:    maxEntries,
		Keys:        make(map[string][]uint32),
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
func (l *LRUCacheInt32Array) Add(key string, value []uint32) ([]uint32, bool) {
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
func (l *LRUCacheInt32Array) Get(key string) ([]uint32, bool) {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()

	value, keyExists := l.Keys[key]
	if keyExists {
		l.KeyTimeouts.UpdateNodeTimeout(key)
	}

	return value, keyExists
}

// RemoveLeastUsed removes the least high prioritized key in the LRU cache.
// Because we use an underlying map of string : uint32 (unix timestamp), we also
// remove any keys from that map, as well.
func (l *LRUCacheInt32Array) RemoveLeastUsed() {
	deletedNode := l.KeyTimeouts.EvictMinNode()
	delete(l.Keys, deletedNode.Key)
}
