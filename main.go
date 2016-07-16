package main

import (
	"github.com/GrappigPanda/Olivia/lib/network"
	"github.com/GrappigPanda/Olivia/lib/network/message_handler"
)

func Init() {
	messageHandler := message_handler.NewMessageHandler()
	go networkHandler.StartIncomingNetwork(messageHandler)
}

func main() {
}
