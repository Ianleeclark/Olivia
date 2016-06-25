package olilib_lru


import (
        "time"
        "sync"
)

// MAXINT64 Signifies the maximum value for an int64 in Go
var MAXINT64 = int64(1 << 63 - 1)

// LRUCache is a simple implementation of an LRU cache which will be used in
// the cache based whenever we want to cache values that we don't care too much
// if they're frequently thrown away, so long as the most frequently sought
// keys are preserved within the datastructure.
type LRUCacheString struct {
       KeyCount int
       Keys map[string]string
       KeyTimeouts map[string]int64
       Mutex *sync.Mutex
}

// New simply allocates a new instance of an LRU cache with `maxEntries` total
// slots.
func NewString(maxEntries int) *LRUCacheString {
        return &LRUCacheString{
                KeyCount: maxEntries,
                Keys: make(map[string]string),
                KeyTimeouts: make(map[string]int64),
                Mutex: &sync.Mutex{},
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
                l.KeyTimeouts[key] = getCurrentUnixTime()
                return foundValue, true
        }

        if(len(l.Keys) == l.KeyCount) {
                l.RemoveLeastUsed()
        }

        l.Keys[key] = value
        l.KeyTimeouts[key] = getCurrentUnixTime()

        return key, false
}

// Get Retrieves a key from the LRU cache and increases its priority.
func (l *LRUCacheString) Get(key string) (string, bool) {
        l.Mutex.Lock()
        defer l.Mutex.Unlock()

        value, keyExists := l.Keys[key]
        if keyExists {
                l.KeyTimeouts[key] = getCurrentUnixTime()
        }

        return value, keyExists
}

// RemoveLeastUsed removes the least high prioritized key in the LRU cache.
// Because we use an underlying map of string : int64 (unix timestamp), we also
// remove any keys from that map, as well.
func (l *LRUCacheString) RemoveLeastUsed() {
        var lowest int64
        var lowestKey string

        lowest = MAXINT64

        for k := range l.KeyTimeouts {
                if l.KeyTimeouts[k] < lowest {
                        lowestKey = k
                        lowest = l.KeyTimeouts[k]
                }
        }

        delete(l.Keys, lowestKey)
        delete(l.KeyTimeouts, lowestKey)
}

func getCurrentUnixTime() int64 {
        return time.Now().UnixNano()
}
