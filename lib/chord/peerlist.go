package chord

import (
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
func (p *PeerList) ConnectAllPeers() {
	for x := range p.Peers {
		if err := p.Peers[x].Connect(); err != nil {
			log.Println(err)
		}
	}
}

// DisconnectAllPeers disconnects all peers
func (p *PeerList) DisconnectAllPeers() {
	for x := range p.Peers {
		if err := p.Peers[x].Connect(); err != nil {
			log.Println(err)
		}
	}
}
