package main

import (
	"github.com/GrappigPanda/Olivia/cache"
	"github.com/GrappigPanda/Olivia/config"
	"github.com/GrappigPanda/Olivia/network"
	"github.com/GrappigPanda/Olivia/network/message_handler"
)

func Init() {
	config := config.ReadConfig()

	messageHandler := message_handler.NewMessageHandler()

	internalCache := cache.NewCache(messageHandler, config)

	networkHandler.StartIncomingNetwork(
		messageHandler,
		internalCache,
		config,
		nil,
	)
}

func main() {
	Init()
}
