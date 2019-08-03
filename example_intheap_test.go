// Example developed by borrowing from "container/heap/example_intheap_test.go"
// Portions Copyright 2012 by the Go authors
//
// This example demonstrates an integer doubbly-ended heap built using the deheap interface.

package deheap_test

import (
	"fmt"

	"github.com/aalpar/deheap"
)

// An IntHeap is a min-heap of ints.
type IntDeheap []int

func (h IntDeheap) Len() int           { return len(h) }
func (h IntDeheap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntDeheap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *IntDeheap) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(int))
}

func (h *IntDeheap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

// This example inserts several ints into an IntHeap, checks the minimum,
// and removes them in order of priority.
func Example_intHeap() {
	h := &IntDeheap{2, 1, 5, 6}
	deheap.Init(h)
	deheap.Push(h, 3)
	fmt.Printf("minimum: %d\n", (*h)[0])
	for h.Len() > 3 {
		fmt.Printf("%d ", deheap.PopMax(h))
	}
	for h.Len() > 1 {
		fmt.Printf("%d ", deheap.Pop(h))
	}
	fmt.Printf("middle value: %d\n", (*h)[0])
	// Output:
	// minimum: 1
	// 6 5 1 2 middle value: 3
}
