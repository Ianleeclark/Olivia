package olilib_lru

import (
	"fmt"
	"testing"
	"time"
)

func TestNewNode(t *testing.T) {
	expectedReturn := Node{
		Key:     "TestingNewNodeKey",
		timeout: time.Now().UTC(),
	}

	retVal := NewNode("TestingNewNodeKey", time.Now().UTC())

	if expectedReturn.Key != retVal.Key {
		t.Errorf("Expected %v, got %v", expectedReturn.Key, retVal.Key)
	}

	if expectedReturn.timeout.Nanosecond() >= retVal.timeout.Nanosecond() {
		t.Errorf("Expected expectedReturn (%v) to be lower than retval (%v)",
			expectedReturn.timeout,
			retVal.timeout,
		)
	}
}

func TestNewHeap(t *testing.T) {
	expectedReturn := Heap{
		index: 0,
		Tree:  make([]*Node, 10),
	}

	retVal := NewHeap(10)

	if expectedReturn.index != retVal.index {
		t.Errorf("Expected %v, got %v", expectedReturn.index, retVal.index)
	}

	if len(expectedReturn.Tree) != len(retVal.Tree) {
		t.Errorf("Expected %v, got %v", len(expectedReturn.Tree), len(retVal.Tree))
	}
}

func TestHeapInsertThenReallocate(t *testing.T) {
	testHeap := NewHeapReallocate(1)
	testNode := NewNode("TestHeapInsertThenReallocate", time.Now().UTC())
	time.Sleep(1 * time.Millisecond)
	testNode2 := NewNode("TestHeapInsertThenReallocate2", time.Now().UTC())

	testHeap.Insert(testNode)
	testHeap.Insert(testNode2)

	if len(testHeap.Tree) == 1 {
		t.Errorf("Incorrect allocation strategy, heap didn't reallocate: %v",
			testHeap,
		)
	}
}

func TestInsertAndMinNode(t *testing.T) {
	testHeap := NewHeap(10)
	testNode := NewNode("TestHeapInsert", time.Now().UTC())

	testHeap.Insert(testNode)

	if testHeap.MinNode().Key != "TestHeapInsert" {
		t.Errorf("Failed retrieving min node, got back %v. Tree: %v",
			testHeap.MinNode().Key,
			testHeap.Tree,
		)
	}
}

func TestMinNodeFailNoRootNode(t *testing.T) {
	testHeap := NewHeap(1)

	if testHeap.MinNode() != nil {
		t.Errorf("Expected nil, got %v with a heap of %v",
			testHeap.MinNode(),
			testHeap,
		)
	}
}

func TestSwap(t *testing.T) {
	testHeap := NewHeapReallocate(1)
	testNode := NewNode("Testswap", time.Now().UTC())
	time.Sleep(1 * time.Millisecond)
	testNode2 := NewNode("Testswap2", time.Now().UTC())

	testHeap.Insert(testNode)
	testHeap.Insert(testNode2)

	minNode := testHeap.MinNode()

	testHeap.swapTwoNodes(0, 1)

	newMinNode := testHeap.MinNode()

	if minNode == newMinNode {
		t.Errorf("Expected nodes to swap: MinNode %v - NewMinNode %v -  Heap %v",
			minNode,
			newMinNode,
			testHeap,
		)
	}
}

func TestPercolateUp(t *testing.T) {
	testHeap := NewHeapReallocate(25)

	originalNode := NewNode("Least expiration time", time.Now().UTC())
	time.Sleep(1 * time.Millisecond)

	for i := 0; i < 5; i++ {
		testNode := NewNode(fmt.Sprintf("Node-%v", i), time.Now().UTC())
		testHeap.Insert(testNode)
		time.Sleep(1 * time.Millisecond)
	}

	testHeap.Insert(originalNode)

	if testHeap.MinNode() != originalNode {
		t.Errorf("Expected %v, got %v with a heap of %v",
			originalNode,
			testHeap.MinNode(),
			testHeap.Tree,
		)
	}
}

func TestIsEmpty(t *testing.T) {
	testHeap := NewHeap(10)

	if testHeap.IsEmpty() != true {
		t.Errorf("Expected an empty heap, got %v", testHeap)
	}
}

func TestIsEmptyHasNode(t *testing.T) {
	testHeap := NewHeapReallocate(1)
	testNode := NewNode("Testswap", time.Now().UTC())

	testHeap.Insert(testNode)

	if testHeap.IsEmpty() == true {
		t.Errorf("Expected a non-empty heap, got %v", testHeap)
	}
}

func TestPercolateDown(t *testing.T) {
	testHeap := NewHeapReallocate(25)

	for i := 0; i < 5; i++ {
		testNode := NewNode(fmt.Sprintf("Node-%v", i), time.Now().UTC())
		testHeap.Insert(testNode)
		time.Sleep(1 * time.Millisecond)
	}

	for i := 0; i < 5; i++ {
		expectedReturn := fmt.Sprintf("Node-%v", i)

		if testHeap.MinNode().Key != expectedReturn {
			t.Errorf("Expected %v, got %v with a heap of %v",
				testHeap.MinNode().Key,
				expectedReturn,
				testHeap.Tree,
			)
		}

		retVal := testHeap.EvictMinNode()

		if retVal.Key != expectedReturn {
			t.Errorf("[After Evict] Expected %v, got %v with a heap of %v",
				retVal.Key,
				expectedReturn,
				testHeap.Tree,
			)
		}
	}
}
