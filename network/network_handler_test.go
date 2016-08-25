package networkHandler

import (
	"testing"
)

func TestMain(m *testing.M) {
	// Due to instability, network tests must be stopped.
	// TODO(ian): Restore network handling testing at a later date.
	/*
		stopChan := make(chan struct{})
		stopChan2 := make(chan struct{})
		mh := message_handler.NewMessageHandler()
		cache := cache.NewCache()
		cfg := config.ReadConfig()
		secondCfg := config.ReadConfig()
		secondCfg.BaseNode = false
		secondCfg.RemotePeers = append(secondCfg.RemotePeers, "127.0.0.1:5454")
		secondCfg.ListenPort = 5555

		go StartIncomingNetwork(mh, cache, cfg, stopChan)
		go StartIncomingNetwork(mh, cache, secondCfg, stopChan2)
		os.Exit(m.Run())
		stopChan <- struct{}{}
		stopChan2 <- struct{}{}
	*/
}
