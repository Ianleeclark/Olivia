package cache

import (
	"fmt"
)

// TODO(ian): Replace this with something else
type Cache struct {
	Cache *map[string]string
}

func NewCache() *Cache {
	cacheMap := make(map[string]string)
	return &Cache{
		&cacheMap,
	}
}

func (c *Cache) Get(key string) (string, error) {
	var value string = ""
	if value, ok := (*c.Cache)[key]; !ok {
		return value, fmt.Errorf("Key not found in cache")
	}

	return value, nil
}

func (c *Cache) Set(key string, value string) error {
	(*c.Cache)[key] = value
	return nil
}
