package cache

import (
	"github.com/GrappigPanda/Olivia/dht"
)

type bloomfilterNode struct {
	bitIndex uint
	refs     []*dht.Peer
}

type BloomfilterSearch struct {
	nodes []*bloomfilterNode
}

type BloomSearch interface {
	Get(int) []*dht.Peer
	Recalculate()
}

func NewBloomfilterSearch(peerList dht.PeerList) *BloomfilterSearch {
	return calculateSearchArray(peerList)
}

func (b *BloomfilterSearch) Recalculate(peerList dht.PeerList) {
	b = calculateSearchArray(peerList)
}

func (b *BloomfilterSearch) Get(bitIndex uint) []*dht.Peer {
	if bitIndex > uint(len(b.nodes)) {
		return nil
	}

	foundNode := b.nodes[bitIndex]

	if foundNode != nil {
		return foundNode.refs
	}

	return nil
}

func calculateSearchArray(peerList dht.PeerList) *BloomfilterSearch {
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

	return &BloomfilterSearch{
		nodes: bfNodes,
	}
}
