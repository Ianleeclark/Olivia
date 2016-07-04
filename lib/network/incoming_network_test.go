package olilib_network

import (
        "time"
        "fmt"
        "testing"
        "net"
        "bufio"
)

func TestRouteNewConection(t *testing.T) {
        go func() {
                _, err := net.DialTimeout("tcp", "localhost:5454", time.Second * 1)
                if err != nil {
                        t.Fatalf("Failed to connect to network")
                }
        }()
        fmt.Println("test")

        go StartNetworkRouter()
}

func TestSendGetCommand(t *testing.T) {
        go func() {
                conn, err := net.DialTimeout("tcp", "localhost:5454", time.Second * 1)
                if err != nil {
                        t.Fatalf("Failed to connect to network")
                }

                writer := bufio.NewWriter(conn)
                reader := bufio.NewReader(conn)

                writer.WriteString("GET key1\n")
                fmt.Println(reader.ReadBytes('\n'))
        }()
}
