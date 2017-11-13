package incomingNetwork

import (
	"github.com/GrappigPanda/Olivia/bloomfilter"
	"github.com/GrappigPanda/Olivia/cache"
	"github.com/GrappigPanda/Olivia/parser"
	"testing"
)

var CTX = &ConnectionCtx{
	nil,
	cache.NewCache(nil, nil),
}

func TestExecuteGetAllSucceed(t *testing.T) {
	expectedReturn := "hash:GOT key1:test1,key2:test14\n"
	expectedReturn2 := "hash:GOT key2:test14,key1:test1\n"

	CTX.Cache.Set("key1", "test1")
	CTX.Cache.Set("key2", "test14")

	command := parser.CommandData{"hash", "GET", map[string]string{"key1": "", "key2": ""}, make(map[string]string), nil}
	result := CTX.ExecuteCommand(command)

	if expectedReturn != result {
		if result != expectedReturn2 {
			t.Fatalf("Expected <%v> or <%v>, got <%v>", expectedReturn, expectedReturn2, result)
		}
	}
}

func TestExecuteGetAllSkipNonexistingKey(t *testing.T) {
	expectedReturn := "hash:GOT key1:test1,key2:test14\n"
	expectedReturn2 := "hash:GOT key2:test14,key1:test1\n"

	CTX.Cache.Set("key1", "test1")
	CTX.Cache.Set("key2", "test14")

	command := parser.CommandData{"hash", "GET", map[string]string{"key1": "", "key2": ""}, make(map[string]string), nil}
	result := CTX.ExecuteCommand(command)

	if expectedReturn != result {
		if result != expectedReturn2 {
			t.Fatalf("Expected [%s] or [%s], got [%s]", expectedReturn, expectedReturn2, result)
		}
	}
}

func TestExecuteSetKey(t *testing.T) {
	expectedReturn := "hash:SAT key4:test4,key7:test126654\n"
	expectedReturn2 := "hash:SAT key7:test126654,key4:test4\n"

	command := parser.CommandData{"hash", "SET", map[string]string{"key4": "test4", "key7": "test126654"}, make(map[string]string), nil}
	result := CTX.ExecuteCommand(command)

	if expectedReturn != result {
		if result != expectedReturn2 {
			t.Fatalf("Expected [%s] or [%s], got [%s]", expectedReturn, expectedReturn2, result)
		}
	}
}

func TestExecuteSetKeyWithExpiration(t *testing.T) {
	expectedReturn := "hash:SATEX key1:test1:30,key2:test2:30\n"
	expectedReturn2 := "hash:SATEX key2:test2:30,key1:test1:30\n"

	command := parser.CommandData{"hash", "SETEX", map[string]string{"key1": "test1", "key2": "test2"}, map[string]string{"key1": "30", "key2": "30"}, nil}
	result := CTX.ExecuteCommand(command)

	if expectedReturn != result {
		if result != expectedReturn2 {
			t.Fatalf("Expected [%s] or [%s], got [%s]", expectedReturn, expectedReturn2, result)
		}
	}
}

func TestRequestBloomFilter(t *testing.T) {
	bf := bloomfilter.NewByFailRate(1000, 0.01)

	bf.AddKey([]byte("keyalksdjfl"))
	bf.AddKey([]byte("key1"))
	bf.AddKey([]byte("key2"))
	bf.AddKey([]byte("key3"))
	bf.AddKey([]byte("key4"))

	testCache := cache.NewCache(nil, CONFIG)
	testCache.Set("key1", "value1")
	testCache.Set("keyalksdjfl", "value1")
	testCache.Set("key2", "value2")
	testCache.Set("key3", "value2")
	testCache.Set("key4", "value3")

	ctx := &ConnectionCtx{
		nil,
		testCache,
	}

	command := parser.CommandData{"hash", "REQUEST", map[string]string{"bloomfilter": ""}, make(map[string]string), nil}
	newBfStr := ctx.ExecuteCommand(command)
	if newBfStr == "Invalid command sent in.\n" {
		t.Fatalf("Sending in a bad command :(")
	}

	requestData, _ := parser.NewParser(nil).Parse(newBfStr, nil)

	var bfToParse string
	for k := range requestData.Args {
		bfToParse = k
		break
	}

	newBloomfilter, err := bloomfilter.Deserialize(bfToParse, uint(CONFIG.BloomfilterSize))
	if err != nil {
		t.Fatalf("%v", err)
	}

	val, _ := newBloomfilter.HasKey([]byte("key1"))
	if !val {
		t.Fatalf("newBloomfilter doesnt have key1!")
	}

	if !bf.Compare(newBloomfilter) {
		t.Fatalf("Two bfs are not equal")
	}
}
