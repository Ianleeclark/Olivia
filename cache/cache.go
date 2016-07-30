package cache

import (
	"fmt"
	"sync"
)

// TODO(ian): Replace this with something else
// Cache is actually just a map[string]string. Don't tell anyone.
type Cache struct {
	Cache *map[string]string
	readCache *map[string]string
	sync.Mutex
}

// NewCache creates a new cache and internal readCache.
func NewCache() *Cache {
	cacheMap := make(map[string]string)
	writeCache := make(map[string]string)
	return &Cache{
		Cache: &cacheMap,
		readCache: &writeCache,
	}
}

// Get handles retrieving a value by its key from the internal cache. It reads
// from the readCache which is for copy-on-write optimizations so that
// reading doesn't lock the cache.
func (c *Cache) Get(key string) (string, error) {
	var value string
	if value, ok := (*c.readCache)[key]; !ok {
		return value, fmt.Errorf("Key not found in cache")
	}

	return value, nil
}


// copyCache handles creating a copy of the cache
func (c *Cache) copyCache() {
	c.Lock()
	for k, v := range (*c.Cache) {
		(*c.readCache)[k] = v
	}
	c.Unlock()
}

// Set handles adding a key/value pair to the cache and updating the internal
// readCache.
func (c *Cache) Set(key string, value string) error {
	c.Lock()
	(*c.Cache)[key] = value
	c.Unlock()
	go c.copyCache()

	return nil
}
