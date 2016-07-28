package olilib_lru

import (
	"time"
)

type HeapAllocationStrategy int

const (
	Realloc HeapAllocationStrategy = iota
	Maintain
)
// Node represents each binary heap node.
type Node struct {
	parent *Node
	child_left *Node
	child_right *Node

	Key string
	timeout time.Time
}

type Heap struct {
	Tree []*Node
	currentSize int
	index int
	allocStrategy HeapAllocationStrategy
}

// NewNode Allocates a new Node. It is not placed in the binary heap at
// allocation. Rather, the caller is expected to later Insert the newly created
// node into the binary heap.
func NewNode(key string, timeout time.Time) *Node {
	return &Node{
		parent: nil,
		child_left: nil,
		child_right: nil,
		Key: key,
		timeout: timeout,
	}
}

func NewHeap(maxSize int) *Heap {
	return &Heap{
		index: maxSize,
		Tree: make([]*Node, maxSize),
	}
}

func (h *Heap) Insert(node *Node) *Node {
	if h.index + 1 >= len(h.Tree) {
		// If we run into the bounds of our heap, we need to either
		// reallocate (if that's what we're wanting to do, or
		// maintain the size and
		if h.allocStrategy == Realloc {
			// The default behavior is to expand the heap by
			// 1.5 times.
			h.ReAllocate(h.index + len(h.Tree) / 2)
		} else {
			// Otherwise, if we're maintaining, we want to evict
			// the root node (The Min Node).
			return h.EvictMinNode()
		}
	}

	h.Tree[h.index] = node
	h.percolateUp(h.index)
	h.index++
	h.currentSize++

	return nil
}

func (h *Heap) EvictMinNode() *Node {
	if h.index == 0 {
		return nil
	}

	retVal := h.Tree[1]

	h.rebalanceHeap(1)

	h.index--
	h.currentSize--
	return retVal
}

// IsEmpty Notifies the caller if the binary heap is empty.
func (h *Heap) IsEmpty() bool {
	return h.currentSize == 0
}
// ReAllocate Handles increasing the size of the underlying binary heap.
func (h *Heap) ReAllocate(maxSize int) {
	h.Tree = append(h.Tree, make([]*Node, maxSize))
}

// percolateUp handles sorting a newly inserted node into its correct position.
func (h *Heap) percolateUp(newNodeIndex int) {
	currentNode := h.Tree[newNodeIndex]
	for currentNode.timeout > currentNode.timeout {
		h.swapTwoNodes(newNodeIndex - 1, newNodeIndex)
		newNodeIndex--

		// If newNodeIndex hits 1, then we're establishing a new root
		// node.
		if newNodeIndex == 1 {
			continue
		}
	}
}

// swapTwoNodes swaps j into i
func (h *Heap) swapTwoNodes(i int, j int) {
	h.Tree[j], h.Tree[i] = h.Tree[i], h.Tree[j]
}
