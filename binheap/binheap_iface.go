package binheap

type LRUStorage interface {
	// Return a copy of the current LRUStorage
	// Copy() LRUStorage
	// NOTE: Copy() is disregarded, but needs to be implemented. It
	// seems that, as a contractual obligation, interfaces in golang
	// aren't working optimally and can't return the interface type.
	// MinNode Returns the root node.
	// NOTE: This is a minimum LRUStorage.
	MinNode() *Node
	// Insert inserts a new LRUStorageNode into the LRUStorage.
	// Moreover, if no realloc strategy is declared, it returns the
	// node to the caller. Verify correct insertion against `nil`.
	Insert(*Node) *Node
	// EvictMinNode removes the root node.
	EvictMinNode() *Node
	// Peek views the node at specified index.
	// Errors are only returned if index is not existing in LRUStorage
	Peek(int) (*Node, error)
	// Checks if the LRUStorage is empty.
	IsEmpty() bool
	// Reallocate the size for the LRUStorage.
	// NOTE: If LRUStorage size goes down, the implementation **ought** to
	// evict according to however the implementation sees fit.
	ReAllocate(int)
	UpdateNodeTimeout(string) *Node
	Get(string) (*Node, bool)
	// NOTE: Percolate methods are not required, as a ring-buffer
	// implementation will allow for non-tree-based operations for the
	// LRUStorage.
	CurrentSize() int
}
