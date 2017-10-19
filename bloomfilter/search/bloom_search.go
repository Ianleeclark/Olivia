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

func calculateSearchArray(peerList dht.PeerList) *Search {
	bfSize := peerList.Peers[0].BloomFilter.GetMaxSize()

	var bfNodes []*bloomfilterNode

	for i := uint(0); i <= bfSize; i++ {
		var nodes []*dht.Peer

		for _, peer := range peerList.Peers {
			bf := peer.BloomFilter
			bitset := bf.GetStorage()

			if bitset.IsSet(i) {
				nodes = append(nodes, peer)
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
