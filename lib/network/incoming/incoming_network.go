package incomingNetwork

import (
	"bufio"
	"github.com/GrappigPanda/Olivia/cache"
	"github.com/GrappigPanda/Olivia/lib/bloomfilter"
	"github.com/GrappigPanda/Olivia/lib/network/message_handler"
	"github.com/GrappigPanda/Olivia/lib/queryLanguage"
	"log"
	"net"
	"github.com/GrappigPanda/Olivia/lib/chord"
)

// ConnectionCtx handles maintaining a persistent state per incoming
// connection.
type ConnectionCtx struct {
	Parser      *queryLanguage.Parser
	Cache       *cache.Cache
	Bloomfilter *olilib.BloomFilter
	MessageBus  *message_handler.MessageHandler
	PeerList *chord.PeerList
}

// StartNetworkRouter initializes everything necessary for our incoming network
// router to process and begins our network router.
func StartNetworkRouter(
	mh *message_handler.MessageHandler,
	cache *cache.Cache,
	peerList *chord.PeerList,
) {

	listen, err := net.Listen("tcp", ":5454")
	if err != nil {
		panic(err)
	}
	defer listen.Close()

	bf := olilib.NewByFailRate(10000, 0.01)

	ctx := &ConnectionCtx{
		queryLanguage.NewParser(mh),
		cache,
		bf,
		mh,
		peerList,
	}

	log.Println("Starting connection router.")

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go ctx.handleConnection(&conn)
	}
}

// handleConnection handles handling state of the incoming network FSM,
// verifying passwords, &c.
func (ctx *ConnectionCtx) handleConnection(conn *net.Conn) {
	defer (*conn).Close()
	// TODO(ian): Implement authentication (new issue).
	connProc := NewProcessorFSM(PROCESSING)
	reader := bufio.NewReader(*conn)

	for {
		// TODO(ian): Replace this with a new language processor for incoming
		// commands
		password := "TestBcryptPassword"
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Connection %v failed to readline, closing connection.", *conn)
			break
		}

		switch connProc.State {
		case UNAUTHENTICATED:
			connProc.Authenticate(password)
			break
		case PROCESSING:
			command, err := ctx.Parser.Parse(line, conn)
			if err != nil {
				log.Println(err)
			}
			response := ctx.ExecuteCommand(*command)
			(*conn).Write([]byte(response))
			break
		}
	}
}
