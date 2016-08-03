package shared

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestNewNode(t *testing.T) {
	expectedReturn := Node{
		Key:     "TestingNewNodeKey",
		Timeout: time.Now().UTC(),
	}

	time.Sleep(5 * time.Millisecond)
	retVal := NewNode("TestingNewNodeKey", time.Now().UTC())

	if expectedReturn.Key != retVal.Key {
		t.Errorf("Expected %v, got %v", expectedReturn.Key, retVal.Key)
	}

	if expectedReturn.Timeout.Nanosecond() >= retVal.Timeout.Nanosecond() {
		t.Errorf("Expected expectedReturn (%v) to be lower than retval (%v)",
			expectedReturn.Timeout,
			retVal.Timeout,
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
	time.Sleep(5 * time.Millisecond)
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
	time.Sleep(5 * time.Millisecond)
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
	time.Sleep(5 * time.Millisecond)

	for i := 0; i < 5; i++ {
		testNode := NewNode(fmt.Sprintf("Node-%v", i), time.Now().UTC())
		testHeap.Insert(testNode)
		time.Sleep(5 * time.Millisecond)
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
		time.Sleep(5 * time.Millisecond)
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

func TestKeyLookup(t *testing.T) {
	testHeap := NewHeapReallocate(5)
	testNode1 := NewNode("TestNode1", time.Now().UTC())
	testNode2 := NewNode("TestNode2", time.Now().UTC())
	testNode3 := NewNode("TestNode3", time.Now().UTC())

	testHeap.Insert(testNode1)
	testHeap.Insert(testNode2)
	testHeap.Insert(testNode3)

	node1Index := testHeap.keyLookup["TestNode1"]
	node2Index := testHeap.keyLookup["TestNode2"]
	node3Index := testHeap.keyLookup["TestNode3"]

	if node1Index != 0 {
		t.Errorf("Incorrect index for node1 %v", node1Index)
	}

	if node2Index != 1 {
		t.Errorf("Incorrect index for node2 %v", node2Index)
	}

	if node3Index != 2 {
		t.Errorf("Incorrect index for node2 %v", node3Index)
	}
}

func TestKeyLookupIndexesProperly(t *testing.T) {
	testHeap := NewHeapReallocate(25)

	keyValues := make([]string, 25)
	for i := 0; i < 25; i++ {
		keyName := fmt.Sprintf("Node-%v", i)
		testNode := NewNode(keyName, time.Now().UTC())
		keyValues[i] = keyName
		testHeap.Insert(testNode)
		time.Sleep(5 * time.Millisecond)
	}

	for i := 0; i < 25; i++ {
		key := keyValues[i]
		keyIndex := testHeap.keyLookup[key]

		if keyIndex != i {
			t.Errorf("Expected key %v to have an index of %v but had index of %v",
				key,
				keyIndex,
			)
		}
	}
}

func TestKeyLookupReadjustsOnEviction(t *testing.T) {
	testHeap := NewHeapReallocate(25)

	keyValues := make([]string, 25)
	for i := 0; i < 25; i++ {
		keyName := fmt.Sprintf("Node-%v", i)
		testNode := NewNode(keyName, time.Now().UTC())
		keyValues[i] = keyName
		testHeap.Insert(testNode)
		time.Sleep(5 * time.Millisecond)
	}

	testHeap.EvictMinNode()

	for i := 0; i < 25; i++ {
		if i == 0 {
			continue
		}

		key := keyValues[i]
		keyIndex := testHeap.keyLookup[key]

		if keyIndex != i-1 {
			t.Errorf("Expected key %v to have an index of %v but had index of %v",
				key,
				keyIndex-1,
				keyIndex,
			)
		}
	}
}

func TestKeyLookupReadjustsOnInsertion(t *testing.T) {
	testHeap := NewHeapReallocate(25)

	originalNode := NewNode("OriginalNode", time.Now().UTC())
	time.Sleep(5 * time.Millisecond)

	keyValues := make([]string, 25)
	for i := 0; i < 25; i++ {
		keyName := fmt.Sprintf("Node-%v", i)
		testNode := NewNode(keyName, time.Now().UTC())
		keyValues[i] = keyName
		testHeap.Insert(testNode)
		time.Sleep(5 * time.Millisecond)
	}

	testHeap.Insert(originalNode)

	for i := 0; i < 25; i++ {
		if i == 0 {
			continue
		}

		key := keyValues[i]
		keyIndex := testHeap.keyLookup[key]

		if keyIndex != i+1 {
			t.Errorf("Expected key %v to have an index of %v but had index of %v",
				key,
				keyIndex+1,
				keyIndex,
			)
		}
	}
}

func TestKeyUpdateTimeoutDoesntBlowUpEverything(t *testing.T) {
	testHeap := NewHeapReallocate(25)

	keyValues := make([]string, 25)
	for i := 0; i < 25; i++ {
		keyName := fmt.Sprintf("Node-%v", i)
		testNode := NewNode(keyName, time.Now().UTC())
		keyValues[i] = keyName
		testHeap.Insert(testNode)
		time.Sleep(5 * time.Millisecond)
	}

	for i := 0; i < 5; i++ {
		index := rand.Intn(25)
		ok := testHeap.UpdateNodeTimeout(keyValues[index])
		if ok == nil {
			t.Errorf("Got weird error, %v index %v", keyValues, index)
		}
		time.Sleep(5 * time.Millisecond)
	}

	for i := 0; i < len(testHeap.Tree)-1; i++ {
		for j := i + 1; j < len(testHeap.Tree)-1; j++ {
			if testHeap.compareTwoTimes(i, j) {
				t.Errorf("%v - %v", testHeap.Tree[i].Key, testHeap.Tree[i].Timeout)
				break
			}
		}
	}
}
