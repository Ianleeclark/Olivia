package dht

import (
	"fmt"
	"github.com/GrappigPanda/Olivia/network/message_handler"
	"log"
	"strings"
	"sync"
)

// PeerList is a data structure which represents remote olivia nodes.
type PeerList struct {
	Peers       []*Peer
	BackupPeers []*Peer
	PeerMap     *map[string]bool
	MessageBus  *message_handler.MessageHandler
	sync.Mutex
}

// NewPeerList Creates a new peer list
func NewPeerList(mh *message_handler.MessageHandler) *PeerList {
	peerlist := make([]*Peer, 3)
	// We originally allocate 10 slots for backup peers, but if necessary
	// we readjust whenever we request peers from a new node.
	backupList := make([]*Peer, 10)

	peerMap := make(map[string]bool)

	return &PeerList{
		Peers:       peerlist,
		BackupPeers: backupList,
		PeerMap:     &peerMap,
		MessageBus:  mh,
	}
}

// AddPeer handles intelligently putting a peer into our peer list. Priority
// of insertion is towards Peers first and then BackupPeers.
func (p *PeerList) AddPeer(ipPort string) {
	if _, ok := (*p.PeerMap)[ipPort]; ok {
		// If we already have the peer stored, we don't need to
		// add it again.
		return
	}

	newPeer := NewPeerByIP(ipPort, p.MessageBus)

	p.Lock()
	defer p.Unlock()
	if len(p.Peers)+1 <= 3 {
		p.Peers = append(p.Peers, newPeer)
		return
	}

	if len(p.BackupPeers)+1 >= cap(p.BackupPeers) {
		p.BackupPeers = append(
			p.BackupPeers,
			make([]*Peer, cap(p.BackupPeers)*2)...,
		)

		p.BackupPeers = append(p.BackupPeers, newPeer)
	}

	newPeer.Connect()

	return
}

// ConnectAllPeers connects all peers (or at least attempts to)
func (p *PeerList) ConnectAllPeers() error {
	responseChannel := make(chan string)
	go p.handlePeerQueries(responseChannel)

	failureCount := 0
	successCount := 0

	p.Lock()
	defer p.Unlock()
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

		log.Println(
			"Connected to ",
			p.Peers[x].IPPort,
			"Requesting peer list",
		)

		log.Println("Sending Request Connect")
		p.Peers[x].SendCommand("0:REQUEST CONNECT\n")
		p.Peers[x].GetPeerList(responseChannel)
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

// handlePeerQueries handles the responses for each peer list.
func (p *PeerList) handlePeerQueries(responseChannel chan string) {
	p.Lock()
	defer p.Unlock()
	for response := range responseChannel {
		splitResponse := strings.SplitN(response, " ", 2)
		if len(splitResponse) != 2 {
			continue
		}

		peers := strings.Split(splitResponse[1], ",")

		for i := range peers {
			p.AddPeer(peers[i])
		}
	}

	return
}
