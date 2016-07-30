package incomingNetwork

import (
	"github.com/GrappigPanda/Olivia/bloomfilter"
	"github.com/GrappigPanda/Olivia/cache"
	"github.com/GrappigPanda/Olivia/chord"
	"github.com/GrappigPanda/Olivia/network/message_handler"
	"github.com/GrappigPanda/Olivia/parser"
	"testing"
)

var CACHE = make(map[string]string)
var READCACHE = make(map[string]string)

var CTX = &ConnectionCtx{
	nil,
	&cache.Cache{
		Cache: &CACHE,
		ReadCache: &READCACHE,
	},
	olilib.NewByFailRate(10000, 0.01),
	message_handler.NewMessageHandler(),
	chord.NewPeerList(),
}

func TestExecuteGetAllSucceed(t *testing.T) {
	expectedReturn := "hash:GOT key1:test1,key2:test14\n"
	expectedReturn2 := "hash:GOT key2:test14,key1:test1\n"

	CACHE["key1"] = "test1"
	CACHE["key2"] = "test14"

	command := parser.CommandData{"hash", "GET", map[string]string{"key1": "", "key2": ""}, nil}
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

	CACHE["key1"] = "test1"
	CACHE["key2"] = "test14"

	command := parser.CommandData{"hash", "GET", map[string]string{"key1": "", "key3": "", "key2": ""}, nil}
	result := CTX.ExecuteCommand(command)

	if expectedReturn != result {
		if result != expectedReturn2 {
			t.Fatalf("Expected <%s> or <%s>, got <%s>", expectedReturn, expectedReturn2, result)
		}
	}
}

func TestExecuteSetKey(t *testing.T) {
	expectedReturn := "hash:SAT key4:test4,key7:test126654\n"
	expectedReturn2 := "hash:SAT key7:test126654,key4:test4\n"

	command := parser.CommandData{"hash", "SET", map[string]string{"key4": "test4", "key7": "test126654"}, nil}
	result := CTX.ExecuteCommand(command)

	if expectedReturn != result {
		if result != expectedReturn2 {
			t.Fatalf("Expected <%s> or <%s>, got <%s>", expectedReturn, expectedReturn2, result)
		}
	}
}

func TestRequestBloomFilter(t *testing.T) {
	bf := olilib.NewByFailRate(10000, 0.01)

	bf.AddKey([]byte("keyalksdjfl"))
	bf.AddKey([]byte("key1"))
	bf.AddKey([]byte("key2"))
	bf.AddKey([]byte("key3"))
	bf.AddKey([]byte("key4"))

	ctx := &ConnectionCtx{
		nil,
		nil,
		bf,
		nil,
		nil,
	}

	command := parser.CommandData{"hash", "REQUEST", map[string]string{"bloomfilter": ""}, nil}
	newBfStr := ctx.ExecuteCommand(command)
	if newBfStr == "Invalid command sent in.\n" {
		t.Fatalf("Sending in a bad command :(")
	}

	newBloomfilter, err := olilib.ConvertStringtoBF(newBfStr)
	if err != nil {
		t.Fatalf("%v", err)
	}

	val, _ := newBloomfilter.HasKey([]byte("key1"))
	if !val {
		t.Fatalf("newBloomfilter doesnt have key1!")
	}

	for i := range bf.Filter {
		if bf.Filter[i] != newBloomfilter.Filter[i] {
			t.Fatalf("Two bfs are not equal")
		}
	}
}
