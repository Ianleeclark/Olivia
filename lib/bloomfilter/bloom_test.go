package olilib

import (
        "testing"
)

func TestNewBloomFilter(t *testing.T) {
        expectedReturn := BloomFilter{
                MaxSize: 10000,
                HashFunctions: 3,
                Filter: make([]int, 10000),
        }

        result := New(10000, 3)

        if expectedReturn.MaxSize != result.MaxSize {
                t.Fatalf("Expected %v got %v", expectedReturn.MaxSize, result.MaxSize)
        }

        if expectedReturn.HashFunctions != result.HashFunctions {
                t.Fatalf("Expected %v got %v", expectedReturn.HashFunctions, result.HashFunctions)
        }
}

func TestNewBloomFilterByFailRate(t *testing.T) {
        expectedReturn := BloomFilter{
                MaxSize: 95850,
                HashFunctions: 3,
                Filter: make([]int, 10000),
        }

        result := NewByFailRate(10000, 0.01)

        if expectedReturn.MaxSize != result.MaxSize {
                t.Fatalf("Expected %v got %v", expectedReturn.MaxSize, result.MaxSize)
        }
}

func TestAddKey(t *testing.T) {
        bf := NewByFailRate(10000, 0.01)

        addKeyRet, addIndexes := bf.AddKey([]byte("TestKey"))

        hasKeyRet, hasIndexes := bf.HasKey([]byte("TestKey"))

        if !addKeyRet {
                t.Fatalf("Adding keys failed with indexes %v", addIndexes)
        }

        if !hasKeyRet {
                t.Fatalf("Adding keys failed with indexes %v", addIndexes)
        }

        for index, _ := range hasIndexes {
                if hasIndexes[index] != addIndexes[index] {
                        t.Fatalf("Expected indexes %v, got %v", hasIndexes[index], addIndexes[index])
                }
        }
}
