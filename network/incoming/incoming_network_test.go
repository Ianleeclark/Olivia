package incomingNetwork

import (
	"github.com/GrappigPanda/Olivia/cache"
	"github.com/GrappigPanda/Olivia/config"
	"github.com/GrappigPanda/Olivia/dht"
	"github.com/GrappigPanda/Olivia/network/message_handler"
	"os"
	"testing"
)

func TestStartStopNetworkRouter(t *testing.T) {
	t.Errorf("test")
}

func TestGetBloomfilter(t *testing.T) {

}

func TestSetKey(t *testing.T) {

}

func TestGetKeyFromRemoteNode(t *testing.T) {
	// It's not yet possible to do this, for this to be done, we need to first
	// add listening ports and base nodes to be a part of the config file.
}

func TestMain(m *testing.M) {
	mh := message_handler.NewMessageHandler()
	cache := cache.NewCache()
	peerList := dht.NewPeerList(mh)
	config := config.ReadConfig()
	stopchan := StartNetworkRouter(mh, cache, peerList, config)

	os.Exit(m.Run())

	stopchan <- struct{}{}
}
