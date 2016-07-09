package network_receiver

import (
        "testing"
        "crypto/md5"
        "encoding/hex"
        "fmt"
)


var CHANNELMAP = NewChannelMap()
var RECEIVER = NewReceiver(CHANNELMAP)
var CHANNEL1 = make(chan string)

func TestProcessIncomingString(t *testing.T) {
        for k, _ := range *CHANNELMAP {
                RECEIVER.processIncomingString(fmt.Sprintf("%s:%s", k, "TestCommand"))

                x :=  <-CHANNEL1
                if x != "TestCommand" {
                        t.Fatalf("Expected TestCommand, got %v", x)
                }
        }
}

func TestMain(m *testing.M) {
        for i := 0; i < 10; i++ {
                x := fmt.Sprintf("test-%v", i)
                hasher := md5.New()
                hasher.Write([]byte(x))

                (*CHANNELMAP)[hex.EncodeToString(hasher.Sum(nil))] = CHANNEL1
        }
}
