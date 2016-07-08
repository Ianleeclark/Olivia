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

func TestHasKeyFailNoKey(t *testing.T) {
        bf := NewByFailRate(10000, 0.01)

        hasKeyRet, _ := bf.HasKey([]byte("TestKey"))

        if hasKeyRet {
                t.Fatalf("Somehow it has the key?")
        }
}

func TestConvertToString(t *testing.T) {
        bf := NewByFailRate(10000, 0.01)

        new_bf_str := bf.ConvertToString()

        new_bf, err := ConvertStringtoBF(new_bf_str)
        if err != nil {
                t.Fatalf("%v", err)
        }

        for i := range bf.Filter {
                if bf.Filter[i] != new_bf.Filter[i] {
                        t.Fatalf("Two bfs are not equal")
                }
        }
}

func TestConvertWithContainedValues(t *testing.T) {
        bf := NewByFailRate(10000, 0.01)

        bf.AddKey([]byte("keyalksdjfl"))
        bf.AddKey([]byte("key1"))
        bf.AddKey([]byte("key2"))
        bf.AddKey([]byte("key3"))
        bf.AddKey([]byte("key4"))

        new_bf_str := bf.ConvertToString()

        new_bf, err := ConvertStringtoBF(new_bf_str)
        if err != nil {
                t.Fatalf("%v", err)
        }

        val, _ := new_bf.HasKey([]byte("key1"))
        if !val {
                t.Fatalf("new_bf doesnt have key1!")
        }

        for i := range bf.Filter {
                if bf.Filter[i] != new_bf.Filter[i] {
                        t.Fatalf("Two bfs are not equal")
                }
        }
}

