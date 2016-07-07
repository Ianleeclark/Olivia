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

        go StartNetworkRouter()
}

func TestSendGetCommand(t *testing.T) {
        testChannel := make(chan string)

        fmt.Println("laksdjfklj")
        go func(ch chan string) {
                conn, err := net.DialTimeout("tcp", "localhost:5454", time.Second * 1)
                if err != nil {
                        ch <- "Failed to connect to network"
                }
                fmt.Println("test")

                writer := bufio.NewWriter(conn)
                reader := bufio.NewReader(conn)

                writer.WriteString("GET key1\n")
                val, err := reader.ReadString('\n')
                if err != nil {
                        ch <- "FAILURE"
                }

                ch <- string(val)
        }(testChannel)

        fmt.Println(<-testChannel)
}
