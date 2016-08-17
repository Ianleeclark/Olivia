package olilib

import (
	"github.com/willf/bitset"
)

// Bitset is a simple wrapper around the willf bitset library.
type Bitset struct {
	BS *bitset.BitSet
}

// NewBitset constructs a new bitset to be used with bloom filters.
func NewBitset(maxSize uint) *Bitset {
	return &Bitset{
		bitset.New(maxSize),
	}
}

// Add handles adding a new hashed index into the bitset.
func (b *Bitset) Add(index uint) {
	b.BS.Set(index)
}

// Contains verifies if a hash index is actually in the bitset or not.
func (b *Bitset) Contains(index uint) bool {
	return b.BS.Test(index)
}

// ToString handles converting the bitset to a RLE usable string.
func (b *Bitset) ToString() string {
	json, err := b.BS.MarshalJSON()
	if err != nil {
		panic(err)
	}

	return string(json[1 : len(string(json))-2])
}

// FromString handles converting a (valid json) string to a valid underlying
// bitset.
func (b *Bitset) FromString(inputString string) {
	err := b.BS.UnmarshalJSON([]byte(inputString))
	if err != nil {
		panic(err)
	}
}
