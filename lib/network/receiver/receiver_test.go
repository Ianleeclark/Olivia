package network_receiver

import (
        "testing"
        "crypto/md5"
        "encoding/hex"
        "fmt"
)



func TestProcessIncomingString(t *testing.T) {
        var CHANNELMAP = NewChannelMap()
        var RECEIVER = NewReceiver(CHANNELMAP)
        var CHANNEL1 = make(chan string)

        for i := 0; i < 10; i++ {
                x := fmt.Sprintf("test-%v", i)
                hasher := md5.New()
                hasher.Write([]byte(x))

                (*CHANNELMAP.HashLookup)[hex.EncodeToString(hasher.Sum(nil))] = CHANNEL1
        }

        for k, _ := range *CHANNELMAP.HashLookup {
                go RECEIVER.processIncomingString(fmt.Sprintf("%s:%s", k, "TestCommand"))

                x :=  <-CHANNEL1
                if x != "TestCommand" {
                        t.Fatalf("Expected TestCommand, got %v", x)
                }

                RECEIVER.MessageStore.Lock()
                _, ok := (*RECEIVER.MessageStore.HashLookup)[k]
                RECEIVER.MessageStore.Unlock()
                if ok {
                        t.Fatalf("Hash wasn't removed from the `MessageStore` after being replied to")
                }
        }
}
