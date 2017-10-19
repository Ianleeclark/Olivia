package cache

import (
	"time"
)

// executeRepeatedly Allows repeated calls to any function which doesn't accept
// arguments. Allows for remote stopping of the execution and passing back
// total number of executions.
func (c *Cache) executeRepeatedly(
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
func (c *Cache) heartbeatRemoteNodes(interval time.Duration) {
	c.executeRepeatedly(
		interval,
		func() {
			if c.PeerList != nil {
				for _, peer := range c.PeerList.Peers {
					if peer != nil {
						go peer.TestConnection()
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
func (c *Cache) getRemoteBloomFilters(interval time.Duration) {
	c.executeRepeatedly(
		interval,
		func() {
			if c.PeerList != nil {
				for _, peer := range c.PeerList.Peers {
					if peer != nil {
						go peer.GetBloomFilter()
					}
				}
			}

			c.bloomfilterSearch.Recalculate(*c.PeerList)
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
func (c *Cache) Heartbeat() {
	go c.heartbeatRemoteNodes(time.Duration(200) * time.Millisecond)
	go c.getRemoteBloomFilters(time.Duration(30) * time.Second)
}
