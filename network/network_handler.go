package networkHandler

import (
	"github.com/GrappigPanda/Olivia/cache"
	"github.com/GrappigPanda/Olivia/config"
	"github.com/GrappigPanda/Olivia/dht"
	"github.com/GrappigPanda/Olivia/network/incoming"
	"github.com/GrappigPanda/Olivia/network/message_handler"
	"log"
	"time"
)

// executeRepeatedly Allows repeated calls to any function which doesn't accept
// arguments. Allows for remote stopping of the execution and passing back
// total number of executions.
func executeRepeatedly(
	sleepDuration time.Duration,
	toExecute func(),
	stopExecution chan interface{},
	responseChannel chan int,
) {
	executionCount := 0

	for {
		select {
		default:
			time.Sleep(sleepDuration)
			toExecute()

			if responseChannel != nil {
				responseChannel <- executionCount
				executionCount++
			}
		case <-stopExecution:
			return
		}
	}
}

// Heartbeat handles time-critical events, such as sending a heartbeat to a
// remote node or expiring keys. heartbeatInterval is the rate at which we need
// to send heartbeat updates to important remote nodes and cycleDuration is the
// rate at which we need to update remote nodes. By default, keys expire every
// second. By default, we send a heartbeat to every important node every second
// on the second. This allows us to asynchronously send our commands and then
// pre-emptively select any keys which will expire the following second.
// Adjusting the heartbeatinterval may have strange, unintended side effects.
func Heartbeat(heartbeatInterval time.Duration, cycleDuration time.Duration) {

}

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
