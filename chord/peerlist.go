package chord

import (
	"fmt"
	"log"
)

// PeerList is a data structure which represents remote olivia nodes.
type PeerList struct {
	Peers []*Peer
	// TODO(ian): Add a backup list
}

// NewPeerList Creates a new peer list
func NewPeerList() *PeerList {
	peerlist := make([]*Peer, 3)

	return &PeerList{
		peerlist,
	}
}

// ConnectAllPeers connects all peers (or at least attempts to)
func (p *PeerList) ConnectAllPeers() error {
	failureCount := 0
	for x := range p.Peers {
		if p.Peers[x] == nil {
			failureCount++
			continue
		}

		if err := p.Peers[x].Connect(); err != nil {
			log.Println(err)
			failureCount++
		}
	}

	if failureCount == len(p.Peers) {
		log.Println("Failed to connect to any nodes.")
		return fmt.Errorf("No connectable nodes.")
	}

	return nil
}

// DisconnectAllPeers disconnects all peers
func (p *PeerList) DisconnectAllPeers() {
	for x := range p.Peers {
		if err := p.Peers[x].Connect(); err != nil {
			log.Println(err)
		}
	}
}
