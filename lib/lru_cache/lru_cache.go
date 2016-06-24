package olilib_lru


import (
        "time"
)

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

func (l *LRUCache) Get(key string) (string, error) {
        value, keyExists := l.Keys[key]
        if keyExists {
                l.KeyTimeouts[key] = getCurrentUnixTime()
        }

        return value, nil
}

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
        return time.Now().Unix()
}
