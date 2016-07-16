package requester

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/GrappigPanda/Olivia/lib/network/outgoing"
	. "github.com/GrappigPanda/Olivia/lib/network/outgoing/message_handler"
	"github.com/GrappigPanda/Olivia/lib/network/outgoing/receiver"
)

// SendRequest handles taking in a peer object and a command and sending a
// command which will be responded to the calling channel once the request has
// been fulfilled
func SendRequest(peer *peer.Peer, Command string, responseChannel chan string, mh *MessageHandler) {
	receiver := network_receiver.NewReceiver(mh, peer.Conn)

	hash := hashRequest(Command)
	addCommandToMessageHandler(hash, responseChannel, mh)

	go func() {
		receiver.Run()
	}()

	(*peer.Conn).Write([]byte(fmt.Sprintf("%s:%s", hash, Command)))
}

// addCommandToMessageHandler send a command to the message container to store
// the callback channel.
func addCommandToMessageHandler(hash string, responseChannel chan string, mh *MessageHandler) {
	keyVal := NewKeyValPair(hash, responseChannel, nil)

	mh.AddKeyChannel <- keyVal
}

// hashRequest hashes the command so that later the channel can be responded to
// from the message container
func hashRequest(Command string) string {
	hasher := md5.New()
	hasher.Write([]byte(Command))

	return hex.EncodeToString(hasher.Sum(nil))
}
