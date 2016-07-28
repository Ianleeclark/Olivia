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
	keyLookup     map[string]int
}

// keyLookup helps our LRU cache find nodes quicker than traversing the
// entire binary heap. Allows o(1) lookup rather than o(n)

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
		index:     0,
		Tree:      make([]*Node, maxSize),
		keyLookup: make(map[string]int),
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
		keyLookup:     make(map[string]int),
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
	h.keyLookup[node.Key] = h.index
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
	delete(h.keyLookup, retVal.Key)

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

// UpdateNodeTimeout allows changing of the keys timeout in the
func (h *Heap) UpdateNodeTimeout(key string) *Node {
	nodeIndex, ok := h.keyLookup[key]
	if !ok {
		return nil
	}

	h.Tree[nodeIndex].timeout = time.Now().UTC()

	if nodeIndex+1 < h.currentSize {
		if h.compareTwoTimes(nodeIndex, nodeIndex+1) {
			h.percolateDown(nodeIndex)
		} else if h.compareTwoTimes(nodeIndex-1, nodeIndex) {
			h.percolateUp(nodeIndex)
		}
	}

	return h.Tree[h.keyLookup[key]]
}

// Get handles retrieving a Node by its key. Not extensively used, but it was a
// nice-to-have.
func (h *Heap) Get(key string) (*Node, bool) {
	if index, ok := h.keyLookup[key]; ok {
		return h.Tree[index], ok
	} else {
		return nil, ok
	}
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
	if fromIndex == len(h.Tree)-1 {
		return
	}

	// trackerNode is the node which we're currently tracking as it percolates
	// down the binary heap.
	trackerNode := h.Tree[fromIndex]
	// preExistingNode is _any_ node which we're not currently tracking.
	preExistingNode := h.Tree[fromIndex+1]
	// NOTE: I'm not using `compareTwoTimes` here because I think it makes it
	// more readable. I know this is an aggregious abuse of intermediary state
	// or some other nonsense, but it makes it easier for me to read.

	// If our tracker node is nil (meaning it was the slot of the recently
	// evited min node) we want to automatically percolate it to the bottom of
	// the binary heap.
	if trackerNode == nil {
		h.swapTwoNodes(fromIndex, fromIndex+1)
		h.percolateDown(fromIndex + 1)
		return
	}

	// Unlikely to ever do anything. But it asserts that the minimum
	if trackerNode.timeout.Nanosecond() > preExistingNode.timeout.Nanosecond() {
		h.swapTwoNodes(fromIndex, fromIndex+1)
		h.percolateDown(fromIndex + 1)
	}
}

// swapTwoNodes swaps j into i and vice versa. Moreover, it handles updating
// the keyLookup field in the heap so that we can continue to quickly retrieve
// key timeouts.
func (h *Heap) swapTwoNodes(i int, j int) {
	// If we find a value at Tree[i], we can update it in the keylookup,
	// otherwise disregard, as it's a recently evicted node.
	if h.Tree[i] != nil {
		h.keyLookup[h.Tree[i].Key] = j
	}
	if h.Tree[j] != nil {
		h.keyLookup[h.Tree[j].Key] = i
	}

	h.Tree[j], h.Tree[i] = h.Tree[i], h.Tree[j]
}

// compareTwoTimes takes two indexes and compares the `.Nanosecond()` value of
// each in the tree. If the left (i) has an expiration time _after_ the right
// (j), then we return True. Otherwise, if the left (i) has an expiration time
// _before_ the right (j) we return a False.
func (h *Heap) compareTwoTimes(i int, j int) bool {
	if h.Tree[i].timeout.Nanosecond() > h.Tree[j].timeout.Nanosecond() {
		return true
	} else {
		return false
	}
}
