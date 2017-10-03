package cache

import (
	"github.com/GrappigPanda/Olivia/bloomfilter"
	"github.com/GrappigPanda/Olivia/dht"
)

type bloomfilterNode struct {
	bitIndex int
	refs     []*dht.Peer
}

type BloomfilterSearch struct {
	nodes []*bloomfilterNode
}

type BloomSearch interface {
	Get(int) []*dht.Peer
	Recalculate()
}

func NewBloomfilterSearch(bf *bloomfilter.BloomFilter) *BloomfilterSearch {
	return nil
}

func (b *BloomfilterSearch) Get(bitIndex int) []*dht.Peer {
	return nil
}
