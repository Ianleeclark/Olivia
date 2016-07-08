package olilib_network

import (
        "testing"
        "github.com/GrappigPanda/Olivia/lib/bloomfilter"
)

var CACHE = make(map[string]string)

var CTX = &ConnectionCtx{
       nil,
       &Cache{
               &CACHE,
       },
       nil,
}

func TestExecuteGetAllSucceed(t *testing.T) {
        expectedReturn := "GOT key1:test1,key2:test14\n"
        expectedReturn2 := "GOT key2:test14,key1:test1\n"

        CACHE["key1"] = "test1"
        CACHE["key2"] = "test14"

        result := CTX.ExecuteCommand("GET", map[string]string{"key1": "", "key2": ""})

        if expectedReturn != result {
                if result != expectedReturn2 {
                        t.Fatalf("Expected <%s> or <%s>, got <%s>", expectedReturn, expectedReturn2, result)
                }
        }
}

func TestExecuteGetAllSkipNonexistingKey(t *testing.T) {
        expectedReturn := "GOT key1:test1,key2:test14\n"
        expectedReturn2 := "GOT key2:test14,key1:test1\n"


        CACHE["key1"] = "test1"
        CACHE["key2"] = "test14"

        result := CTX.ExecuteCommand("GET", map[string]string{"key1": "", "key3": "", "key2": ""})

        if expectedReturn != result {
                if result != expectedReturn2 {
                        t.Fatalf("Expected <%s> or <%s>, got <%s>", expectedReturn, expectedReturn2, result)
                }
        }
}

func TestExecuteSetKey(t *testing.T) {
        expectedReturn := "SAT key4:test4,key7:test126654\n"
        expectedReturn2 := "SAT key7:test126654,key4:test4\n"

        result := CTX.ExecuteCommand("SET", map[string]string{"key4": "test4", "key7": "test126654"})

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
        }

        new_bf_str := ctx.ExecuteCommand("REQUEST", map[string]string{"bloomfilter": ""})
        if new_bf_str == "Invalid command sent in.\n" {
                t.Fatalf("Sending in a bad command :(")
        }

        new_bf, err := olilib.ConvertStringtoBF(new_bf_str)
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


