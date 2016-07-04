package olilib_network

import (
        "net"
        "github.com/GrappigPanda/Olivia/lib/query_language"
        "bufio"
        "log"
)

type ConnectionCtx struct {
        Parser *query_language.Parser
        Cache *query_language.Cache
}

func StartNetworkRouter() {
        listen, err := net.Listen("tcp",  ":5454")
        if err != nil {
                panic(err)
        }
        defer listen.Close()

        _cache := make(map[string]string)
        cache := query_language.Cache{
                &_cache,
        }

        ctx := &ConnectionCtx{
                query_language.NewParser(),
                &cache,
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
                                 ctx.Parser.Parse(line)
                                 break
                }
        }
}
