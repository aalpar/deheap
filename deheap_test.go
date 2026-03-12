// Copyright 2019-2026 Aaron H. Alpar
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files
// (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge,
// publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included
// in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
// CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
// TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
// SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
package deheap

import (
	"container/heap"
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"sort"
	"testing"
	"time"
)

type IntHeap []int

func (h IntHeap) Less(i, j int) bool {
	return h[i] < h[j]
}

func (h IntHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h IntHeap) Len() int {
	return len(h)
}

func (h *IntHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}

func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func _newRand() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().Unix()))
}

// TestParent verifies the grandparent (skip-level) index formula:
// parent(i) = ((i+1)/4) - 1. Returns -1 for nodes without a grandparent.
func TestParent(t *testing.T) {
	x := parent(0)
	if x != -1 {
		t.Fatalf("unexpected value: %d", x)
	}
	x = parent(1)
	if x != -1 {
		t.Fatalf("unexpected value: %d", x)
	}
	x = parent(2)
	if x != -1 {
		t.Fatalf("unexpected value: %d", x)
	}
	x = parent(3)
	if x != 0 {
		t.Fatalf("unexpected value: %d", x)
	}
	x = parent(4)
	if x != 0 {
		t.Fatalf("unexpected value: %d", x)
	}
	x = parent(5)
	if x != 0 {
		t.Fatalf("unexpected value: %d", x)
	}
	x = parent(6)
	if x != 0 {
		t.Fatalf("unexpected value: %d", x)
	}
	x = parent(7)
	if x != 1 {
		t.Fatalf("unexpected value: %d", x)
	}
}

// TestHParent verifies the standard binary-tree parent formula:
// hparent(i) = (i-1)/2.
func TestHParent(t *testing.T) {
	x := hparent(0)
	if x != 0 {
		t.Fatalf("unexpected value: %d", x)
	}
	x = hparent(1)
	if x != 0 {
		t.Fatalf("unexpected value: %d", x)
	}
	x = hparent(2)
	if x != 0 {
		t.Fatalf("unexpected value: %d", x)
	}
	x = hparent(3)
	if x != 1 {
		t.Fatalf("unexpected value: %d", x)
	}
	x = hparent(4)
	if x != 1 {
		t.Fatalf("unexpected value: %d", x)
	}
	x = hparent(5)
	if x != 2 {
		t.Fatalf("unexpected value: %d", x)
	}
	x = hparent(6)
	if x != 2 {
		t.Fatalf("unexpected value: %d", x)
	}
	x = hparent(7)
	if x != 3 {
		t.Fatalf("unexpected value: %d", x)
	}
}

// TestLChild verifies the grandchild (skip-level) left child formula:
// lchild(i) = ((i+1)*4) - 1.
func TestLChild(t *testing.T) {
	x := lchild(0)
	if x != 3 {
		t.Fatalf("unexpected value")
	}
	x = lchild(3)
	if x != 15 {
		t.Fatalf("unexpected value")
	}
}

// TestHLChild verifies the standard binary-tree left child formula:
// hlchild(i) = (i*2) + 1.
func TestHLChild(t *testing.T) {
	x := hlchild(0)
	if x != 1 {
		t.Fatalf("unexpected value")
	}
	x = hlchild(1)
	if x != 3 {
		t.Fatalf("unexpected value")
	}
}

// TestMin2 verifies min2, which finds the extremum of two adjacent elements.
// Tests min mode, max mode (out-of-bounds index), and equal elements.
func TestMin2(t *testing.T) {
	h := &IntHeap{10, 1}
	x := min2(h, h.Len(), true, 0)
	if x != 1 {
		t.Fatalf("unexpected value")
	}

	h = &IntHeap{1, 10}
	x = min2(h, h.Len(), true, 0)
	if x != 0 {
		t.Fatalf("unexpected value")
	}

	x = min2(h, h.Len(), false, 9)
	if x != 9 {
		t.Fatalf("unexpected value")
	}

	h = &IntHeap{10, 10}
	x = min2(h, h.Len(), true, 0)
	if x != 0 {
		t.Fatalf("unexpected value")
	}

}

// TestLevel verifies the tree level computation: level(i) = floor(log2(i+1)).
func TestLevel(t *testing.T) {
	x := level(0)
	if x != 0 {
		t.Fatalf("unexpected value")
	}
	x = level(1)
	if x != 1 {
		t.Fatalf("unexpected value")
	}
	x = level(2)
	if x != 1 {
		t.Fatalf("unexpected value")
	}
}

// TestBubbleUp verifies the upward sift operation on hand-crafted heaps.
// Cases: no-op (already in place), swap with binary parent only,
// swap across grandparent levels (7-element and 15-element heaps).
func TestBubbleUp(t *testing.T) {

	h := &IntHeap{1, 15, 2}
	bubbleup(h, isMinHeap(2), 2)
	if !reflect.DeepEqual(h, &IntHeap{1, 15, 2}) {
		t.Fatalf("unexpected value: %v", h)
	}

	h = &IntHeap{1, 4, 10}
	bubbleup(h, isMinHeap(2), 2)
	if !reflect.DeepEqual(h, &IntHeap{1, 4, 10}) {
		t.Fatalf("unexpected value: %v", h)
	}

	h = &IntHeap{2, 15, 13, 4, 6, 8, 1}
	bubbleup(h, isMinHeap(6), 6)
	if !Verify(h) {
		t.Fatalf("unexpected value: %v", h)
	}

	h = &IntHeap{1, 15, 14, 2, 3, 4, 5, 13, 12, 11, 10, 6, 7, 8, 9}
	bubbleup(h, isMinHeap(14), 14)
	if !Verify(h) {
		t.Fatalf("unexpected value: %v", h)
	}

}

// TestBubbleDown verifies the downward sift operation on hand-crafted heaps.
// Cases: 3-element swap, 7-element with grandchild swap, 12-element with
// multi-level cascade, and two larger heaps validated via Verify.
func TestBubbleDown(t *testing.T) {

	h := &IntHeap{15, 1, 2}
	bubbledown(h, h.Len(), isMinHeap(0), 0)
	if !reflect.DeepEqual(h, &IntHeap{1, 15, 2}) {
		t.Fatalf("unexpected value: %v", h)
	}

	h = &IntHeap{5, 7, 4, 6, 1, 3, 2}
	bubbledown(h, h.Len(), isMinHeap(0), 0)
	if !reflect.DeepEqual(h, &IntHeap{1, 7, 4, 6, 5, 3, 2}) {
		t.Fatalf("unexpected value: %v", h)
	}

	h = &IntHeap{10, 8, 12, 1, 2, 9, 10, 5, 3, 4, 6, 11}
	bubbledown(h, h.Len(), isMinHeap(0), 0)
	if !reflect.DeepEqual(h, &IntHeap{1, 10, 12, 3, 2, 9, 10, 5, 8, 4, 6, 11}) {
		t.Fatalf("unexpected value: %v", h)
	}

	h = &IntHeap{14, 15, 12, 4, 2, 3, 5, 13}
	bubbledown(h, h.Len(), isMinHeap(0), 0)
	if !Verify(h) {
		t.Fatalf("unexpected value: %v", h)
	}

	h = &IntHeap{13, 14, 15, 3, 4, 5, 6, 7, 8, 9, 10}
	bubbledown(h, h.Len(), isMinHeap(0), 0)
	if !Verify(h) {
		t.Fatalf("unexpected value: %v", h)
	}

}

// TestMin4 verifies min4, which finds the extremum of up to 4 consecutive
// elements. Tests min at each of the 4 positions, equal values, and a
// single-element slice (exercises early-return branches).
func TestMin4(t *testing.T) {
	h := &IntHeap{3, 1, 2, 4}
	x := min4(h, h.Len(), true, 0)
	if x != 1 {
		t.Fatalf("unexpected value")
	}

	h = &IntHeap{1, 3, 2, 4}
	x = min4(h, h.Len(), true, 0)
	if x != 0 {
		t.Fatalf("unexpected value")
	}

	h = &IntHeap{2, 3, 1, 4}
	x = min4(h, h.Len(), true, 0)
	if x != 2 {
		t.Fatalf("unexpected value")
	}

	h = &IntHeap{2, 3, 4, 1}
	x = min4(h, h.Len(), true, 0)
	if x != 3 {
		t.Fatalf("unexpected value")
	}

	h = &IntHeap{1, 1, 2, 2}
	x = min4(h, h.Len(), true, 0)
	if x != 0 {
		t.Fatalf("unexpected value")
	}

	h = &IntHeap{2, 2, 1, 1}
	x = min4(h, h.Len(), true, 0)
	if x != 2 {
		t.Fatalf("unexpected value")
	}

	h = &IntHeap{2}
	x = min4(h, h.Len(), true, 0)
	if x != 0 {
		t.Fatalf("unexpected value")
	}
}

// TestMin3 verifies min3, which finds the extremum of up to 3 elements
// given explicit indices. Tests all permutations, equal values, and
// out-of-bounds indices.
func TestMin3(t *testing.T) {
	h := &IntHeap{3, 1, 2}
	x := min3(h, h.Len(), true, 0, 1, 2)
	if x != 1 {
		t.Fatalf("unexpected value")
	}

	h = &IntHeap{1, 3, 2}
	x = min3(h, h.Len(), true, 0, 1, 2)
	if x != 0 {
		t.Fatalf("unexpected value")
	}

	h = &IntHeap{2, 3, 1}
	x = min3(h, h.Len(), true, 0, 1, 2)
	if x != 2 {
		t.Fatalf("unexpected value")
	}

	h = &IntHeap{1, 1, 2}
	x = min3(h, h.Len(), true, 0, 1, 2)
	if x != 0 && x != 1 {
		t.Fatalf("unexpected value")
	}

	h = &IntHeap{2, 1, 1}
	x = min3(h, h.Len(), true, 0, 1, 2)
	if x != 2 && x != 1 {
		t.Fatalf("unexpected value")
	}

	h = &IntHeap{2}
	x = min3(h, h.Len(), true, 0, 1, 2)
	if x != 0 {
		t.Fatalf("unexpected value")
	}
}

// TestInit verifies that Init produces a valid min-max heap from an
// arbitrary unordered slice.
func TestInit(t *testing.T) {
	h := &IntHeap{15, 1, 2, 14, 13, 12, 11, 3, 4, 5, 6, 7, 8, 9, 10}
	Init(h)
	if !Verify(h) {
		t.Fatalf("not a valid heap: %v", *h)
	}
}

// TestEmptyPop verifies that Pop and PopMax on an empty heap return nil
// without panicking.
func TestEmptyPop(t *testing.T) {
	h := &IntHeap{}
	Init(h)
	if v := Pop(h); v != nil {
		t.Fatalf("Pop on empty heap = %v, want nil", v)
	}
	if v := PopMax(h); v != nil {
		t.Fatalf("PopMax on empty heap = %v, want nil", v)
	}
}

// TestSingleElement verifies Pop and PopMax each work correctly on a
// single-element heap, and that the heap is empty afterward.
func TestSingleElement(t *testing.T) {
	h := &IntHeap{42}
	Init(h)
	if v := Pop(h).(int); v != 42 {
		t.Fatalf("Pop single = %d, want 42", v)
	}
	if h.Len() != 0 {
		t.Fatalf("Len after Pop = %d, want 0", h.Len())
	}

	h = &IntHeap{42}
	Init(h)
	if v := PopMax(h).(int); v != 42 {
		t.Fatalf("PopMax single = %d, want 42", v)
	}
	if h.Len() != 0 {
		t.Fatalf("Len after PopMax = %d, want 0", h.Len())
	}
}

// TestTwoElements verifies Pop and PopMax on a two-element heap return
// elements in the correct order (min first for Pop, max first for PopMax).
func TestTwoElements(t *testing.T) {
	h := &IntHeap{3, 7}
	Init(h)
	if v := Pop(h).(int); v != 3 {
		t.Fatalf("Pop first = %d, want 3", v)
	}
	if v := Pop(h).(int); v != 7 {
		t.Fatalf("Pop second = %d, want 7", v)
	}

	h = &IntHeap{3, 7}
	Init(h)
	if v := PopMax(h).(int); v != 7 {
		t.Fatalf("PopMax first = %d, want 7", v)
	}
	if v := PopMax(h).(int); v != 3 {
		t.Fatalf("PopMax second = %d, want 3", v)
	}

	// Two equal elements
	h = &IntHeap{5, 5}
	Init(h)
	if v := Pop(h).(int); v != 5 {
		t.Fatalf("Pop equal = %d, want 5", v)
	}
	if v := PopMax(h).(int); v != 5 {
		t.Fatalf("PopMax equal = %d, want 5", v)
	}
}

// TestRemoveSingleElement verifies Remove(h, 0) on a 1-element heap returns
// the correct value and leaves the heap empty.
func TestRemoveSingleElement(t *testing.T) {
	h := &IntHeap{42}
	Init(h)
	v := Remove(h, 0).(int)
	if v != 42 {
		t.Fatalf("Remove(0) = %d, want 42", v)
	}
	if h.Len() != 0 {
		t.Fatalf("Len after Remove = %d, want 0", h.Len())
	}
}

// TestRemoveTwoElements verifies Remove on a 2-element heap at both indices,
// including the case where both elements are equal.
func TestRemoveTwoElements(t *testing.T) {
	// Remove index 0 (min element)
	h := &IntHeap{3, 7}
	Init(h)
	v := Remove(h, 0).(int)
	if v != 3 {
		t.Fatalf("Remove(0) = %d, want 3", v)
	}
	if h.Len() != 1 {
		t.Fatalf("Len = %d, want 1", h.Len())
	}
	if r := Pop(h).(int); r != 7 {
		t.Fatalf("remaining = %d, want 7", r)
	}

	// Remove index 1 (max element)
	h = &IntHeap{3, 7}
	Init(h)
	v = Remove(h, 1).(int)
	if v != 7 {
		t.Fatalf("Remove(1) = %d, want 7", v)
	}
	if h.Len() != 1 {
		t.Fatalf("Len = %d, want 1", h.Len())
	}
	if r := Pop(h).(int); r != 3 {
		t.Fatalf("remaining = %d, want 3", r)
	}

	// Equal elements
	h = &IntHeap{5, 5}
	Init(h)
	v = Remove(h, 0).(int)
	if v != 5 {
		t.Fatalf("Remove(0) equal = %d, want 5", v)
	}
	if r := Pop(h).(int); r != 5 {
		t.Fatalf("remaining equal = %d, want 5", r)
	}

	h = &IntHeap{5, 5}
	Init(h)
	v = Remove(h, 1).(int)
	if v != 5 {
		t.Fatalf("Remove(1) equal = %d, want 5", v)
	}
	if r := Pop(h).(int); r != 5 {
		t.Fatalf("remaining equal = %d, want 5", r)
	}
}

// TestFix verifies Fix restores the heap property after modifying a single
// element. Cases cover both directions (bubbledown and bubbleup), both
// level types (min and max), root vs interior vs leaf, and no-op.
func TestFix(t *testing.T) {
	// Increase root (min level) — needs bubbledown.
	h := &IntHeap{1, 9, 5, 4, 6, 3, 2}
	(*h)[0] = 100
	Fix(h, 0)
	if !Verify(h) {
		t.Fatalf("Fix root increase: not a heap: %v", *h)
	}

	// Decrease root — already min, no movement needed.
	h = &IntHeap{1, 9, 5, 4, 6, 3, 2}
	(*h)[0] = -1
	Fix(h, 0)
	if !Verify(h) {
		t.Fatalf("Fix root decrease: not a heap: %v", *h)
	}

	// Decrease min-level node to new global min — bubbleup through min chain.
	h = &IntHeap{1, 9, 5, 4, 6, 3, 2}
	(*h)[3] = -1
	Fix(h, 3)
	if !Verify(h) {
		t.Fatalf("Fix min node decrease: not a heap: %v", *h)
	}
	if (*h)[0] != -1 {
		t.Fatalf("Fix min node decrease: root = %d, want -1", (*h)[0])
	}

	// Increase min-level leaf — bubbleup through max chain.
	h = &IntHeap{1, 9, 5, 4, 6, 3, 2}
	(*h)[6] = 100
	Fix(h, 6)
	if !Verify(h) {
		t.Fatalf("Fix min leaf increase: not a heap: %v", *h)
	}

	// Decrease max-level node — needs bubbledown or cross-level swap.
	h = &IntHeap{1, 9, 5, 4, 6, 3, 2}
	(*h)[1] = 0
	Fix(h, 1)
	if !Verify(h) {
		t.Fatalf("Fix max node decrease: not a heap: %v", *h)
	}

	// Increase max-level node — no movement needed (already max).
	h = &IntHeap{1, 9, 5, 4, 6, 3, 2}
	(*h)[1] = 100
	Fix(h, 1)
	if !Verify(h) {
		t.Fatalf("Fix max node increase: not a heap: %v", *h)
	}

	// No-op: value unchanged.
	h = &IntHeap{1, 9, 5, 4, 6, 3, 2}
	Fix(h, 3)
	if !Verify(h) {
		t.Fatalf("Fix no-op: not a heap: %v", *h)
	}

	// Single element.
	h = &IntHeap{42}
	Init(h)
	(*h)[0] = 99
	Fix(h, 0)
	if !Verify(h) {
		t.Fatalf("Fix single: not a heap: %v", *h)
	}

	// Two elements: fix min.
	h = &IntHeap{3, 7}
	Init(h)
	(*h)[0] = 10
	Fix(h, 0)
	if !Verify(h) {
		t.Fatalf("Fix two (min): not a heap: %v", *h)
	}

	// Two elements: fix max.
	h = &IntHeap{3, 7}
	Init(h)
	(*h)[1] = 1
	Fix(h, 1)
	if !Verify(h) {
		t.Fatalf("Fix two (max): not a heap: %v", *h)
	}
}

// TestFixRandomized modifies random elements in random heaps and verifies
// Fix restores the heap property. Validates both structure (Verify) and
// content (sorted drain matches oracle).
func TestFixRandomized(t *testing.T) {
	s := rand.New(rand.NewSource(time.Now().UnixNano()))

	for iter := 0; iter < 1000; iter++ {
		n := s.Intn(64) + 2
		h := randIntHeapWithDups(t, n, 0.1)

		// Build sorted oracle from current heap state.
		oracle := make([]int, h.Len())
		copy(oracle, *h)
		sort.Ints(oracle)

		// Modify a random element and fix.
		idx := s.Intn(h.Len())
		oldVal := (*h)[idx]
		newVal := s.Intn(n*2) - n
		(*h)[idx] = newVal
		Fix(h, idx)

		// Update oracle: remove old, insert new.
		j := sort.SearchInts(oracle, oldVal)
		oracle = append(oracle[:j], oracle[j+1:]...)
		k := sort.SearchInts(oracle, newVal)
		oracle = append(oracle, 0)
		copy(oracle[k+1:], oracle[k:])
		oracle[k] = newVal

		if !Verify(h) {
			t.Fatalf("iter %d: Fix(%d) broke heap: %v", iter, idx, *h)
		}

		for di, want := range oracle {
			got := Pop(h).(int)
			if got != want {
				t.Fatalf("iter %d: Pop[%d] = %d, want %d", iter, di, got, want)
			}
		}
	}
}

// TestPush verifies Push maintains the heap property under three insertion
// patterns: ascending order, descending order (validated after each push),
// and alternating high/low values.
func TestPush(t *testing.T) {

	h := &IntHeap{}
	for i := 0; i < 32; i++ {
		Push(h, i)
	}
	if !Verify(h) {
		t.Fatalf("not a valid heap: %v", *h)
	}

	h = &IntHeap{}
	for i := 3; i >= 0; i-- {
		Push(h, i)
		if !Verify(h) {
			t.Fatalf("not a valid heap: %v", *h)
		}
	}

	h = &IntHeap{}
	for i := 32; i >= 0; i-- {
		q := (i%2 == 0)
		k := i
		if q {
			k = 32 - i
		}
		Push(h, k)
	}
	if !Verify(h) {
		t.Fatalf("not a valid heap: %v", *h)
	}

}

// TestV1PushPopEmpty verifies PushPop on an empty heap returns o.
func TestV1PushPopEmpty(t *testing.T) {
	h := &IntHeap{}
	if v := PushPop(h, 42).(int); v != 42 {
		t.Fatalf("PushPop empty = %d, want 42", v)
	}
	if h.Len() != 0 {
		t.Fatalf("Len = %d, want 0", h.Len())
	}
}

// TestV1PushPopNewMin verifies PushPop returns o when it's the new min.
func TestV1PushPopNewMin(t *testing.T) {
	h := &IntHeap{1, 9, 5, 4, 6, 3, 2}
	if v := PushPop(h, 0).(int); v != 0 {
		t.Fatalf("PushPop new min = %d, want 0", v)
	}
	if h.Len() != 7 {
		t.Fatalf("Len = %d, want 7", h.Len())
	}
	if !Verify(h) {
		t.Fatal("Verify() failed")
	}
}

// TestV1PushPopReplaces verifies PushPop evicts the current min when o > min.
func TestV1PushPopReplaces(t *testing.T) {
	h := &IntHeap{1, 9, 5, 4, 6, 3, 2}
	if v := PushPop(h, 4).(int); v != 1 {
		t.Fatalf("PushPop = %d, want 1", v)
	}
	if h.Len() != 7 {
		t.Fatalf("Len = %d, want 7", h.Len())
	}
	if !Verify(h) {
		t.Fatal("Verify() failed")
	}
}

// TestV1PushPopMaxEmpty verifies PushPopMax on an empty heap returns o.
func TestV1PushPopMaxEmpty(t *testing.T) {
	h := &IntHeap{}
	if v := PushPopMax(h, 42).(int); v != 42 {
		t.Fatalf("PushPopMax empty = %d, want 42", v)
	}
	if h.Len() != 0 {
		t.Fatalf("Len = %d, want 0", h.Len())
	}
}

// TestV1PushPopMaxNewMax verifies PushPopMax returns o when it's the new max.
func TestV1PushPopMaxNewMax(t *testing.T) {
	h := &IntHeap{1, 9, 5, 4, 6, 3, 2}
	if v := PushPopMax(h, 100).(int); v != 100 {
		t.Fatalf("PushPopMax new max = %d, want 100", v)
	}
	if h.Len() != 7 {
		t.Fatalf("Len = %d, want 7", h.Len())
	}
	if !Verify(h) {
		t.Fatal("Verify() failed")
	}
}

// TestV1PushPopMaxReplaces verifies PushPopMax evicts the current max when o < max.
func TestV1PushPopMaxReplaces(t *testing.T) {
	h := &IntHeap{1, 9, 5, 4, 6, 3, 2}
	if v := PushPopMax(h, 4).(int); v != 9 {
		t.Fatalf("PushPopMax = %d, want 9", v)
	}
	if h.Len() != 7 {
		t.Fatalf("Len = %d, want 7", h.Len())
	}
	if !Verify(h) {
		t.Fatal("Verify() failed")
	}
}

// TestV1PushPopRandomized verifies PushPop against Push+Pop oracle.
func TestV1PushPopRandomized(t *testing.T) {
	s := rand.New(rand.NewSource(time.Now().UnixNano()))
	for iter := 0; iter < 1000; iter++ {
		n := s.Intn(64) + 1
		h := &IntHeap{}
		oracle := &IntHeap{}
		for i := 0; i < n; i++ {
			v := s.Intn(n)
			Push(h, v)
			Push(oracle, v)
		}
		o := s.Intn(n * 2)
		got := PushPop(h, o).(int)
		Push(oracle, o)
		want := Pop(oracle).(int)
		if got != want {
			t.Fatalf("iter %d: PushPop(%d) = %d, want %d", iter, o, got, want)
		}
		if !Verify(h) {
			t.Fatalf("iter %d: Verify() failed", iter)
		}
	}
}

// TestV1PushPopMaxRandomized verifies PushPopMax against Push+PopMax oracle.
func TestV1PushPopMaxRandomized(t *testing.T) {
	s := rand.New(rand.NewSource(time.Now().UnixNano()))
	for iter := 0; iter < 1000; iter++ {
		n := s.Intn(64) + 1
		h := &IntHeap{}
		oracle := &IntHeap{}
		for i := 0; i < n; i++ {
			v := s.Intn(n)
			Push(h, v)
			Push(oracle, v)
		}
		o := s.Intn(n * 2)
		got := PushPopMax(h, o).(int)
		Push(oracle, o)
		want := PopMax(oracle).(int)
		if got != want {
			t.Fatalf("iter %d: PushPopMax(%d) = %d, want %d", iter, o, got, want)
		}
		if !Verify(h) {
			t.Fatalf("iter %d: Verify() failed", iter)
		}
	}
}

// TestPops is the main table-driven correctness test. For each of 7
// hand-crafted valid min-max heaps (including cases with duplicates),
// it verifies:
//   - Pop drains elements in ascending order (q0)
//   - PopMax drains elements in descending order (q1)
//   - Remove at every index produces the correct element and the
//     remaining elements Pop in sorted order (q2)
func TestPops(t *testing.T) {
	ts := []struct {
		h  IntHeap
		q0 []int
		q1 []int
		q2 []int
	}{
		{
			IntHeap{1, 4, 3, 2, 4, 3},
			[]int{1, 2, 3, 3, 4, 4},
			[]int{4, 4, 3, 3, 2, 1},
			[]int{1, 2, 3, 3, 4, 4},
		},
		{
			IntHeap{1, 3, 2, 2, 3},
			[]int{1, 2, 2, 3, 3},
			[]int{3, 3, 2, 2, 1},
			[]int{1, 2, 2, 3, 3},
		},
		{
			IntHeap{1, 5, 4, 2, 3, 3},
			[]int{1, 2, 3, 3, 4, 5},
			[]int{5, 4, 3, 3, 2, 1},
			[]int{1, 2, 3, 3, 4, 5},
		},
		{
			IntHeap{1, 4, 4, 2, 3, 3},
			[]int{1, 2, 3, 3, 4, 4},
			[]int{4, 4, 3, 3, 2, 1},
			[]int{1, 2, 3, 3, 4, 4},
		},
		{
			IntHeap{1, 4, 4, 3, 3, 2},
			[]int{1, 2, 3, 3, 4, 4},
			[]int{4, 4, 3, 3, 2, 1},
			[]int{1, 2, 3, 3, 4, 4},
		},
		{
			IntHeap{1, 8, 9, 2, 3, 4, 5, 6, 7},
			[]int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			[]int{9, 8, 7, 6, 5, 4, 3, 2, 1},
			[]int{1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		{
			IntHeap{1, 9, 5, 4, 7, 3, 2, 6, 8},
			[]int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			[]int{9, 8, 7, 6, 5, 4, 3, 2, 1},
			[]int{1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
	}

	for ti, tv := range ts {

		if !Verify(&tv.h) {
			t.Fatalf("not a valid heap: %v", tv.h)
		}

		t.Run(fmt.Sprintf("%d Pop", ti), func(t1 *testing.T) {

			hp0 := make(IntHeap, len(tv.h))
			copy(hp0, tv.h)

			for _, y := range tv.q0 {
				if y != Pop(&hp0) {
					t.Fatalf("unexpected value")
				}
			}

		})

		t.Run(fmt.Sprintf("%d PopMax", ti), func(t1 *testing.T) {

			hp1 := make(IntHeap, len(tv.h))
			copy(hp1, tv.h)

			for _, y := range tv.q1 {
				if y != PopMax(&hp1) {
					t.Fatalf("unexpected value")
				}
			}

		})

		for k := 0; k < len(tv.h); k++ {

			t.Run(fmt.Sprintf("%d Remove %d", ti, k), func(t1 *testing.T) {

				hp2 := make(IntHeap, len(tv.h))
				copy(hp2, tv.h)
				q2 := make([]int, len(tv.q2))
				copy(q2, tv.q2)

				c := Remove(&hp2, k).(int)
				j := sort.SearchInts(q2, c)
				q2 = append(q2[:j], q2[j+1:]...)

				for j, y := range q2 {
					x := Pop(&hp2)
					if y != x {
						t.Fatalf("unexpected value: %d %d", j, x)
					}
				}
			})
		}

	}
}

// TestRemove verifies Remove on a 10-element heap by checking the exact
// internal array layout after each removal. This is the only test that
// asserts on the specific heap structure (not just sorted output),
// serving as a regression test for the sift algorithm.
func TestRemove(t *testing.T) {

	h := &IntHeap{0, 9, 5, 6, 1, 2, 4, 8, 7, 3}

	x := Remove(h, 9).(int)
	if x != 3 {
		t.Fatalf("unexpected value")
	}
	if !Verify(h) {
		t.Fatalf("unexpected value")
	}
	if !reflect.DeepEqual(h, &IntHeap{0, 9, 5, 6, 1, 2, 4, 8, 7}) {
		t.Fatalf("unexpected value")
	}

	x = Remove(h, 2).(int)
	if x != 5 {
		t.Fatalf("unexpected value")
	}
	if !Verify(h) {
		t.Fatalf("unexpected value")
	}
	if !reflect.DeepEqual(h, &IntHeap{0, 9, 7, 6, 1, 2, 4, 8}) {
		t.Fatalf("unexpected value")
	}

	x = Remove(h, 0).(int)
	if x != 0 {
		t.Fatalf("unexpected value")
	}
	if !Verify(h) {
		t.Fatalf("unexpected value")
	}
	if !reflect.DeepEqual(h, &IntHeap{1, 9, 7, 6, 8, 2, 4}) && !reflect.DeepEqual(h, &IntHeap{1, 9, 8, 6, 7, 2, 4}) {
		t.Fatalf("unexpected value: %v", h)
	}

}

// TestOps is the comprehensive randomized stress test. Over 1000 iterations
// with random heap sizes (2-65 elements, with ~10% duplicate values):
//   - First half: randomly interleaves Pop and PopMax, verifying monotonicity
//     of returned values and heap validity after each operation.
//   - Second half: randomly Removes elements one at a time, verifying heap
//     validity after each removal.
func TestOps(t *testing.T) {

	for k := 0; k < 1000; k++ {

		s := rand.New(rand.NewSource(time.Now().Unix()))

		N := s.Intn(64) + 2

		h := randIntHeapWithDups(t, N, 0.1)
		y0 := math.MinInt32
		y1 := math.MaxInt32
		for h.Len() > 0 {
			if s.Intn(2) == 0 {
				x := Pop(h).(int)
				if x < y0 || x > y1 {
					t.Fatalf("unexpected value: %d %d", x, y0)
				}
				y0 = x
			} else {
				x := PopMax(h).(int)
				if x > y1 || x < y0 {
					t.Fatalf("unexpected value: %d %d", x, y1)
				}
				y1 = x
			}
			if !Verify(h) {
				t.Fatalf("unexpected value")
			}
		}

		h = randIntHeapWithDups(t, N, 0.1)
		for h.Len() > 0 {
			h0 := make([]int, h.Len())
			copy(h0, []int(*h))
			x := s.Intn(h.Len())
			Remove(h, x)
			if !Verify(h) {
				t.Fatalf("unexpected value: removed %d\n%v\n%v", x, h0, h)
			}
		}

	}

}

func BenchmarkMin4(b *testing.B) {

	r := &[]int{}
	for i := 0; i < b.N; i++ {
		*r = append(*r, i)
	}

	s := _newRand()
	s.Shuffle(len(*r), func(i, j int) { (*r)[i], (*r)[j] = (*r)[j], (*r)[i] })

	h := (*IntHeap)(r)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		min4(h, h.Len(), true, i)
	}

}

func BenchmarkBaselinePush(b *testing.B) {

	r := []int{}
	for i := 0; i < b.N; i++ {
		r = append(r, i)
	}

	s := _newRand()
	s.Shuffle(len(r), func(i, j int) { r[i], r[j] = r[j], r[i] })

	b.ResetTimer()

	a := []int{}
	for _, q := range r {
		a = append(a, q)
	}

}

func BenchmarkPush(b *testing.B) {

	r := []int{}
	for i := 0; i < b.N; i++ {
		r = append(r, i)
	}

	s := _newRand()
	s.Shuffle(len(r), func(i, j int) { r[i], r[j] = r[j], r[i] })

	b.ResetTimer()

	h := &IntHeap{}
	for _, q := range r {
		Push(h, q)
	}

}

func BenchmarkPop(b *testing.B) {

	r := []int{}
	for i := 0; i < b.N; i++ {
		r = append(r, i)
	}

	s := _newRand()
	s.Shuffle(len(r), func(i, j int) { r[i], r[j] = r[j], r[i] })

	h := &IntHeap{}
	for _, q := range r {
		Push(h, q)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Pop(h)
	}

}

func BenchmarkPopMax(b *testing.B) {

	r := []int{}
	for i := 0; i < b.N; i++ {
		r = append(r, i)
	}

	s := _newRand()
	s.Shuffle(len(r), func(i, j int) { r[i], r[j] = r[j], r[i] })

	h := &IntHeap{}
	for _, q := range r {
		Push(h, q)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		PopMax(h)
	}

}

func BenchmarkPushPop(b *testing.B) {

	r := []int{}
	for i := 0; i < b.N; i++ {
		r = append(r, i)
	}

	s := _newRand()
	s.Shuffle(len(r), func(i, j int) { r[i], r[j] = r[j], r[i] })

	b.ResetTimer()

	h := &IntHeap{}
	for _, q := range r {
		Push(h, q)
	}
	for i := 0; i < b.N; i++ {
		Pop(h)
	}

}

func BenchmarkHeapPushPop(b *testing.B) {

	r := []int{}
	for i := 0; i < b.N; i++ {
		r = append(r, i)
	}

	s := _newRand()
	s.Shuffle(len(r), func(i, j int) { r[i], r[j] = r[j], r[i] })

	b.ResetTimer()

	h := &IntHeap{}
	for _, q := range r {
		heap.Push(h, q)
	}
	for i := 0; i < b.N; i++ {
		heap.Pop(h)
	}

}

func BenchmarkHeapPop(b *testing.B) {

	r := []int{}
	for i := 0; i < b.N; i++ {
		r = append(r, i)
	}

	s := _newRand()
	s.Shuffle(len(r), func(i, j int) { r[i], r[j] = r[j], r[i] })

	h := &IntHeap{}
	for _, q := range r {
		heap.Push(h, q)
	}

	b.ResetTimer()

	x := Pop(h).(int)
	for i := 0; i < b.N-1; i++ {
		y := heap.Pop(h).(int)
		if x > y {
			panic("bad")
		}
	}

}

func BenchmarkHeapPush(b *testing.B) {

	r := []int{}
	for i := 0; i < b.N; i++ {
		r = append(r, i)
	}

	s := _newRand()
	s.Shuffle(len(r), func(i, j int) { r[i], r[j] = r[j], r[i] })

	b.ResetTimer()

	h := &IntHeap{}
	for _, q := range r {
		heap.Push(h, q)
	}

}

// FuzzV1PushPop interprets a byte sequence as heap commands: '<' = Pop,
// '>' = PopMax, anything else = Push(byte). Validates every result against
// a sorted-slice oracle.
func FuzzV1PushPop(f *testing.F) {
	f.Add([]byte{10, 5, 3, 8, '<', '>', 1, '<', 7, '>'})
	f.Add([]byte{1, 1, 1, '<', '<', '<'})
	f.Add([]byte{})

	f.Fuzz(func(t *testing.T, data []byte) {
		h := &IntHeap{}
		Init(h)
		var oracle sortedOracle

		for _, c := range data {
			switch c {
			case '<':
				if h.Len() > 0 {
					got := Pop(h).(int)
					want := oracle.popMin()
					if got != want {
						t.Fatalf("Pop: got %d, want %d", got, want)
					}
				}
			case '>':
				if h.Len() > 0 {
					got := PopMax(h).(int)
					want := oracle.popMax()
					if got != want {
						t.Fatalf("PopMax: got %d, want %d", got, want)
					}
				}
			default:
				v := int(c)
				Push(h, v)
				oracle.push(v)
			}
			if !Verify(h) {
				t.Fatalf("Verify() failed")
			}
			if h.Len() != oracle.len() {
				t.Fatalf("length mismatch: heap=%d, oracle=%d", h.Len(), oracle.len())
			}
		}
	})
}

// FuzzV1Remove pushes all bytes onto the heap, then removes elements at
// deterministic-random indices until drained. Any panic from an invalid
// heap state indicates a bug in the Remove/sift logic.
func FuzzV1Remove(f *testing.F) {
	f.Add([]byte{10, 5, 3, 8, 1, 7})
	f.Add([]byte{1, 1, 1, 1})
	f.Add([]byte{})

	f.Fuzz(func(t *testing.T, data []byte) {
		h := &IntHeap{}
		for _, c := range data {
			Push(h, int(c))
		}

		// Remove elements at random-ish positions and verify heap drains correctly
		s := rand.New(rand.NewSource(int64(len(data))))
		for h.Len() > 0 {
			idx := s.Intn(h.Len())
			Remove(h, idx)
			if !Verify(h) {
				t.Fatalf("Verify() failed after Remove(%d)", idx)
			}
		}
	})
}

// FuzzV1Fix pushes the first half of bytes onto a heap, then interprets the
// second half as (index, newValue) pairs for Fix operations. Validates the
// heap property after every Fix.
func FuzzV1Fix(f *testing.F) {
	f.Add([]byte{5, 3, 8, 1, 9, 2, 7, 4, 6})
	f.Add([]byte{1, 1, 1, 1, 1, 1})
	f.Add([]byte{})

	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) < 2 {
			return
		}
		mid := len(data) / 2
		h := &IntHeap{}
		for _, c := range data[:mid] {
			Push(h, int(c))
		}
		for i := mid; i+1 < len(data); i += 2 {
			if h.Len() == 0 {
				break
			}
			idx := int(data[i]) % h.Len()
			(*h)[idx] = int(data[i+1])
			Fix(h, idx)
			if !Verify(h) {
				t.Fatalf("Fix(%d) broke heap", idx)
			}
		}
	})
}

// sortedOracle is a simple sorted-slice reference implementation for fuzz testing.
type sortedOracle []int

func (s *sortedOracle) push(v int) {
	data := *s
	i := sort.SearchInts(data, v)
	data = append(data, 0)
	copy(data[i+1:], data[i:])
	data[i] = v
	*s = data
}

func (s *sortedOracle) popMin() int {
	data := *s
	v := data[0]
	*s = data[1:]
	return v
}

func (s *sortedOracle) popMax() int {
	data := *s
	v := data[len(data)-1]
	*s = data[:len(data)-1]
	return v
}

func (s *sortedOracle) len() int {
	return len(*s)
}

// randIntHeapWithDups builds a random heap of n elements. The fraction
// parameter controls how many duplicate values are inserted (e.g. 0.1
// means ~10% extra duplicates).
func randIntHeapWithDups(t *testing.T, n int, fraction float64) *IntHeap {
	t.Helper()

	s := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
	r := []int{}
	for i := 0; i < n; i++ {
		r = append(r, i+1)
		for s.Float64() < fraction {
			r = append(r, i)
		}
	}

	s.Shuffle(len(r), func(i, j int) { r[i], r[j] = r[j], r[i] })

	h := &IntHeap{}
	for _, q := range r {
		Push(h, q)
	}

	if !Verify(h) {
		panic(fmt.Sprintf("not a heap!: %v", *h))
	}

	return h
}
