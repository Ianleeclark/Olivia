package query_language

import (
        "testing"
)

var CACHE = make(map[string]string)

func TestExecuteGetAllSucceed(t *testing.T) {
        expectedReturn := "GOT key1:test1,key2:test14\n"
        expectedReturn2 := "GOT key2:test14,key1:test1\n"

        CACHE["key1"] = "test1"
        CACHE["key2"] = "test14"

        testCache := &Cache{
                &CACHE,
        }

        result := testCache.ExecuteCommand("GET", map[string]string{"key1": "", "key2": ""})

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

        testCache := &Cache{
                &CACHE,
        }

        result := testCache.ExecuteCommand("GET", map[string]string{"key1": "", "key3": "", "key2": ""})

        if expectedReturn != result {
                if result != expectedReturn2 {
                        t.Fatalf("Expected <%s> or <%s>, got <%s>", expectedReturn, expectedReturn2, result)
                }
        }
}

func TestExecuteSetKey(t *testing.T) {
        expectedReturn := "SAT key4:test4,key7:test126654\n"
        expectedReturn2 := "SAT key7:test126654,key4:test4\n"

        testCache := &Cache{
                &CACHE,
        }

        result := testCache.ExecuteCommand("SET", map[string]string{"key4": "test4", "key7": "test126654"})

        if expectedReturn != result {
                if result != expectedReturn2 {
                        t.Fatalf("Expected <%s> or <%s>, got <%s>", expectedReturn, expectedReturn2, result)
                }
        }
}
