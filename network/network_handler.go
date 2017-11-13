package networkHandler

import (
	"github.com/GrappigPanda/Olivia/cache"
	"github.com/GrappigPanda/Olivia/config"
	"github.com/GrappigPanda/Olivia/network/incoming"
	"github.com/GrappigPanda/Olivia/network/message_handler"
)

// doing any error checking (in the future) as well as ensuring that it
// continues running.
func StartIncomingNetwork(
	mh *message_handler.MessageHandler,
	cache *cache.Cache,
	config *config.Cfg,
	mainStopChan chan struct{},
) {
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
