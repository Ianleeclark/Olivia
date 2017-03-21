package cache

import (
	"fmt"
	"testing"
	"time"
)

func TestNewCache(t *testing.T) {
	_ = NewCache()
}

func TestSetGet(t *testing.T) {
	cache := NewCache()

	key := "TestKey"
	cache.Set(key, "1024")
	if value, err := cache.Get(key); err != nil || value != "1024" {
		t.Fatalf("expected %v, got %v", "1024", value)
		t.Fatalf("Expected True, got False")
	}

	secondValue, err := cache.Get(key)
	if err != nil {
		t.Fatalf(
			"Got error from GETing key: %v",
			err,
		)
	}

	if secondValue != "1024" {
		t.Fatalf("Expected %v, got %v", "1024", secondValue)
	}
}

func TestCache_SetExpiration(t *testing.T) {
	cache := NewCache()

	key := "TestKey"
	testValue := "1024"

	for i := 0; i < 5; i++ {
		err := cache.SetExpiration(
			fmt.Sprintf("%v-%v",
				key,
				i,
			),
			testValue,
			1,
		)
		if err != nil {
			t.Errorf("Got an error SetExpiration'ing")
		}
	}

	time.Sleep(2 * time.Second)

	cache.EvictExpiredkeys(time.Now().UTC())

	for i := 0; i < 5; i++ {
		value, err := cache.Get(
			fmt.Sprintf("%v-%v",
				key,
				i,
			),
		)

		if err == nil {
			t.Fatalf("Expected err, got %v", value)
		}
	}

}
