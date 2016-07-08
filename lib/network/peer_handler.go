package olilib_network

import (
        "time"
        "net"
        "github.com/GrappigPanda/Olivia/lib/bloomfilter"
        "log"
)

type ConnectionStatus int

const (
        DEAD ConnectionStatus = iota
        ALIVE
)

type Peer struct {
        Address net.Addr
        Conn *net.Conn
        Bloomfilter *olilib.BloomFilter
        LastUpdate time.Time
        Status ConnectionStatus
}

// NewPeer Obviously handles constructing a new Peer object
func NewPeer(Address net.Addr) (*Peer, error) {
        conn, err := net.DialTimeout("tcp", Address.String(), time.Second * 5)
        if err != nil {
                log.Printf("")
                return nil, err
        }

        return &Peer{
                Address,
                &conn,
                getBloomFilter(&conn),
                time.Now().UTC(),
                ALIVE,
        }, nil
}

// UpdatePeer handles refreshing a connection (if a connection is alive),
// sending a copy of our bloom filter to the remote peer, and receiving a copy
// of the remote host's bloom filter.
func (p *Peer) UpdatePeer() {

}
