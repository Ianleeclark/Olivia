package binheap

import (
	"time"
)

// Node represents each LRUStorage node.
type Node struct {
	Key     string
	Timeout time.Time
}

// NewNode Allocates a new Node. It is not placed in the binary heap at
// allocation. Rather, the caller is expected to later Insert the newly created
// node into the binary heap.
func NewNode(key string, Timeout time.Time) *Node {
	return &Node{
		Key:     key,
		Timeout: Timeout,
	}
}
