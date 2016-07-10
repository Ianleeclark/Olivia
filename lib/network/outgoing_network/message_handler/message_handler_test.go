package message_handler

import (
        "testing"
)

var MESSAGEHANDLER = NewMessageHandler()
var RESPONSECHANNEL = make(chan string)
var CALLERRESPONSECHAN = make(chan chan string)

func TestAddKey(t *testing.T) {
        key1 := NewKeyValPair("key1", RESPONSECHANNEL, CALLERRESPONSECHAN)
        keyNoResponseChan := NewKeyValPair("key2", RESPONSECHANNEL, nil)
        key1repeat := NewKeyValPair("key1", RESPONSECHANNEL, CALLERRESPONSECHAN)

        go func() {
                MESSAGEHANDLER.AddKeyChannel <- key1
        }()

        go func() {
                MESSAGEHANDLER.AddKeyChannel <- keyNoResponseChan
        }()

        go func() {
                MESSAGEHANDLER.AddKeyChannel <- key1repeat
        }()
}

func TestRemoveKey(t *testing.T) {
        keyToDelete := NewKeyValPair("keyToDelete", RESPONSECHANNEL, nil)
        MESSAGEHANDLER.AddKeyChannel <- keyToDelete
        keyToDelete2 := NewKeyValPair("keyToDelete", RESPONSECHANNEL, CALLERRESPONSECHAN)
        MESSAGEHANDLER.AddKeyChannel <- keyToDelete2

        MESSAGEHANDLER.RemoveKeyChannel <- keyToDelete
        MESSAGEHANDLER.RemoveKeyChannel <- keyToDelete2
}

func TestRemoveKeyAssertCallerResponse(t *testing.T) {
        responseKey := NewKeyValPair("keyToRespondTo", RESPONSECHANNEL, CALLERRESPONSECHAN)
        MESSAGEHANDLER.AddKeyChannel <- responseKey

        go func() {
                MESSAGEHANDLER.RemoveKeyChannel <- responseKey
        }()

        callerResponse := <-CALLERRESPONSECHAN

        callerResponse <- "testString"

        responseChannelResponse := <-RESPONSECHANNEL

        if responseChannelResponse != "testString" {
                t.Fatalf("expected testString, got %v", responseChannelResponse)
        }
}

func TestRemoveKeyKeyNoExists(t *testing.T) {

}

func TestRemoveNonExistentKeyRespondsWithNil(t *testing.T) {

}
