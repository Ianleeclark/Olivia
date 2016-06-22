package olilib_lru


import (

)

// LRUCache is a simple implementation of an LRU cache which will be used in
// the cache based whenever we want to cache values that we don't care too much
// if they're frequently thrown away, so long as the most frequently sought
// keys are preserved within the datastructure.
type LRUCache struct {
       KeyCount int
       Keys []interface{}
}

func New(maxEntries int) *LRUCache {
        return &LRUCache{
                KeyCount: maxEntries,
                Keys: make([]interface{}, maxEntries),
        }
}

func (l *LRUCache) Add(key interface{}) {

}

func (l *LRUCache) Get(key interface{}) {

}

func (l *LRUCache) RemoveLeastUsed() {

}
