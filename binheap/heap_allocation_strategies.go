package binheap

// HeapAllocationStrategy is a type definition used for enum values in handling
// heap reallocation.
type HeapAllocationStrategy int

const (
	Maintain HeapAllocationStrategy = iota
	Realloc
)
