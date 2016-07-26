package message_handler

import (
	"fmt"
	"testing"
	"time"
)

var MESSAGEHANDLER = NewMessageHandler()
var RESPONSECHANNEL = make(chan string)
var CALLERRESPONSECHAN = make(chan chan string)

func TestAddKey(t *testing.T) {
	key1 := NewKeyValPair("key1", RESPONSECHANNEL, CALLERRESPONSECHAN)
	keyNoResponseChan := NewKeyValPair("key2", RESPONSECHANNEL, nil)
	key1repeat := NewKeyValPair("key1", RESPONSECHANNEL, CALLERRESPONSECHAN)

	MESSAGEHANDLER.AddKeyChannel <- key1
	MESSAGEHANDLER.AddKeyChannel <- keyNoResponseChan
	MESSAGEHANDLER.AddKeyChannel <- key1repeat

	time.Sleep(1 * time.Second)

	MESSAGEHANDLER.Lock()
	if _, keyExists := (*MESSAGEHANDLER.messageResponseStore)["key1"]; !keyExists {
		t.Fatalf("Expected to find key key1, no key exists!")
		fmt.Println(*MESSAGEHANDLER.messageResponseStore)
	}

	if _, keyExists := (*MESSAGEHANDLER.messageResponseStore)["key2"]; !keyExists {
		t.Fatalf("Expected to find key key2, no key exists!")
		fmt.Println(*MESSAGEHANDLER.messageResponseStore)
	}
	MESSAGEHANDLER.Unlock()
}

func TestRemoveKey(t *testing.T) {
	keyToDelete := NewKeyValPair("keyToDelete", RESPONSECHANNEL, nil)
	MESSAGEHANDLER.AddKeyChannel <- keyToDelete

	keyToDelete2 := NewKeyValPair("keyToDelete2", RESPONSECHANNEL, CALLERRESPONSECHAN)
	MESSAGEHANDLER.AddKeyChannel <- keyToDelete2

	time.Sleep(1 * time.Second)

	MESSAGEHANDLER.Lock()
	if _, keyExists := (*MESSAGEHANDLER.messageResponseStore)["keyToDelete"]; !keyExists {
		t.Fatalf("Expected to find key keyToDelete, no key exists!")
		fmt.Println(*MESSAGEHANDLER.messageResponseStore)
	}

	if _, keyExists := (*MESSAGEHANDLER.messageResponseStore)["keyToDelete2"]; !keyExists {
		t.Fatalf("Expected to find key keyToDelete2, no key exists!")
		fmt.Println(*MESSAGEHANDLER.messageResponseStore)
	}
	MESSAGEHANDLER.Unlock()

	MESSAGEHANDLER.RemoveKeyChannel <- keyToDelete
	MESSAGEHANDLER.RemoveKeyChannel <- keyToDelete2

	time.Sleep(2 * time.Second)

	MESSAGEHANDLER.Lock()
	if _, keyExists := (*MESSAGEHANDLER.messageResponseStore)["keyToDelete"]; keyExists {
		t.Fatalf("Expected to not find key keyToDelete, no key exists!")
		fmt.Println(*MESSAGEHANDLER.messageResponseStore)
	}

	if _, keyExists := (*MESSAGEHANDLER.messageResponseStore)["keyToDelete2"]; keyExists {
		t.Fatalf("Expected to not find key keyToDelete2, no key exists!")
		fmt.Println(*MESSAGEHANDLER.messageResponseStore)
	}
	MESSAGEHANDLER.Unlock()
}

func TestRemoveKeyAssertCallerResponse(t *testing.T) {
	endResponseChannel := make(chan string)
	callbackChan := make(chan chan string)

	keyToDelete := NewKeyValPair("keyToRespondTo", endResponseChannel, callbackChan)
	MESSAGEHANDLER.AddKeyChannel <- keyToDelete

	time.Sleep(1 * time.Second)

	MESSAGEHANDLER.Lock()
	if _, keyExists := (*MESSAGEHANDLER.messageResponseStore)["keyToRespondTo"]; !keyExists {
		t.Fatalf("Expected to find key keyToRespondTo, no key exists!")
		fmt.Println(*MESSAGEHANDLER.messageResponseStore)
	}
	MESSAGEHANDLER.Unlock()

	MESSAGEHANDLER.RemoveKeyChannel <- keyToDelete

	endChannel := <-callbackChan
	go func() {
		returnValue := <-endResponseChannel
		if returnValue != "testString" {
			t.Fatalf("Expected testString, got %v\n", returnValue)
		}
	}()

	endChannel <- "testString"
}

func TestRemoveKeyKeyNoExistsRespondsNil(t *testing.T) {
	endResponseChannel := make(chan string)
	callbackChan := make(chan chan string)
	keyToDelete := NewKeyValPair("KeyThatDoesntExist", endResponseChannel, callbackChan)

	MESSAGEHANDLER.RemoveKeyChannel <- keyToDelete

	endChannel := <-callbackChan
	if endChannel != nil {
		t.Fatalf("Expected nil, got %v\n", endChannel)
	}
}
