package network_receiver

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	. "github.com/GrappigPanda/Olivia/lib/network/outgoing/message_handler"
	"testing"
)

func TestProcessIncomingString(t *testing.T) {
	var CHANNELMAP = NewMessageHandler()
	var RECEIVER = NewReceiver(CHANNELMAP, nil)
	var CHANNEL1 = make(chan string)

	keys := make([]string, 10)
	for i := 0; i < 10; i++ {
		x := fmt.Sprintf("test-%v", i)
		hasher := md5.New()
		hasher.Write([]byte(x))

		hash := hex.EncodeToString(hasher.Sum(nil))
		keys[i] = hash
		(*RECEIVER.MessageStore).AddKeyChannel <- NewKeyValPair(hash, CHANNEL1, nil)
	}

	for i := range keys {
		RECEIVER.processIncomingString(fmt.Sprintf("%s:%s", keys[i], "key"))

	}
}
