package networkHandler

import (
	"github.com/GrappigPanda/Olivia/cache"
	"github.com/GrappigPanda/Olivia/dht"
	"github.com/GrappigPanda/Olivia/network/incoming"
	"github.com/GrappigPanda/Olivia/network/message_handler"
	"log"
	"time"
	"github.com/GrappigPanda/Olivia/config"
)

// StartIncomingNetwork handles spinning up an incoming network router and
// doing any error checking (in the future) as well as ensuring that it
// continues running.
func StartIncomingNetwork(
	mh *message_handler.MessageHandler,
	cache *cache.Cache,
	config *config.Cfg,
) {
	peerList := dht.NewPeerList(mh)
	peer := dht.NewPeerByIP("127.0.0.1:5454", mh)
	peerList.Peers[0] = peer
	(*peerList.PeerMap)["127.0.0.1:5454"] = true

	err := peerList.ConnectAllPeers()
	if err != nil {
		for err != nil {
			log.Println("Sleeping for 60 seconds and attempting to reconnect")
			time.Sleep(time.Second * 60)
			err = peerList.ConnectAllPeers()
		}
	}

	incomingNetwork.StartNetworkRouter(mh, cache, peerList, config)
}
