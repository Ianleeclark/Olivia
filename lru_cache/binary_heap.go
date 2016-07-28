package olilib_lru

import (
	"time"
)

// HeapAllocationStrategy is a type definition used for enum values in handling
// heap reallocation.
type HeapAllocationStrategy int

const (
	Maintain HeapAllocationStrategy = iota
	Realloc
)

// Node represents each binary heap node.
type Node struct {
	Key     string
	timeout time.Time
}

// Heap represents our binary heap object.
type Heap struct {
	Tree          []*Node
	currentSize   int
	index         int
	allocStrategy HeapAllocationStrategy
}

// NewNode Allocates a new Node. It is not placed in the binary heap at
// allocation. Rather, the caller is expected to later Insert the newly created
// node into the binary heap.
func NewNode(key string, timeout time.Time) *Node {
	return &Node{
		Key:     key,
		timeout: timeout,
	}
}

// NewHeap handles initializing a new Heap object. The default strategy for
// reallocation is to `Maintain` which means that whenever there is an attempt
// to insert a new node, the min node (the root index) will be evicted from the
// heap.
func NewHeap(maxSize int) *Heap {
	return &Heap{
		index: 0,
		Tree:  make([]*Node, maxSize),
	}
}

// NewHeapReallocate handles intiailizing a new heap object which is able to
// reallocate itself.
func NewHeapReallocate(maxSize int) *Heap {
	return &Heap{
		index:         0,
		Tree:          make([]*Node, maxSize),
		currentSize:   0,
		allocStrategy: Realloc,
	}
}

// MinNode returns the root node. In this implementation, we opted for a
// minimum binary heap instead of a generic implementation.
func (h *Heap) MinNode() *Node {
	if h.currentSize == 0 {
		return nil
	}

	return h.Tree[0]
}

// Insert handles placing a new node into the heap. If the allocation strategy
// is set to `Maintain`, then and only then will `Insert` return a *Node.
// Moreover, a *Node is only returned if the binary heap is full and can no
// longer place new nodes into it.
func (h *Heap) Insert(node *Node) *Node {
	if h.index+0 >= len(h.Tree) {
		// If we run into the bounds of our heap, we need to either
		// reallocate (if that's what we're wanting to do, or
		// maintain the size and
		if h.allocStrategy == Realloc {
			// The default behavior is to expand the heap by
			// 0.5 times.
			h.ReAllocate(h.index + len(h.Tree)/2)
		} else {
			// Otherwise, if we're maintaining, we want to evict
			// the root node (The Min Node).
			return h.EvictMinNode()
		}
	}

	h.Tree[h.index] = node
	// It's unlikely that percolating up is ever necessary, as we don't
	// typically insert nodes with an expiration time sooner than nodes already
	// living in the binary heap, but it's important to have, regardless.
	h.percolateUp(h.index)

	h.index++
	h.currentSize++

	return nil
}

// EvictMinNode takes the root node (index 0) and removes it from the binary
// heap. It then reorganizes the binary heap so that everything stays in order
// correctly.
func (h *Heap) EvictMinNode() *Node {
	if h.index == 0 {
		return nil
	}

	retVal := h.Tree[0]
	h.Tree[0] = nil

	// Decrement the current size specifically before percolating down, as
	// percolating down needs bounds checking against the end of the binary
	// heap.
	h.currentSize--
	h.index--
	h.percolateDown(0)

	return retVal
}

// IsEmpty Notifies the caller if the binary heap is empty.
func (h *Heap) IsEmpty() bool {
	return h.currentSize == 0
}

// ReAllocate Handles increasing the size of the underlying binary heap.
func (h *Heap) ReAllocate(maxSize int) {
	h.Tree = append(h.Tree, make([]*Node, maxSize)...)
}

// percolateUp handles sorting a newly inserted node into its correct position.
// It's very unlikely this function actually ever does anything, as it's only
// called by `Insert`, so newly inserted nodes don't typically have an
// expiration time sooner than nodes already living in the heap.
func (h *Heap) percolateUp(newNodeIndex int) {
	if newNodeIndex == 0 {
		return
	}

	newlyInsertedNode := h.Tree[newNodeIndex]
	preExistingNode := h.Tree[newNodeIndex-1]

	// Unlikely to ever do anything.
	if newlyInsertedNode.timeout.Nanosecond() < preExistingNode.timeout.Nanosecond() {
		h.swapTwoNodes(newNodeIndex, newNodeIndex-1)
		h.percolateUp(newNodeIndex - 1)
	}
}

// percolateDown handles moving a node starting at index `fromIndex` down into
// its correct spot in the binary heap.
func (h *Heap) percolateDown(fromIndex int) {
	if fromIndex == h.currentSize {
		return
	}

	// trackerNode is the node which we're currently tracking as it percolates
	// down the binary heap.
	trackerNode := h.Tree[fromIndex]
	// preExistingNode is _any_ node which we're not currently tracking.
	preExistingNode := h.Tree[fromIndex+1]

	// If our tracker node is nil (meaning it was the slot of the recently
	// evited min node) we want to automatically percolate it to the bottom of
	// the binary heap.
	if trackerNode == nil {
		h.swapTwoNodes(fromIndex, fromIndex+1)
		h.percolateDown(fromIndex + 1)
		return
	}

	// Unlikely to ever do anything.
	if trackerNode.timeout.Nanosecond() > preExistingNode.timeout.Nanosecond() {
		h.swapTwoNodes(fromIndex, fromIndex+1)
		h.percolateDown(fromIndex + 1)
	}
}

// swapTwoNodes swaps j into i and vice versa
func (h *Heap) swapTwoNodes(i int, j int) {
	h.Tree[j], h.Tree[i] = h.Tree[i], h.Tree[j]
}
