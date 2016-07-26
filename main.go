package main

import (
	"github.com/GrappigPanda/Olivia/cache"
	"github.com/GrappigPanda/Olivia/network"
	"github.com/GrappigPanda/Olivia/network/message_handler"
)

func Init() {
	internalCache := cache.NewCache()

	messageHandler := message_handler.NewMessageHandler()
	go networkHandler.StartIncomingNetwork(messageHandler, internalCache)
}

func main() {
	Init()
}
