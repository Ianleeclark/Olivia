package shared

type BinHeap interface {
	// Return a copy of the current BinHeap
	// Copy() BinHeap
	// NOTE: Copy() is disregarded, but needs to be implemented. It
	// seems that, as a contractual obligation, interfaces in golang
	// aren't working optimally and can't return the interface type.
	// MinNode Returns the root node.
	// NOTE: This is a minimum binheap.
	MinNode() *Node
	// Insert inserts a new BinHeapNode into the BinHeap.
	// Moreover, if no realloc strategy is declared, it returns the
	// node to the caller. Verify correct insertion against `nil`.
	Insert(*Node) *Node
	// EvictMinNode removes the root node.
	EvictMinNode() *Node
	// Peek views the node at specified index.
	// Errors are only returned if index is not existing in BinHeap
	Peek(int) (*Node, error)
	// Checks if the binheap is empty.
	IsEmpty() bool
	// Reallocate the size for the binheap.
	// NOTE: If binheap size goes down, the implementation **ought** to
	// evict according to however the implementation sees fit.
	ReAllocate(int)
	UpdateNodeTimeout(string) *Node
	Get(string) (*Node, bool)
	// NOTE: Percolate methods are not required, as a ring-buffer
	// implementation will allow for non-tree-based operations for the
	// binheap.
}
