package bloomfilter

import (
	"github.com/GrappigPanda/Olivia/config"
	"testing"
)

var CONFIG = config.ReadConfig()

func TestNewBloomFilter(t *testing.T) {
	expectedReturn := SimpleBloomFilter{
		maxSize:       uint(CONFIG.BloomfilterSize),
		HashFunctions: 3,
		filter:        NewWFBitset(10000),
	}

	result := NewSimpleBF(uint(CONFIG.BloomfilterSize), 3)

	if expectedReturn.GetMaxSize() != result.GetMaxSize() {
		t.Fatalf("Expected %v got %v", expectedReturn.GetMaxSize(), result.GetMaxSize())
	}

	if expectedReturn.HashFunctions != result.HashFunctions {
		t.Fatalf("Expected %v got %v", expectedReturn.HashFunctions, result.HashFunctions)
	}
}

func TestNewBloomFilterByFailRate(t *testing.T) {
	expectedReturn := SimpleBloomFilter{
		maxSize:       9585,
		HashFunctions: 3,
		filter:        NewWFBitset(10000),
	}

	result := NewByFailRate(uint(CONFIG.BloomfilterSize), 0.01)

	if expectedReturn.GetMaxSize() != result.GetMaxSize() {
		t.Fatalf("Expected %v got %v", expectedReturn.GetMaxSize(), result.GetMaxSize())
	}
}

func TestAddKey(t *testing.T) {
	bf := NewByFailRate(uint(CONFIG.BloomfilterSize), 0.01)

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
	bf := NewByFailRate(uint(CONFIG.BloomfilterSize), 0.01)

	hasKeyRet, _ := bf.HasKey([]byte("TestKey"))

	if hasKeyRet {
		t.Fatalf("Somehow it has the key?")
	}
}

func TestConvertToString(t *testing.T) {
	bf := NewByFailRate(uint(CONFIG.BloomfilterSize), 0.01)

	new_bf_str := bf.Serialize()

	new_bf, err := Deserialize(new_bf_str, uint(CONFIG.BloomfilterSize))
	if err != nil {
		t.Fatalf("%v", err)
	}

	if !new_bf.filter.Compare(bf.filter) {
		t.Fatalf("Two bfs are not equal")
	}
}

func TestConvertWithContainedValues(t *testing.T) {
	bf := NewByFailRate(uint(CONFIG.BloomfilterSize), 0.01)

	bf.AddKey([]byte("keyalksdjfl"))
	bf.AddKey([]byte("key1"))
	bf.AddKey([]byte("key2"))
	bf.AddKey([]byte("key3"))
	bf.AddKey([]byte("key4"))

	new_bf_str := bf.Serialize()

	new_bf, err := Deserialize(new_bf_str, uint(CONFIG.BloomfilterSize))
	if err != nil {
		t.Fatalf("%v", err)
	}

	val, _ := new_bf.HasKey([]byte("key1"))
	if !val {
		t.Fatalf("new_bf doesnt have key1!")
	}

	if !new_bf.filter.Compare(bf.filter) {
		t.Fatalf("Two bfs are not equal")
	}
}
