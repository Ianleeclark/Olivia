package cache

import (
	"time"
	"fmt"
	"sync"
	binheap "github.com/GrappigPanda/Olivia/shared"
)

// TODO(ian): Replace this with something else
// Cache is actually just a map[string]string. Don't tell anyone.
type Cache struct {
	Cache *map[string]string
	ReadCache *map[string]string
	binHeap *binheap.Heap
	sync.Mutex
}

// NewCache creates a new cache and internal ReadCache.
func NewCache() *Cache {
	cacheMap := make(map[string]string)
	writeCache := make(map[string]string)
	return &Cache{
		Cache: &cacheMap,
		ReadCache: &writeCache,
		binHeap: binheap.NewHeapReallocate(100),
	}
}

// Get handles retrieving a value by its key from the internal cache. It reads
// from the ReadCache which is for copy-on-write optimizations so that
// reading doesn't lock the cache.
func (c *Cache) Get(key string) (string, error) {
	if value, ok := (*c.ReadCache)[key]; !ok {
		return "", fmt.Errorf("Key not found in cache")
	} else {
		return value, nil
	}
}


// copyCache handles creating a copy of the cache
func (c *Cache) copyCache() {
	c.Lock()
	for k, v := range (*c.Cache) {
		(*c.ReadCache)[k] = v
	}
	c.Unlock()
}

// Set handles adding a key/value pair to the cache and updating the internal
// ReadCache.
func (c *Cache) Set(key string, value string) error {
	c.Lock()
	(*c.Cache)[key] = value
	c.Unlock()

	c.copyCache()

	return nil
}

// SetExpiration handles setting a key with an expiration time.
func (c *Cache) SetExpiration(key string, value string, timeout int) error {
	c.binHeap.Insert(binheap.NewNode(key, time.Now().UTC()))
	err := c.Set(key, value)

	c.copyCache()
	return err
}

// EvictExpiredKeys handles
func (c *Cache) EvictExpiredkeys(expirationDate time.Time) {
	c.Lock()
	keysToExpire := make([]string, len(c.binHeap.Tree))

	i := 0
	for {
		node, err := c.binHeap.Peek(i)
		if err != nil {
			break
		}

		if node.Timeout.Second() > expirationDate.Second() {
			break
		} else {
			keysToExpire = append(keysToExpire, node.Key)
		}

		i++
	}

	for _, key := range keysToExpire {
		delete((*c.Cache), key)
	}
	c.Unlock()
}
