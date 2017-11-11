package binheapv2

import (
	"fmt"
	. "github.com/GrappigPanda/Olivia/binheap"
	"math"
	"sync"
	"time"
)

type Direction int

const (
	INCREMENT Direction = iota
	DECREMENT
)

type BinheapOptimized struct {
	Tree          []*Node
	maxIndex      int
	minIndex      int
	allocStrategy HeapAllocationStrategy
	keyLookup     map[string]int
	sync.Mutex
}

func NewBinheapOptimized(maxSize int) *BinheapOptimized {
	return &BinheapOptimized{
		maxIndex:  0,
		minIndex:  0,
		Tree:      make([]*Node, maxSize),
		keyLookup: make(map[string]int),
	}
}

// NewDoubleIndexArrayReallocate handles intiailizing a new dia object which is able to
// reallocate itself.
func NewBinheapOptimizedReallocate(maxSize int) *BinheapOptimized {
	return &BinheapOptimized{
		maxIndex:      0,
		minIndex:      0,
		Tree:          make([]*Node, maxSize),
		keyLookup:     make(map[string]int),
		allocStrategy: Realloc,
	}
}

func (d *BinheapOptimized) Copy() *BinheapOptimized {
	d.Lock()

	newStorage := NewBinheapOptimized(len(d.Tree))

	for index, element := range d.Tree {
		newStorage.Tree[index] = element
	}

	for k, v := range d.keyLookup {
		newStorage.keyLookup[k] = v
	}

	newStorage.maxIndex = d.maxIndex
	newStorage.minIndex = d.minIndex

	d.Unlock()

	return newStorage
}

func (d *BinheapOptimized) MinNode() *Node {
	return d.Tree[d.maxIndex]
}

func (d *BinheapOptimized) Insert(newNode *Node) *Node {
	var newlyInsertedIndex int

	d.Lock()
	if d.IsEmpty() {
		d.Tree[0] = newNode
		d.maxIndex = 0
		d.minIndex = 0

		newlyInsertedIndex = 0
	} else if d.maxIndex == d.minIndex && !d.IsEmpty() {
		if d.IsFull() {
			newlyInsertedIndex = d.insertNodeAboveOrBelowSingleNode(newNode)
		} else {
			newlyInsertedIndex = d.insertNodeAboveOrBelowSingleNode(newNode)
		}
	} else {
		newlyInsertedIndex = 0
		if d.IsFull() {
			// TODO(ian): Grab the newly freed index
			// TODO(ian): Insert node here
		} else {
			// TODO(ian): Assign newNodeIndex here
		}

		// TODO(ian): Reassign newlyInsertedIndex to whatever the result of the percolation is
	}

	d.keyLookup[newNode.Key] = newlyInsertedIndex

	d.Unlock()
	return newNode
}

func (d *BinheapOptimized) EvictMinNode() (*Node, int) {
	d.Lock()
	defer d.Unlock()

	return d.evictMinNodeLockless()
}

func (d *BinheapOptimized) evictMinNodeLockless() (*Node, int) {
	minNode := d.Tree[d.maxIndex]

	println(fmt.Sprintf("%v", d.keyLookup))
	println(fmt.Sprintf("%v", d.Tree))
	println(minNode == nil)
	println(d.maxIndex)
	println(d.minIndex)
	println(d.Tree[d.maxIndex])
	println(fmt.Sprintf("%v", minNode))

	evictedIndex := d.maxIndex

	delete(d.keyLookup, minNode.Key)
	d.Tree[d.maxIndex] = nil
	d.maxIndex = safeIndex(cap(d.Tree), d.maxIndex, INCREMENT)

	return minNode, evictedIndex
}

// Peek handles looking at the index of the tree.
func (d *BinheapOptimized) Peek(index int) (*Node, error) {
	if index > d.CurrentSize() {
		return nil, fmt.Errorf("Index greater than size of heap.")
	}
	return d.Tree[index], nil
}

// IsEmpty checks to see if the binheap is empty.
func (d *BinheapOptimized) IsEmpty() bool {
	return d.maxIndex == d.minIndex && d.Tree[d.maxIndex] == nil
}

// IsFull checks to see if the binheap is full.
func (d *BinheapOptimized) IsFull() bool {
	// TODO(ian): The easiest way to do this is safeindex either min or max index and see if that == minindex/maxindex.

	count := 0
	for i := 0; i < cap(d.Tree); i++ {
		if d.Tree[i] != nil {
			count++
		}
	}

	return count == cap(d.Tree)
}

// ReAllocate Handles increasing the size of the underlying binary heap.
func (d *BinheapOptimized) ReAllocate(maxSize int) {
	d.Lock()

	// TODO(ian): If `maxSize` decreases, we should do something!
	d.reAllocateLockless(maxSize)

	d.Unlock()
}

// reAllocateLockless Handles increasing the size of the underlying binary heap without a lock. Be careful!
func (d *BinheapOptimized) reAllocateLockless(maxSize int) {
	d.Tree = append(d.Tree, make([]*Node, maxSize)...)
}

// UpdateNodeTimeout allows changing of the keys Timeout in the
func (d *BinheapOptimized) UpdateNodeTimeout(key string) *Node {
	d.Lock()
	nodeIndex, ok := d.keyLookup[key]
	if !ok {
		return nil
	}

	d.Tree[nodeIndex].Timeout = time.Now().UTC()

	nextNodeIndex := safeIndex(cap(d.Tree), nodeIndex, INCREMENT)

	if nextNodeIndex < d.CurrentSize() {
		if d.compareTwoTimes(nodeIndex, nextNodeIndex) {
			d.percolateRight(nodeIndex)
		} else {
			d.percolateLeft(nodeIndex)
		}
	}

	node, _ := d.Get(key)

	d.Unlock()
	return node
}

// Get handles retrieving a Node by its key. Not extensively used, but it was a
// nice-to-have.
func (d *BinheapOptimized) Get(key string) (*Node, bool) {
	if index, ok := d.keyLookup[key]; ok {
		return d.Tree[index], ok
	} else {
		return nil, ok
	}
}

func (d *BinheapOptimized) CurrentSize() int {
	// NOTE: The only time the dereferenced value at d.maxIndex is nil is whenever the binheap is empty.
	if d.Tree[d.maxIndex] != nil {
		// NOTE: If we wrap around the array, we nee to handle that.
		if d.maxIndex < d.minIndex {
			return d.maxIndex + 1 + (cap(d.Tree))
		} else {
			return d.maxIndex - d.minIndex
		}
	} else {
		return 0
	}
}

// compareTwoTimes takes two indexes and compares the `.Nanosecond()` value of
// each in the tree. If the left (i) has an expiration time _after_ the right
// (j), then we return True. Otherwise, if the left (i) has an expiration time
// _before_ the right (j) we return a False.
func (h *BinheapOptimized) compareTwoTimes(i int, j int) bool {
	return compareTimeouts(h.Tree[i].Timeout, h.Tree[j].Timeout)
}

// percolateLeft handles sorting a newly inserted node into its correct position.
// It's very unlikely this function actually ever does anything, as it's only
// called by `Insert`, so newly inserted nodes don't typically have an
// expiration time sooner than nodes already living in the heap.
func (d *BinheapOptimized) percolateLeft(newNodeIndex int) {
	d.Lock()

	if d.maxIndex == d.minIndex && d.maxIndex == newNodeIndex {
		return
	}

	for {
		leftIndex := safeIndex(cap(d.Tree), newNodeIndex, DECREMENT)
		if compareTimeouts(d.Tree[newNodeIndex].Timeout, d.Tree[leftIndex].Timeout) {
			d.swapTwoNodes(newNodeIndex, leftIndex)
		} else {
			break
		}

		newNodeIndex = leftIndex
	}

	d.Unlock()
}

// percolateRight handles moving a node starting at index `fromIndex` down into
// its correct spot in the binary heap.
func (d *BinheapOptimized) percolateRight(fromIndex int) {
	d.Lock()

	if d.maxIndex == d.minIndex && d.maxIndex == fromIndex {
		return
	}

	for {
		break
	}

	d.Unlock()
}

// swapTwoNodes swaps j into i and vice versa. Moreover, it handles updating
// the keyLookup field in the heap so that we can continue to quickly retrieve
// key Timeouts.
func (d *BinheapOptimized) swapTwoNodes(i int, j int) {
	// If we find a value at Tree[i], we can update it in the keylookup,
	// otherwise disregard, as it's a recently evicted node.
	if d.Tree[i] != nil {
		d.keyLookup[d.Tree[i].Key] = j
	}
	if d.Tree[j] != nil {
		d.keyLookup[d.Tree[j].Key] = i
	}

	d.Tree[j], d.Tree[i] = d.Tree[i], d.Tree[j]
}

func (d *BinheapOptimized) insertNodeAboveOrBelowSingleNode(newNode *Node) int {
	var safeidx int
	if compareTimeouts(d.Tree[d.maxIndex].Timeout, newNode.Timeout) {
		safeidx = safeIndex(cap(d.Tree), d.maxIndex, DECREMENT)
		d.Tree[safeidx] = newNode
		d.minIndex--
	} else {
		if d.IsFull() {
			if d.allocStrategy == Realloc {
				if cap(d.Tree) < 10 {
					d.reAllocateLockless(10)
				} else {
					d.reAllocateLockless(int(math.Ceil(float64(cap(d.Tree)) * 1.5)))
				}
			} else {
				_, safeidx = d.evictMinNodeLockless()
			}
		} else {
			safeidx = safeIndex(cap(d.Tree), d.maxIndex, INCREMENT)
		}

		d.Tree[safeidx] = newNode
		d.maxIndex++
	}

	return safeidx
}

func compareTimeouts(time1 time.Time, time2 time.Time) bool {
	return time1.Sub(time2) > 0
}

func safeIndex(treeCapacity, i int, direction Direction) int {
	if direction == INCREMENT {
		nextVal := i + 1
		if nextVal > treeCapacity {
			return nextVal % treeCapacity
		} else {
			return nextVal
		}
	} else {
		nextVal := i - 1
		if nextVal < 0 {
			return treeCapacity - (nextVal * -1)
		} else {
			return nextVal
		}
	}
}
