package peer

import (
	"net"
	"time"
)

// State represents the state that the remote peer is in.
type State int

const (
	// Disconnected signifies that the remote node is not yet connected.
	Disconnected State = iota
	// Connected signifies that there is a working connection between the
	// remote peer and our current node.
	Connected
	// Timeout signifies that the remote node has timed out and a
	// connection couldn't be established.
	Timeout
)

// Peer Houses the state for remote Peers
type Peer struct {
	Status State
	Conn   *net.Conn
	ipPort string
}

// Connect opens a connection to a remote peer
func (p *Peer) Connect() {
	conn, err := net.DialTimeout("tcp", p.ipPort, 5*time.Second)
	if err != nil {
		if err, ok := err.(net.Error); ok && err.Timeout() {
			p.Status = Timeout
		}
		return
	}

	p.Conn = &conn
}

// Disconnect closes a connection to a remote peer.
func (p *Peer) Disconnect() {
	(*p.Conn).Close()
}

// RefreshConnection closes and reopens a connection to a remote peer. Surprise
// Surprise.
func (p *Peer) RefreshConnection() {
	p.Disconnect()
	p.Connect()
}
