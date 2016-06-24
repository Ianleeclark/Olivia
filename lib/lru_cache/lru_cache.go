package olilib_lru


import (
        "time"
)

// MAXINT64 Signifies the maximum value for an int64 in Go
var MAXINT64 = int64(1 << 63 - 1)

// LRUCache is a simple implementation of an LRU cache which will be used in
// the cache based whenever we want to cache values that we don't care too much
// if they're frequently thrown away, so long as the most frequently sought
// keys are preserved within the datastructure.
type LRUCache struct {
       KeyCount int
       Keys map[string]string
       KeyTimeouts map[string]int64
}

// New simply allocates a new instance of an LRU cache with `maxEntries` total
// slots.
func New(maxEntries int) *LRUCache {
        return &LRUCache{
                KeyCount: maxEntries,
                Keys: make(map[string]string),
                KeyTimeouts: make(map[string]int64),
        }
}

// Add handles adding keys to the cache and verifying that any values already
// existing in the map are prioritized higher. If too many (max amount) of keys
// are already in the LRU Cache, we will remove the least high prioritized to
// make room for a new key.
func (l *LRUCache) Add(key string, value string) (string, error) {
        value, keyExists := l.Keys[key]
        if keyExists {
                l.KeyTimeouts[key] = getCurrentUnixTime()
                return value, nil
        }

        if(len(l.Keys) >= l.KeyCount) {
                l.RemoveLeastUsed()
        }

        l.Keys[key] = value
        l.KeyTimeouts[key] = getCurrentUnixTime()

        return key, nil
}

// Get Retrieves a key from the LRU cache and increases its priority.
func (l *LRUCache) Get(key string) (string, error) {
        value, keyExists := l.Keys[key]
        if keyExists {
                l.KeyTimeouts[key] = getCurrentUnixTime()
        }

        return value, nil
}

// RemoveLeastUsed removes the least high prioritized key in the LRU cache.
// Because we use an underlying map of string : int64 (unix timestamp), we also
// remove any keys from that map, as well.
func (l *LRUCache) RemoveLeastUsed() {
        var lowest int64
        var lowestKey string

        lowest = MAXINT64

        for k := range l.KeyTimeouts {
                if l.KeyTimeouts[k] < lowest {
                        lowestKey = k
                }
        }

        delete(l.Keys, lowestKey)
        delete(l.KeyTimeouts, lowestKey)
}

func getCurrentUnixTime() int64 {
        // TODO(ian): This needs to return time in nano seconds
        return time.Now().Unix()
}
