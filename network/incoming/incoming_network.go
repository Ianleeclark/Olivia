package incomingNetwork

import (
	"bufio"
	"fmt"
	"log"
	"net"

	"github.com/GrappigPanda/Olivia/cache"
	"github.com/GrappigPanda/Olivia/config"
	"github.com/GrappigPanda/Olivia/dht"
	"github.com/GrappigPanda/Olivia/network/message_handler"
	"github.com/GrappigPanda/Olivia/parser"
)

// ConnectionCtx handles maintaining a persistent state per incoming
// connection.
type ConnectionCtx struct {
	Parser      *parser.Parser
	Cache       *cache.Cache
	Bloomfilter *olilib.BloomFilter
	MessageBus  *message_handler.MessageHandler
	PeerList    *dht.PeerList
}

// StartNetworkRouter initializes everything necessary for our incoming network
// router to process and begins our network router.
func StartNetworkRouter(
	mh *message_handler.MessageHandler,
	cache *cache.Cache,
	peerList *dht.PeerList,
	config *config.Cfg,
) chan struct{} {
	stopchan := make(chan struct{})

	// Here we spin up a network router which allows us to start/stop on demand
	// via channels. It's overly indented, this should probably be separated
	// elsewhere.
	go func(stopchan chan struct{}) {
		listen, err := net.Listen("tcp", fmt.Sprintf(":%d", config.ListenPort))
		if err != nil {
			panic(err)
		}
		defer listen.Close()

		bf := olilib.NewByFailRate(1000, 0.01)

		ctx := &ConnectionCtx{
			parser.NewParser(mh),
			cache,
			bf,
			mh,
			peerList,
		}

		log.Println("Starting connection router!")

		for {
			select {
			default:
				conn, err := listen.Accept()
				if err != nil {
					log.Println(err)
					continue
				}
				log.Println("Incoming connection detected from ",
					conn.RemoteAddr().String(),
				)

				go ctx.handleConnection(&conn)
			case <-stopchan:
				log.Printf("Forcefully quitting network router.")
				return
			}
		}
	}(stopchan)

	return stopchan
}

// handleConnection handles handling state of the incoming network FSM,
// verifying passwords, &c.
func (ctx *ConnectionCtx) handleConnection(conn *net.Conn) {
	defer (*conn).Close()
	// TODO(ian): Implement authentication (new issue).
	connProc := NewProcessorFSM(PROCESSING)
	reader := bufio.NewReader(*conn)
	password := "TestBcryptPassword"

	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			log.Printf("Connection %v failed to readline, closing connection.", *conn)
			break
		}

		switch connProc.State {
		case UNAUTHENTICATED:
			connProc.Authenticate(password)
			log.Println(
				"Unauthenticated request from %v",
				(*conn).RemoteAddr().String(),
			)
			break
		case PROCESSING:
			command, err := ctx.Parser.Parse(string(line), conn)
			if err != nil {
				log.Println(err)
			}

			if command.Command != "PING" {
				log.Printf("Received %v from %v", string(line),
					(*conn).RemoteAddr().String(),
				)
			}

			response := ctx.ExecuteCommand(*command)

			if _, ok := command.Args["BLOOMFILTER"]; ok {
				log.Printf("Responding to %v with bloomfilter",
					(*conn).RemoteAddr().String(),
				)
			} else if command.Command != "PING" {
				log.Printf("Responding to %v %v with %v",
					command.Command,
					command.Args,
					response,
				)
			}

			(*conn).Write([]byte(response))
			break
		}
	}
}
