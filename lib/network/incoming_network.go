package olilib_network

import (
	"bufio"
	"github.com/GrappigPanda/Olivia/lib/bloomfilter"
	"github.com/GrappigPanda/Olivia/lib/query_language"
	"log"
	"net"
)

type ConnectionCtx struct {
	Parser      *query_language.Parser
	Cache       *Cache
	Bloomfilter *olilib.BloomFilter
}

func StartNetworkRouter() {
	listen, err := net.Listen("tcp", ":5454")
	if err != nil {
		panic(err)
	}
	defer listen.Close()

	_cache := make(map[string]string)
	cache := Cache{
		&_cache,
	}

	bf := olilib.NewByFailRate(10000, 0.01)

	ctx := &ConnectionCtx{
		query_language.NewParser(),
		&cache,
		bf,
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

func (ctx *ConnectionCtx) handleConnection(conn *net.Conn) {
	defer (*conn).Close()
	// TODO(ian): Implement authentication (new issue).
	conn_proc := NewProcessorFSM(PROCESSING)
	reader := bufio.NewReader(*conn)

	for {
		// TODO(ian): Replace this with a new language processor for incoming
		// commands
		password := "TestBcryptPassword"
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Connection %v failed to readline, closing connection.", *conn)
			break
		}

		switch conn_proc.State {
		case UNAUTHENTICATED:
			conn_proc.Authenticate(password)
			break
		case PROCESSING:
			command, err := ctx.Parser.Parse(line)
			if err != nil {

			}
			response := ctx.ExecuteCommand(command.Command, command.Args)
			(*conn).Write([]byte(response))
			break
		}
	}
}
