package bloomfilter

import (
	"github.com/willf/bitset"
	"log"
)

type Bitset interface {
	Add(uint)
	Contains(uint) bool
	ToString() string
	FromString(string)
	Compare(interface{}) bool
	IsSet(uint) bool
	Len() uint
}

// WFBitset is a simple wrapper around the willf bitset library.
type WFBitset struct {
	bs *bitset.BitSet
}

// NewWFBitset constructs a new bitset to be used with bloom filters.
func NewWFBitset(maxSize uint) *WFBitset {
	return &WFBitset{
		bitset.New(maxSize),
	}
}

// Add handles adding a new hashed index into the bitset.
func (b *WFBitset) Add(index uint) {
	b.bs.Set(index)
}

// Contains verifies if a hash index is actually in the bitset or not.
func (b *WFBitset) Contains(index uint) bool {
	return b.bs.Test(index)
}

// ToString handles converting the bitset to a RLE usable string.
func (b *WFBitset) ToString() string {
	json, err := b.bs.MarshalJSON()
	if err != nil {
		panic(err)
	}

	return string(json[1 : len(string(json))-2])
}

// FromString handles converting a (valid json) string to a valid underlying
// bitset.
func (b *WFBitset) FromString(inputString string) {
	err := b.bs.UnmarshalJSON([]byte(inputString))
	if err != nil {
		log.Println("Invalid bloomfilter received: ", err)
	}
}

func (b *WFBitset) Compare(compareTo interface{}) bool {
	return b.bs.Equal(compareTo.(*WFBitset).bs)
}

func (b *WFBitset) IsSet(bitIndex uint) bool {
	return b.bs.Test(bitIndex)
}

func (b *WFBitset) Len() uint {
	return b.bs.Len()
}
