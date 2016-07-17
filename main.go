package main

import (
	"github.com/GrappigPanda/Olivia/cache"
	"github.com/GrappigPanda/Olivia/lib/network"
	"github.com/GrappigPanda/Olivia/lib/network/message_handler"
)

func Init() {
	internalCache := cache.NewCache()

	messageHandler := message_handler.NewMessageHandler()
	networkHandler.StartIncomingNetwork(messageHandler, internalCache)
}

func main() {
	Init()
}
