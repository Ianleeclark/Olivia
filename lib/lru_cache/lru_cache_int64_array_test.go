package olilib_lru

import (
        "testing"
        "time"
        "sync"
)

var TESTLRUINT64 = NewInt64Array(10)

func TestNewInt64(t *testing.T) {
        expectedReturn := &LRUCacheInt64Array{
                10,
                make(map[string][]uint64, 10),
                make(map[string]int64, 10),
                &sync.Mutex{},
        }

        result := NewString(10)

        if expectedReturn.KeyCount != result.KeyCount {
                t.Fatalf("Expected %v, got %v", expectedReturn.KeyCount, result.KeyCount)
        }
}

func TestAddInt64(t *testing.T) {
        expectedReturn := []uint64{3, 5, 1}

        value, err := TESTLRUINT64.Add("Key", expectedReturn)

        for i := range value {
               if value[i] != expectedReturn[i] {
                        t.Fatalf("Failed adding a key to the LRU cache")
               }
        }

        if err != false {
                t.Fatalf("Weird things happened adding a key")
        }
}

func TestAddPreExistingKeyInt64(t *testing.T) {
        expectedReturn := []uint64{4, 7, 9}
        _, err := TESTLRUINT64.Add("Key", expectedReturn)

        if err != true {
                t.Fatalf("Failed to add pre-existing key.")
        }
}

func TestAddRemoveOldestInt64(t *testing.T) {
        testLRU := NewInt64Array(10)
        testLRU.Add("Key1", []uint64{1, 2, 3})
        time.Sleep(time.Nanosecond * 1)
        testLRU.Add("Key2", []uint64{1, 2, 3})
        testLRU.Add("Key3",  []uint64{1, 2, 3})

        expectedKeys := []string{"Key1", "Key2", "Key3"}

        for i := range expectedKeys {
                if _, ok := testLRU.Keys[expectedKeys[i]]; !ok {
                        t.Fatalf("Key %v not found in the test LRU (Keys)", expectedKeys[i])
                }

                if _, ok := testLRU.KeyTimeouts[expectedKeys[i]]; !ok {
                        t.Fatalf("Key %v not found in the test LRU (KeyTimeouts)", expectedKeys[i])
                }
        }

        testLRU.Add("Key4", []uint64{12, 13, 14})

        expectedKeys = []string{"Key4", "Key2", "Key3"}

        for i := range expectedKeys {
                if _, ok := testLRU.Keys[expectedKeys[i]]; !ok {
                        t.Fatalf("Key %v not found in the test LRU (Keys)", expectedKeys[i])
                }

                if _, ok := testLRU.KeyTimeouts[expectedKeys[i]]; !ok {
                        t.Fatalf("Key %v not found in the test LRU (KeyTimeouts)", expectedKeys[i])
                }
        }


}

func TestGetInt64(t *testing.T) {
        expectedReturn := []uint64{5, 7, 9}
        testLRU := NewInt64Array(5)
        testLRU.Add("Key1", expectedReturn)

        originalTime := testLRU.KeyTimeouts["Key1"]
        value, keyExists := testLRU.Get("Key1")

        if keyExists != true {
                t.Fatalf("Key doesn't exist in the LRU Cache")
        }

        for i := range value {
               if value[i] != expectedReturn[i] {
                        t.Fatalf("Failed adding a key to the LRU cache")
               }
        }

        if testLRU.KeyTimeouts["Key1"] == originalTime {
                t.Fatalf("Time for retrieving a key didnt update, please fix.")
        }
}

func TestGetDoesntExistInt64(t *testing.T) {
        testLRU := NewInt64Array(5)

        _, keyExists := testLRU.Get("Key14")

        if keyExists != false {
                t.Fatalf("Key exists somehow in the array")
        }
}
