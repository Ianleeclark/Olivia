package olilib_lru

import (
        "testing"
)

var TESTLRU = New(10)

func TestNew(t *testing.T) {
        expectedReturn := &LRUCache{
                10,
                make(map[string]string, 10),
                make(map[string]int64, 10),
        }

        result := New(10)

        if expectedReturn.KeyCount != result.KeyCount {
                t.Fatalf("Expected %v, got %v", expectedReturn.KeyCount, result.KeyCount)
        }
}

func TestAdd(t *testing.T) {
        key, err := TESTLRU.Add("Key", "Value")

        if key != "Key" {
                t.Fatalf("Failed adding a key to the LRU cache")
        }

        if err != nil {
                t.Fatalf("Weird things happened adding a key")
        }
}

func TestAddRemoveOldest(t *testing.T) {
        testLRU := New(3)
        testLRU.Add("Key1", "value1")
        testLRU.Add("Key2", "value2")
        testLRU.Add("Key3", "value3")

        expectedKeys := []string{"Key1", "Key2", "Key3"}

        for i := range expectedKeys {
                if _, ok := testLRU.Keys[expectedKeys[i]]; !ok {
                        t.Fatalf("Key %v not found in the test LRU (Keys)", expectedKeys[i])
                }

                if _, ok := testLRU.KeyTimeouts[expectedKeys[i]]; !ok {
                        t.Fatalf("Key %v not found in the test LRU (KeyTimeouts)", expectedKeys[i])
                }
        }
}

func TestRemoveLeastUsed(t *testing.T) {

}

func TestGet(t *testing.T) {

}

func TestGetDoesntExist(t *testing.T) {
}
