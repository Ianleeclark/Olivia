package bfsearch

import (
	"github.com/GrappigPanda/Olivia/dht"
)

type bloomfilterNode struct {
	bitIndex uint
	refs     []*dht.Peer
}

type Search struct {
	nodes []*bloomfilterNode
}

type BloomSearch interface {
	Get(int) []*dht.Peer
	Recalculate()
}

func NewSearch(peerList dht.PeerList) *Search {
	return calculateSearchArray(peerList)
}

func (b *Search) Recalculate(peerList dht.PeerList) {
	b = calculateSearchArray(peerList)
}

func (b *Search) Get(bitIndex uint) []*dht.Peer {
	if bitIndex > uint(len(b.nodes)) {
		return nil
	}

	foundNode := b.nodes[bitIndex]

	if foundNode != nil {
		return foundNode.refs
	}

	return nil
}

func (b *Search) GetFromIndices(bitIndex []uint) []*dht.Peer {
	for _, index := range bitIndex {
		if index > uint(len(b.nodes)) {
			return nil
		}

		foundNodes := b.nodes[index]

		if foundNodes != nil {
			return foundNodes.refs
		}
	}

	return nil
}

func unionPeerLists(peerLists ...[]*dht.Peer) []*dht.Peer {
	peerListRefCounter := make(map[string]int)
	peerListAllPeers := make(map[string]*dht.Peer)

	for _, peerList := range peerLists {
		for _, peer := range peerList {
			// NOTE: fill up our ref counter
			if _, ok := peerListRefCounter[peer.UniqueID]; !ok {
				peerListRefCounter[peer.UniqueID] = 0
			} else {
				peerListRefCounter[peer.UniqueID]++
			}

			// NOTE: Fill up our all peer list
			if _, ok := peerListAllPeers[peer.UniqueID]; !ok {
				peerListAllPeers[peer.UniqueID] = peer
			}
		}
	}

	// NOTE: If any peer is referenced len(peerLists) times, it means they are in all of our lists.
	var unionedPeers []*dht.Peer
	for peerId, refCount := range peerListRefCounter {
		if refCount == len(peerLists) {
			unionedPeers = append(unionedPeers, peerListAllPeers[peerId])
		}
	}

	return unionedPeers
}

func calculateSearchArray(peerList dht.PeerList) *Search {
	var bfNodes []*bloomfilterNode

	if peerList.Peers[0] == nil || len(peerList.Peers) == 0 {
		return &Search{
			nodes: bfNodes,
		}
	}

	peerBF := peerList.Peers[0].BloomFilter
	bfSize := uint(0)
	if peerBF != nil {
		bfSize = peerBF.GetMaxSize()
	}

	for i := uint(0); i <= bfSize; i++ {
		var nodes []*dht.Peer

		for _, peer := range peerList.Peers {
			if peer != nil {
				bf := peer.BloomFilter
				bitset := bf.GetStorage()

				if bitset.IsSet(i) {
					nodes = append(nodes, peer)
				}
			}
		}

		bfNodes = append(
			bfNodes,
			&bloomfilterNode{
				i,
				nodes,
			},
		)
	}

	return &Search{
		nodes: bfNodes,
	}
}
