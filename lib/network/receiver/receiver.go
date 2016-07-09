package network_receiver

import (
        "strings"
        "bufio"
        "net"
        "log"
)

type IncomingChannel chan string
type RequesterResponseChannel chan string

type Receiver struct {
        ReceiverChannel IncomingChannel
        MessageStore *map[string]RequesterResponseChannel
        Listener *net.Listener
}

func NewReceiver(messageStore *map[string]RequesterResponseChannel) *Receiver {
        ln, err := openListeningConnection()
        if err != nil {
                // TODO(ian): Handle this at a supervisor level.
                panic(err)
        }

        return &Receiver{
                make(IncomingChannel),
                messageStore,
                ln,
        }
}

func (r *Receiver) Run() {
        for {
                conn, err := (*r.Listener).Accept()
                if err != nil {
                        log.Printf("")
                }

                go r.handleConnection(&conn)
        }
}

func (r *Receiver) handleConnection(conn *net.Conn) {
        reader := bufio.NewReader(*conn)
        for {
                buffer, err := reader.ReadString('\n')
                if err != nil {
                        (*conn).Write([]byte("Invalid Command"))
                        continue
                }

                go r.processIncomingString(buffer)
        }
}

func (r* Receiver) processIncomingString(incomingString string) {
        splitString := strings.Split(incomingString, ":")
        if len(splitString) != 2 {
                // TODO(ian): Should we have a reference to the conn object and
                // respond on failures to split?
                return
        }

        hash := splitString[0]
        if len(hash) != 32 {
                return
        }

        requesterChannel, hashExists := (*r.MessageStore)[hash]
        if !hashExists {
                return
        }

        requesterChannel <- splitString[1]
}

func openListeningConnection() (*net.Listener, error) {
        ln, err := net.Listen("tcp", ":5555")
        if err != nil {
                return nil, err
        }

        return &ln, nil
}
