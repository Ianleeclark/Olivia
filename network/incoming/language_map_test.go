package incomingNetwork

import (
	"testing"

	"github.com/GrappigPanda/Olivia/cache"
	"github.com/GrappigPanda/Olivia/dht"
	"github.com/GrappigPanda/Olivia/network/message_handler"
	"github.com/GrappigPanda/Olivia/parser"
)

var MESSAGEBUS = message_handler.NewMessageHandler()

var CTX = &ConnectionCtx{
	nil,
	cache.NewCache(),
	olilib.NewByFailRate(1000, 0.01),
	MESSAGEBUS,
	dht.NewPeerList(MESSAGEBUS, *CONFIG),
}

func TestExecuteGetAllSucceed(t *testing.T) {
	expectedReturn := "hash:GOT key1:test1,key2:test14\n"
	expectedReturn2 := "hash:GOT key2:test14,key1:test1\n"

	(*CTX.Cache.Cache)["key1"] = "test1"
	(*CTX.Cache.Cache)["key2"] = "test14"

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

	(*CTX.Cache.Cache)["key1"] = "test1"
	(*CTX.Cache.Cache)["key2"] = "test14"

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
	bf := olilib.NewByFailRate(1000, 0.01)

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

	command := parser.CommandData{"hash", "REQUEST", map[string]string{"bloomfilter": ""}, make(map[string]string), nil}
	newBfStr := ctx.ExecuteCommand(command)
	if newBfStr == "Invalid command sent in.\n" {
		t.Fatalf("Sending in a bad command :(")
	}

	requestData, _ := parser.NewParser(nil).Parse(newBfStr, nil)

	var bfToParse string
	for k, _ := range requestData.Args {
		bfToParse = k
		break
	}

	newBloomfilter, err := olilib.ConvertStringtoBF(bfToParse, uint(CONFIG.BloomfilterSize))
	if err != nil {
		t.Fatalf("%v", err)
	}

	val, _ := newBloomfilter.HasKey([]byte("key1"))
	if !val {
		t.Fatalf("newBloomfilter doesn't have key1!")
	}

	if !bf.Filter.BS.Equal(newBloomfilter.Filter.BS) {
		t.Fatalf("Two bfs are not equal")
	}
}
