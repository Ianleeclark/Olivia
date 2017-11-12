package lru

import (
	binheap "github.com/GrappigPanda/Olivia/shared"
	"sync"
	"testing"
	"time"
)

var TESTLRUINT32 = NewInt32Array(10)

func TestNewInt32(t *testing.T) {
	expectedReturn := &LRUCacheInt32Array{
		10,
		make(map[string][]uint32, 10),
		binheap.NewHeap(10),
		&sync.Mutex{},
	}

	result := NewString(10)

	if expectedReturn.KeyCount != result.KeyCount {
		t.Fatalf("Expected %v, got %v", expectedReturn.KeyCount, result.KeyCount)
	}
}

func TestAddInt32(t *testing.T) {
	expectedReturn := []uint32{3, 5, 1}

	value, err := TESTLRUINT32.Add("Key", expectedReturn)

	for i := range value {
		if value[i] != expectedReturn[i] {
			t.Fatalf("Failed adding a key to the LRU cache")
		}
	}

	if err != false {
		t.Fatalf("Weird things happened adding a key")
	}
}

func TestAddPreExistingKeyInt32(t *testing.T) {
	expectedReturn := []uint32{4, 7, 9}
	_, ok := TESTLRUINT32.Add("Key", expectedReturn)

	if ok != true {
		t.Fatalf("Failed to add pre-existing key. Keys: %v",
			TESTLRUINT32.Keys,
		)
	}
}

func TestAddRemoveOldestInt32(t *testing.T) {
	testLRU := NewInt32Array(10)
	testLRU.Add("Key1", []uint32{1, 2, 3})
	time.Sleep(time.Nanosecond * 1)
	testLRU.Add("Key2", []uint32{1, 2, 3})
	testLRU.Add("Key3", []uint32{1, 2, 3})

	expectedKeys := []string{"Key1", "Key2", "Key3"}

	for i := range expectedKeys {
		if _, ok := testLRU.Keys[expectedKeys[i]]; !ok {
			t.Fatalf("Key %v not found in the test LRU (%v)",
				expectedKeys[i],
				testLRU.Keys,
			)
		}

		if _, ok := testLRU.KeyTimeouts.Get(expectedKeys[i]); !ok {
			t.Fatalf("Key %v not found in the test LRU (KeyTimeouts)", expectedKeys[i])
		}
	}

	testLRU.Add("Key4", []uint32{12, 13, 14})

	expectedKeys = []string{"Key4", "Key2", "Key3"}

	for i := range expectedKeys {
		if _, ok := testLRU.Keys[expectedKeys[i]]; !ok {
			t.Fatalf("Key %v not found in the test LRU (Keys)", expectedKeys[i])
		}

		if _, ok := testLRU.KeyTimeouts.Get(expectedKeys[i]); !ok {
			t.Fatalf("Key %v not found in the test LRU (KeyTimeouts)", expectedKeys[i])
		}
	}

}

func TestGetInt32(t *testing.T) {
	expectedReturn := []uint32{5, 7, 9}
	testLRU := NewInt32Array(5)
	testLRU.Add("Key1", expectedReturn)

	node, _ := testLRU.KeyTimeouts.Get("Key1")
	originalTime := node.Timeout
	time.Sleep(5 * time.Millisecond)
	value, keyExists := testLRU.Get("Key1")

	if keyExists != true {
		t.Fatalf("Key doesn't exist in the LRU Cache")
	}

	for i := range value {
		if value[i] != expectedReturn[i] {
			t.Fatalf("Failed adding a key to the LRU cache")
		}
	}

	if node.Timeout == originalTime {
		t.Fatalf("Time for retrieving a key didnt update, please fix.")
	}
}

func TestGetDoesntExistInt32(t *testing.T) {
	testLRU := NewInt32Array(5)

	_, keyExists := testLRU.Get("Key14")

	if keyExists != false {
		t.Fatalf("Key exists somehow in the array")
	}
}
