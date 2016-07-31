package network_receiver

import (
	"log"
	"bufio"
	. "github.com/GrappigPanda/Olivia/network/message_handler"
	"net"
	"strings"
)

type IncomingChannel chan string
type RequesterResponseChannel chan string

type Receiver struct {
	ReceiverChannel IncomingChannel
	MessageStore    *MessageHandler
	conn            *net.Conn
}

func NewReceiver(messageStore *MessageHandler, conn *net.Conn) *Receiver {
	return &Receiver{
		make(IncomingChannel),
		messageStore,
		conn,
	}
}

func (r *Receiver) Run() {
	reader := bufio.NewReader(*r.conn)
	for {
		buffer, _, err := reader.ReadLine()
		if err != nil {
			(*r.conn).Write([]byte("Invalid Command"))
			continue
		}

		go r.processIncomingString(string(buffer))
	}
}

func (r *Receiver) processIncomingString(incomingString string) {
	splitString := strings.Split(incomingString, ":")
	if len(splitString) != 2 {
		// TODO(ian): Should we have a reference to the conn object and
		// respond on failures to split?
		log.Print("Invalid received command, no hash found")
		return
	} else {
		log.Println(splitString)
	}

	hash := splitString[0]
	if len(hash) != 32 {
		return
	}

	callbackChan := make(chan chan string)
	(*r.MessageStore).RemoveKeyChannel <- NewKeyValPair(hash, nil, callbackChan)

	requesterChannel := <-callbackChan

	go func() {
		requesterChannel <- splitString[1]
	}()
}

func openListeningConnection() (*net.Listener, error) {
	ln, err := net.Listen("tcp", ":5555")
	if err != nil {
		return nil, err
	}

	return &ln, nil
}
