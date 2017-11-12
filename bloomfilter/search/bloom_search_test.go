package bfsearch

import (
	"github.com/GrappigPanda/Olivia/config"
	"github.com/GrappigPanda/Olivia/dht"
	"testing"
)

var CONFIG = config.ReadConfig()
var PEERLIST = dht.NewPeerList(nil, *CONFIG)

func TestNewBFSearch(t *testing.T) {
	bs := NewSearch(*PEERLIST)

	if bs == nil {
		t.Fatalf("Expected non-nil bloomsearch, but got nil")
	}
}

func TestGetIndex(t *testing.T) {
	bs := NewSearch(*PEERLIST)

	bs.setIndex(5, createFakePeers(2))

	retval := bs.Get(5)
	if len(retval) != 2 {
		t.Fatalf("Expected 2 peers found in BF, got %v", len(retval))
	}
}

func TestGetIndexTooLarge(t *testing.T) {
	bs := NewSearch(*PEERLIST)

	retval := bs.Get(10000000000000000)
	if retval != nil {
		t.Fatalf("Expected nil value from too large index, got %v", retval)
	}
}

func (b *Search) fillIndexWithPeers(i int) {
	peers := createFakePeers(i)

	newNode := &bloomfilterNode{
		bitIndex: uint(i),
		refs:     peers,
	}

	if cap(b.nodes) == 0 {
		b.nodes = make([]*bloomfilterNode, 10)
	}

	b.nodes[i] = newNode
}

func (b *Search) setIndex(index uint, fakePeers []*dht.Peer) {
	if uint(cap(b.nodes)) < index {
		b.nodes = make([]*bloomfilterNode, index+5)
	}

	b.nodes[index] = &bloomfilterNode{
		bitIndex: index,
		refs:     fakePeers,
	}
}

func createFakePeers(i int) []*dht.Peer {
	peers := make([]*dht.Peer, 2)
	for i := 0; i < 2; i++ {
		newPeer := &dht.Peer{
			Status:      dht.Disconnected,
			Conn:        nil,
			IPPort:      "",
			BloomFilter: nil,
			MessageBus:  nil,
		}

		peers[i] = newPeer
	}

	return peers
}
