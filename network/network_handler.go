package networkHandler

import (
	"github.com/GrappigPanda/Olivia/cache"
	"github.com/GrappigPanda/Olivia/chord"
	"github.com/GrappigPanda/Olivia/network/incoming"
	"github.com/GrappigPanda/Olivia/network/message_handler"
	"log"
	"time"
)

// StartIncomingNetwork handles spinning up an incoming network router and
// doing any error checking (in the future) as well as ensuring that it
// continues running.
func StartIncomingNetwork(
	mh *message_handler.MessageHandler,
	cache *cache.Cache,
) {
	peerList := chord.NewPeerList()
	err := peerList.ConnectAllPeers()
	if err != nil {
		for err != nil {
			log.Println("Sleeping for 60 seconds and attempting to reconnect")
			time.Sleep(time.Second * 60)
			err = peerList.ConnectAllPeers()
		}
	}

	incomingNetwork.StartNetworkRouter(mh, cache, peerList)
}
