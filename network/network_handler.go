package networkHandler

import (
	"github.com/GrappigPanda/Olivia/cache"
	"github.com/GrappigPanda/Olivia/config"
	"github.com/GrappigPanda/Olivia/dht"
	"github.com/GrappigPanda/Olivia/network/incoming"
	"github.com/GrappigPanda/Olivia/network/message_handler"
	"time"
) // executeRepeatedly Allows repeated calls to any function which doesn't accept
// arguments. Allows for remote stopping of the execution and passing back
// total number of executions.
func executeRepeatedly(
	sleepDuration time.Duration,
	toExecute func(),
	stopExecution chan bool,
	responseChannel chan int,
) {
	for {
		select {
		default:
			time.Sleep(sleepDuration)
			toExecute()

			if responseChannel != nil {
				responseChannel <- 1
			}
			break
		case <-stopExecution:
			return
		}
	}
}

// heartbeatRemoteNodes handles sending a heartbeat to every node in a peer
// list.
func heartbeatRemoteNodes(interval time.Duration, peerLists ...[]*dht.Peer) {
	executeRepeatedly(
		interval,
		func() {
			for _, peerList := range peerLists {
				for peer := range peerList {
					if peerList[peer] != nil {
						go peerList[peer].TestConnection()
					}
				}
			}
		},
		nil,
		nil,
	)
}

// getRemoteBloomFilters requests a remote node's peer list on a timed
// interval.
func getRemoteBloomFilters(interval time.Duration, peerLists ...[]*dht.Peer) {
	executeRepeatedly(
		interval,
		func() {
			for _, peerList := range peerLists {
				for peer := range peerList {
					if peerList[peer] != nil {
						go peerList[peer].GetBloomFilter()
					}
				}
			}
		},
		nil,
		nil,
	)
}

// Heartbeat handles time-critical events, such as sending a heartbeat to a
// remote node or expiring keys. heartbeatInterval is the rate at which we need
// to send heartbeat updates to important remote nodes and cycleDuration is the
// rate at which we need to update remote nodes. By default, keys expire every
// second. By default, we send a heartbeat to every important node every second
// on the second. This allows us to asynchronously send our commands and then
// pre-emptively select any keys which will expire the following second.
// Adjusting the heartbeatinterval may have strange, unintended side effects.
func Heartbeat(
	heartbeatInterval time.Duration,
	cycleDuration time.Duration,
	peerList *dht.PeerList,
) {
	go heartbeatRemoteNodes(
		heartbeatInterval,
		peerList.Peers,
		peerList.BackupPeers,
	)
	go getRemoteBloomFilters(
		cycleDuration,
		peerList.Peers,
		peerList.BackupPeers,
	)
}

// StartIncomingNetwork handles spinning up an incoming network router and
// doing any error checking (in the future) as well as ensuring that it
// continues running.
func StartIncomingNetwork(
	mh *message_handler.MessageHandler,
	cache *cache.Cache,
	config *config.Cfg,
	mainStopChan chan struct{},
) {
	// BaseNode==True signifies that we're not expecting to connect to any
	// remote nodes on connection, so if it's false, we'll skip that step and
	// just wait for incoming connections.
	if !config.BaseNode {
		Heartbeat(
			time.Duration(config.HeartbeatInterval)*time.Millisecond,
			time.Duration(config.HeartbeatLoop)*time.Second,
			cache.PeerList,
		)
	}

	networkRouterStopChan := incomingNetwork.StartNetworkRouter(mh, cache, config)
	// TODO(ian): Clean up this for statement, it's technical debt.
	for {
		select {
		default:
			continue
		case <-mainStopChan:
			networkRouterStopChan <- struct{}{}
			break
		}
	}
}
