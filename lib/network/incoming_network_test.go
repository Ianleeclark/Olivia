package olilib_network

import (
        "testing"
        "net"
)

func TestRouteNewConection(t *testing.T) {
        go StartNetworkRouter()

        _, err := net.Dial("tcp", "localhost:5454")
        if err != nil {
                t.Fatalf("Failed to connect to network")
        }

        _, err = net.Dial("tcp", "localhost:5454")
        if err != nil {
                t.Fatalf("Failed to connect to network")
        }
}
