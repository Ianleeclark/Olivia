package binheapv1

import (
	"fmt"
	"github.com/GrappigPanda/Olivia/binheap"
	"sync"
	"time"
)

type Heap struct {
	Tree          []*binheap.Node
	currentSize   int
	index         int
	allocStrategy binheap.HeapAllocationStrategy
	keyLookup     map[string]int
	sync.Mutex
}

// keyLookup helps our LRU cache find nodes quicker than traversing the
// entire binary heap. Allows o(1) lookup rather than o(n)

// NewHeap handles initializing a new Heap object. The default strategy for
// reallocation is to `Maintain` which means that whenever there is an attempt
// to insert a new node, the min node (the root index) will be evicted from the
// heap.
func NewHeap(maxSize int) *Heap {
	return &Heap{
		index:     0,
		Tree:      make([]*binheap.Node, maxSize),
		keyLookup: make(map[string]int),
	}
}

// NewHeapReallocate handles intiailizing a new heap object which is able to
// reallocate itself.
func NewHeapReallocate(maxSize int) *Heap {
	return &Heap{
		index:         0,
		Tree:          make([]*binheap.Node, maxSize),
		currentSize:   0,
		allocStrategy: binheap.Realloc,
		keyLookup:     make(map[string]int),
	}
}

// Copy handles taking in a binary heap and making a copy of it.
func (h *Heap) Copy() Heap {
	h.Lock()
	defer h.Unlock()
	newHeap := NewHeap(len(h.Tree))

	for index, element := range h.Tree {
		newHeap.Tree[index] = element
	}

	for k, v := range h.keyLookup {
		newHeap.keyLookup[k] = v
	}

	newHeap.index = h.index
	newHeap.currentSize = h.currentSize

	return *newHeap
}

// Minbinheap.Node returns the root node. In this implementation, we opted for a
// minimum binary heap instead of a generic implementation.
func (h *Heap) MinNode() *binheap.Node {
	if h.currentSize == 0 {
		return nil
	}

	return h.Tree[0]
}

// Insert handles placing a new node into the heap. If the allocation strategy
// is set to `Maintain`, then and only then will `Insert` return a *binheap.Node.
// Moreover, a *binheap.Node is only returned if the binary heap is full and can no
// longer place new nodes into it.
func (h *Heap) Insert(node *binheap.Node) *binheap.Node {
	if h.index >= len(h.Tree) {
		// If we run into the bounds of our heap, we need to either
		// reallocate (if that's what we're wanting to do, or
		// maintain the size and
		if h.allocStrategy == binheap.Realloc {
			// The default behavior is to expand the heap by
			// 0.5 times.
			h.ReAllocate(h.index + len(h.Tree)/2)
		} else {
			// Otherwise, if we're maintaining, we want to evict
			// the root node (The Min binheap.Node).
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

// EvictMinbinheap.Node takes the root node (index 0) and removes it from the binary
// heap. It then reorganizes the binary heap so that everything stays in order
// correctly.
func (h *Heap) EvictMinNode() *binheap.Node {
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
	delete(h.keyLookup, retVal.Key)
	h.percolateDown(0)

	return retVal
}

// Peek handles looking at the index of the tree.
func (h *Heap) Peek(index int) (*binheap.Node, error) {
	if index >= h.currentSize {
		return nil, fmt.Errorf("Index greater than size of heap.")
	}

	return h.Tree[index], nil
}

// IsEmpty Notifies the caller if the binary heap is empty.
func (h *Heap) IsEmpty() bool {
	return h.currentSize == 0
}

// ReAllocate Handles increasing the size of the underlying binary heap.
func (h *Heap) ReAllocate(maxSize int) {
	h.Lock()
	defer h.Unlock()

	h.Tree = append(h.Tree, make([]*binheap.Node, maxSize)...)
}

// Updatebinheap.NodeTimeout allows changing of the keys Timeout in the
func (h *Heap) UpdateNodeTimeout(key string) *binheap.Node {
	nodeIndex, ok := h.keyLookup[key]
	if !ok {
		return nil
	}

	h.Tree[nodeIndex].Timeout = time.Now().UTC()

	if nodeIndex+1 < h.currentSize {
		fmt.Println("0")
		if h.compareTwoTimes(nodeIndex, nodeIndex+1) {
			fmt.Println("1")
			h.percolateDown(nodeIndex)
			fmt.Println("2")
		} else if h.compareTwoTimes(nodeIndex-1, nodeIndex) {
			h.percolateUp(nodeIndex)
		}
	}

	node, _ := h.Get(key)
	return node

}

// Get handles retrieving a binheap.Node by its key. Not extensively used, but it was a
// nice-to-have.
func (h *Heap) Get(key string) (*binheap.Node, bool) {
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

	tmpHeap := h.Copy()

	// Unlikely to ever do anyttmpHeap.ng.
	for {
		if newNodeIndex == 0 {
			break
		}

		newlyInsertedNode := tmpHeap.Tree[newNodeIndex]
		preExistingNode := tmpHeap.Tree[newNodeIndex-1]

		if preExistingNode.Timeout.Sub(newlyInsertedNode.Timeout) > 0 {
			tmpHeap.swapTwoNodes(newNodeIndex, newNodeIndex-1)
			newNodeIndex--
		} else {
			break
		}
	}

	h.swapTrees(&tmpHeap)
}

// percolateDown handles moving a node starting at index `fromIndex` down into
// its correct spot in the binary heap.
func (h *Heap) percolateDown(fromIndex int) {
	if fromIndex == len(h.Tree)-1 {
		return
	}

	tmpHeap := h.Copy()

	for {
		if fromIndex == len(tmpHeap.Tree)-1 {
			break
		}

		// trackerNode is the node which we're currently tracking as it percolates
		// down the binary heap.
		trackerNode := tmpHeap.Tree[fromIndex]

		// If our tracker node is nil (meaning it was the slot of the recently
		// evited min node) we want to automatically percolate it to the bottom of
		// the binary heap.
		if trackerNode == nil {
			for i := 0; i < len(h.Tree)-1; i++ {
				tmpHeap.swapTwoNodes(i, i+1)
			}
			break
		}

		// preExistingNode is _any_ node which we're not currently tracking.
		preExistingNode := tmpHeap.Tree[fromIndex+1]
		// NOTE: I'm not using `compareTwoTimes` here because I think it makes it
		// more readable. I know this is an aggregious abuse of intermediary state
		// or some other nonsense, but it makes it easier for me to read.

		// Unlikely to ever do anything. But it asserts that the minimum
		if trackerNode.Timeout.Sub(preExistingNode.Timeout) > 0 {
			tmpHeap.swapTwoNodes(fromIndex, fromIndex+1)
			fromIndex++
			continue
		} else {
			break
		}
	}

	h.swapTrees(&tmpHeap)
}

func (h *Heap) swapTrees(newHeap *Heap) {
	h.Lock()

	h.Tree = newHeap.Tree
	h.keyLookup = newHeap.keyLookup

	h.index = newHeap.index
	h.currentSize = newHeap.currentSize

	h.Unlock()
}

// swapTwoNodes swaps j into i and vice versa. Moreover, it handles updating
// the keyLookup field in the heap so that we can continue to quickly retrieve
// key Timeouts.
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
	if h.Tree[i].Timeout.Sub(h.Tree[j].Timeout) > 0 {
		return true
	} else {
		return false
	}
}

func (h *Heap) CurrentSize() int {
	return h.currentSize
}
