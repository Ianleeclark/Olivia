package dht

import (
	"fmt"
	"log"
)

// PeerList is a data structure which represents remote olivia nodes.
type PeerList struct {
	Peers []*Peer
	BackupPeers []*Peer
}

// NewPeerList Creates a new peer list
func NewPeerList() *PeerList {
	peerlist := make([]*Peer, 3)
	// We originally allocate 10 slots for backup peers, but if necessary
	// we readjust whenever we request peers from a new node.
	backupList := make([]*Peer, 10)

	return &PeerList{
		peerlist,
		backupList,
	}
}

// ConnectAllPeers connects all peers (or at least attempts to)
func (p *PeerList) ConnectAllPeers() error {
	failureCount := 0
	successCount := 0
	for x := range p.Peers {
		if p.Peers[x] == nil {
			failureCount++
			continue
		}
		log.Println("Attempting connection to ", p.Peers[x].IPPort)

		if err := p.Peers[x].Connect(); err != nil {
			log.Println(err)
			failureCount++
			continue
		}

		successCount++

		log.Println("Connected to ", p.Peers[x].IPPort)
	}

	if failureCount == len(p.Peers) {
		log.Println("Failed to connect to any nodes.")
		return fmt.Errorf("No connectable nodes.")
	}

	log.Println("Connected to ", successCount, " nodes.")
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
