package olilib_network

import (
        "net"
        "fmt"
)

func StartNetworkRouter() {
        listen, err := net.Listen("tcp",  ":5454")
        if err != nil {
                panic(err)
        }
        defer listen.Close()

        for {
                conn, err := listen.Accept()
                if err != nil {
                        fmt.Println(err)
                        continue
                }

                go handleConnection(&conn)
        }
}

func handleConnection(conn *net.Conn) {
       defer (*conn).Close()
       // TODO(ian): Implement authentication (new issue).
       conn_proc := NewProcessorFSM(PROCESSING)

       for {
               // TODO(ian): Replace this with a new language processor for incoming
               // commands
               password := "TestBcryptPassword"


               switch conn_proc.State {
                        case UNAUTHENTICATED:
                                conn_proc.Authenticate(password)
                                break
                        case PROCESSING:
                                break
               }
       }

}
