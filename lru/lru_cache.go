package lru

import (
	"github.com/GrappigPanda/Olivia/binheap"
	"github.com/GrappigPanda/Olivia/binheap/binheapv1"
	"sync"
	"time"
)

// MAXINT64 Signifies the maximum value for an int64 in Go
var MAXINT64 = int64(1<<63 - 1)

// LRUCache is a simple implementation of an LRU cache which will be used in
// the cache based whenever we want to cache values that we don't care too much
// if they're frequently thrown away, so long as the most frequently sought
// keys are preserved within the datastructure.
type LRUCacheString struct {
	KeyCount    int
	Keys        map[string]string
	KeyTimeouts binheap.LRUStorage
	Mutex       *sync.Mutex
}

// New simply allocates a new instance of an LRU cache with `maxEntries` total
// slots.
func NewString(maxEntries int) *LRUCacheString {
	return &LRUCacheString{
		KeyCount:    maxEntries,
		Keys:        make(map[string]string),
		KeyTimeouts: binheapv1.NewHeap(maxEntries),
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
func (l *LRUCacheString) Add(key string, value string) (string, bool) {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()

	foundValue, keyExists := l.Keys[key]
	if keyExists {
		l.KeyTimeouts.UpdateNodeTimeout(key)
		return foundValue, true
	}

	if len(l.Keys) == l.KeyCount {
		l.RemoveLeastUsed()
	}

	l.Keys[key] = value
	l.addNewKeyTimeout(key)

	return key, false
}

// addNewKeyTimeout handles adding a key into our priority queue for later
// eviction.
func (l *LRUCacheString) addNewKeyTimeout(key string) {
	l.KeyTimeouts.Insert(binheap.NewNode(key, getCurrentUnixTime()))
}

// Get Retrieves a key from the LRU cache and increases its priority.
func (l *LRUCacheString) Get(key string) (string, bool) {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()

	value, keyExists := l.Keys[key]
	if keyExists {
		l.KeyTimeouts.UpdateNodeTimeout(key)
	}

	return value, keyExists
}

// RemoveLeastUsed removes the least high prioritized key in the LRU cache.
// Because we use an underlying map of string : int64 (unix timestamp), we also
// remove any keys from that map, as well.
func (l *LRUCacheString) RemoveLeastUsed() {
	deletedNode := l.KeyTimeouts.EvictMinNode()
	delete(l.Keys, deletedNode.Key)
}

func getCurrentUnixTime() time.Time {
	return time.Now().UTC()
}
