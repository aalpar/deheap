//
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
//

package deheap

import (
	"math"
	"math/rand"
	"reflect"
	"sort"
	"testing"
	"time"
)

// TestOrderedEmptyPop verifies that Pop and PopMax on an empty heap return
// the zero value without panicking.
func TestOrderedEmptyPop(t *testing.T) {
	h := New[int]()
	if v := h.Pop(); v != 0 {
		t.Fatalf("Pop on empty int heap = %d, want 0", v)
	}
	if v := h.PopMax(); v != 0 {
		t.Fatalf("PopMax on empty int heap = %d, want 0", v)
	}

	hs := New[string]()
	if v := hs.Pop(); v != "" {
		t.Fatalf("Pop on empty string heap = %q, want \"\"", v)
	}
	if v := hs.PopMax(); v != "" {
		t.Fatalf("PopMax on empty string heap = %q, want \"\"", v)
	}
}

// TestOrderedSingleElement verifies Pop and PopMax on a single-element heap.
func TestOrderedSingleElement(t *testing.T) {
	h := From(42)
	if v := h.Pop(); v != 42 {
		t.Fatalf("Pop single = %d, want 42", v)
	}
	if h.Len() != 0 {
		t.Fatalf("Len after Pop = %d, want 0", h.Len())
	}

	h = From(42)
	if v := h.PopMax(); v != 42 {
		t.Fatalf("PopMax single = %d, want 42", v)
	}
	if h.Len() != 0 {
		t.Fatalf("Len after PopMax = %d, want 0", h.Len())
	}
}

// TestOrderedTwoElements verifies Pop and PopMax on a two-element heap,
// including two equal elements.
func TestOrderedTwoElements(t *testing.T) {
	h := From(3, 7)
	if v := h.Pop(); v != 3 {
		t.Fatalf("Pop first = %d, want 3", v)
	}
	if v := h.Pop(); v != 7 {
		t.Fatalf("Pop second = %d, want 7", v)
	}

	h = From(3, 7)
	if v := h.PopMax(); v != 7 {
		t.Fatalf("PopMax first = %d, want 7", v)
	}
	if v := h.PopMax(); v != 3 {
		t.Fatalf("PopMax second = %d, want 3", v)
	}

	h = From(5, 5)
	if v := h.Pop(); v != 5 {
		t.Fatalf("Pop equal = %d, want 5", v)
	}
	if v := h.PopMax(); v != 5 {
		t.Fatalf("PopMax equal = %d, want 5", v)
	}
}

// TestOrderedPushPop verifies basic Pop ordering: push descending values,
// pop all ascending.
func TestOrderedPushPop(t *testing.T) {
	h := New[int]()
	for i := 32; i >= 0; i-- {
		h.Push(i)
	}
	if !h.Verify() {
		t.Fatalf("Verify() failed after Push")
	}
	prev := -1
	for h.Len() > 0 {
		v := h.Pop()
		if v < prev {
			t.Fatalf("Pop returned %d after %d", v, prev)
		}
		prev = v
	}
}

// TestOrderedPopMax verifies basic PopMax ordering: push ascending values,
// pop all descending.
func TestOrderedPopMax(t *testing.T) {
	h := New[int]()
	for i := 0; i < 33; i++ {
		h.Push(i)
	}
	if !h.Verify() {
		t.Fatalf("Verify() failed after Push")
	}
	prev := math.MaxInt
	for h.Len() > 0 {
		v := h.PopMax()
		if v > prev {
			t.Fatalf("PopMax returned %d after %d", v, prev)
		}
		prev = v
	}
}

// TestOrderedFrom verifies the From() constructor initializes a valid heap
// and drains correctly via Pop.
func TestOrderedFrom(t *testing.T) {
	h := From(5, 3, 8, 1, 9, 2, 7, 4, 6)
	if !h.Verify() {
		t.Fatalf("Verify() failed after From")
	}
	want := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	for i, w := range want {
		v := h.Pop()
		if v != w {
			t.Fatalf("Pop[%d] = %d, want %d", i, v, w)
		}
	}
}

// TestOrderedFromEmpty verifies From() with no arguments produces an
// empty heap.
func TestOrderedFromEmpty(t *testing.T) {
	h := From[int]()
	if h.Len() != 0 {
		t.Fatalf("Len = %d, want 0", h.Len())
	}
}

// TestOrderedPeek verifies Peek returns the smallest element without
// removing it.
func TestOrderedPeek(t *testing.T) {
	h := From(5, 3, 1, 4, 2)
	if v := h.Peek(); v != 1 {
		t.Fatalf("Peek = %d, want 1", v)
	}
	if h.Len() != 5 {
		t.Fatalf("Peek modified length")
	}
}

// TestOrderedPeekMax verifies PeekMax returns the largest element without
// removing it. Tests 5 elements, single element, and two elements (boundary
// cases for the if/else branches in PeekMax).
func TestOrderedPeekMax(t *testing.T) {
	h := From(5, 3, 1, 4, 2)
	if v := h.PeekMax(); v != 5 {
		t.Fatalf("PeekMax = %d, want 5", v)
	}
	if h.Len() != 5 {
		t.Fatalf("PeekMax modified length")
	}

	h = From(42)
	if v := h.PeekMax(); v != 42 {
		t.Fatalf("PeekMax single = %d, want 42", v)
	}

	h = From(1, 5)
	if v := h.PeekMax(); v != 5 {
		t.Fatalf("PeekMax two = %d, want 5", v)
	}
}

// TestOrderedRemove tests Remove at every index of a 9-element heap,
// verifying the removed value and that the remaining elements drain
// in sorted order.
func TestOrderedRemove(t *testing.T) {
	items := []int{1, 9, 5, 6, 2, 3, 4, 8, 7}
	for k := 0; k < len(items); k++ {
		h := From(items...)
		removed := h.Remove(k)
		if !h.Verify() {
			t.Fatalf("Verify() failed after Remove(%d)", k)
		}
		// Collect remaining via Pop
		var got []int
		for h.Len() > 0 {
			got = append(got, h.Pop())
		}
		// Verify sorted and removed element is missing
		expected := make([]int, len(items))
		copy(expected, items)
		sort.Ints(expected)
		j := sort.SearchInts(expected, removed)
		expected = append(expected[:j], expected[j+1:]...)
		if len(got) != len(expected) {
			t.Fatalf("Remove(%d): length mismatch", k)
		}
		for i := range got {
			if got[i] != expected[i] {
				t.Fatalf("Remove(%d): got[%d]=%d, want %d", k, i, got[i], expected[i])
			}
		}
	}
}

// TestOrderedRandomPopPopMax is the comprehensive randomized test. Over
// 1000 iterations with random sizes (2-65) and values including duplicates
// (via Intn(n)), it randomly interleaves Pop and PopMax and verifies
// monotonicity of returned values.
func TestOrderedRandomPopPopMax(t *testing.T) {
	s := rand.New(rand.NewSource(time.Now().UnixNano()))

	for iter := 0; iter < 1000; iter++ {
		n := s.Intn(64) + 2
		h := New[int]()
		for i := 0; i < n; i++ {
			h.Push(s.Intn(n))
		}

		lo := math.MinInt32
		hi := math.MaxInt32
		for h.Len() > 0 {
			if s.Intn(2) == 0 {
				v := h.Pop()
				if v < lo {
					t.Fatalf("Pop returned %d < previous min %d", v, lo)
				}
				lo = v
			} else {
				v := h.PopMax()
				if v > hi {
					t.Fatalf("PopMax returned %d > previous max %d", v, hi)
				}
				hi = v
			}
			if !h.Verify() {
				t.Fatalf("iter %d: Verify() failed", iter)
			}
		}
	}
}

// TestOrderedRandomRemove is a randomized stress test for Remove. Over 1000
// iterations, it removes elements at random indices and validates each
// removed value against a sorted-slice oracle.
func TestOrderedRandomRemove(t *testing.T) {
	s := rand.New(rand.NewSource(time.Now().UnixNano()))

	for iter := 0; iter < 1000; iter++ {
		n := s.Intn(64) + 2

		// Build an oracle sorted slice and a deheap with the same elements
		oracle := make([]int, n)
		h := New[int]()
		for i := 0; i < n; i++ {
			v := s.Intn(n)
			oracle[i] = v
			h.Push(v)
		}
		sort.Ints(oracle)

		// Remove random elements, maintaining the oracle in sync
		for h.Len() > 0 {
			idx := s.Intn(h.Len())
			removed := h.Remove(idx)
			if !h.Verify() {
				t.Fatalf("iter %d: Verify() failed after Remove(%d)", iter, idx)
			}
			j := sort.SearchInts(oracle, removed)
			oracle = append(oracle[:j], oracle[j+1:]...)
		}

		if len(oracle) != 0 {
			t.Fatalf("oracle not empty after draining heap")
		}
	}
}

// FuzzOrderedPushPop interprets a byte sequence as heap commands:
//
//	'<' = Pop, '>' = PopMax,
//	'{' + next byte = PushPop(byte), '}' + next byte = PushPopMax(byte),
//	anything else = Push(byte).
//
// Validates every result against a sorted-slice oracle.
func FuzzOrderedPushPop(f *testing.F) {
	f.Add([]byte{10, 5, 3, 8, '<', '>', 1, '<', 7, '>'})
	f.Add([]byte{1, 1, 1, '<', '<', '<'})
	f.Add([]byte{'{', 5, '}', 3, '<', '>'})
	f.Add([]byte{})

	f.Fuzz(func(t *testing.T, data []byte) {
		h := New[int]()
		var oracle []int // kept sorted

		for i := 0; i < len(data); i++ {
			c := data[i]
			switch c {
			case '<':
				if h.Len() > 0 {
					got := h.Pop()
					want := oracle[0]
					oracle = oracle[1:]
					if got != want {
						t.Fatalf("Pop: got %d, want %d", got, want)
					}
				}
			case '>':
				if h.Len() > 0 {
					got := h.PopMax()
					want := oracle[len(oracle)-1]
					oracle = oracle[:len(oracle)-1]
					if got != want {
						t.Fatalf("PopMax: got %d, want %d", got, want)
					}
				}
			case '{':
				i++
				if i >= len(data) {
					continue
				}
				v := int(data[i])
				// Oracle: insert v, pop min.
				j := sort.SearchInts(oracle, v)
				oracle = append(oracle, 0)
				copy(oracle[j+1:], oracle[j:])
				oracle[j] = v
				want := oracle[0]
				oracle = oracle[1:]
				got := h.PushPop(v)
				if got != want {
					t.Fatalf("PushPop(%d): got %d, want %d", v, got, want)
				}
			case '}':
				i++
				if i >= len(data) {
					continue
				}
				v := int(data[i])
				// Oracle: insert v, pop max.
				j := sort.SearchInts(oracle, v)
				oracle = append(oracle, 0)
				copy(oracle[j+1:], oracle[j:])
				oracle[j] = v
				want := oracle[len(oracle)-1]
				oracle = oracle[:len(oracle)-1]
				got := h.PushPopMax(v)
				if got != want {
					t.Fatalf("PushPopMax(%d): got %d, want %d", v, got, want)
				}
			default:
				v := int(c)
				h.Push(v)
				j := sort.SearchInts(oracle, v)
				oracle = append(oracle, 0)
				copy(oracle[j+1:], oracle[j:])
				oracle[j] = v
			}
			if !h.Verify() {
				t.Fatalf("Verify() failed")
			}
			if h.Len() != len(oracle) {
				t.Fatalf("length mismatch: heap=%d, oracle=%d", h.Len(), len(oracle))
			}
		}
	})
}

// FuzzOrderedRemove pushes all bytes, then removes at deterministic-random
// indices, validating each removed value against a sorted oracle.
func FuzzOrderedRemove(f *testing.F) {
	f.Add([]byte{10, 5, 3, 8, 1, 7})
	f.Add([]byte{1, 1, 1, 1})
	f.Add([]byte{})

	f.Fuzz(func(t *testing.T, data []byte) {
		h := New[int]()
		var oracle []int
		for _, c := range data {
			v := int(c)
			h.Push(v)
			i := sort.SearchInts(oracle, v)
			oracle = append(oracle, 0)
			copy(oracle[i+1:], oracle[i:])
			oracle[i] = v
		}

		s := rand.New(rand.NewSource(int64(len(data))))
		for h.Len() > 0 {
			idx := s.Intn(h.Len())
			removed := h.Remove(idx)
			if !h.Verify() {
				t.Fatalf("Verify() failed after Remove")
			}
			j := sort.SearchInts(oracle, removed)
			oracle = append(oracle[:j], oracle[j+1:]...)
		}

		if len(oracle) != 0 {
			t.Fatalf("oracle not empty after draining heap")
		}
	})
}

// FuzzOrderedFix pushes the first half of bytes onto a heap, then interprets
// the second half as (index, newValue) pairs for Fix operations. Drains the
// heap and validates sorted output against an oracle.
func FuzzOrderedFix(f *testing.F) {
	f.Add([]byte{5, 3, 8, 1, 9, 2, 7, 4, 6})
	f.Add([]byte{1, 1, 1, 1, 1, 1})
	f.Add([]byte{})

	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) < 2 {
			return
		}
		mid := len(data) / 2
		h := New[int]()
		var oracle []int
		for _, c := range data[:mid] {
			v := int(c)
			h.Push(v)
			i := sort.SearchInts(oracle, v)
			oracle = append(oracle, 0)
			copy(oracle[i+1:], oracle[i:])
			oracle[i] = v
		}

		for i := mid; i+1 < len(data); i += 2 {
			if h.Len() == 0 {
				break
			}
			idx := int(data[i]) % h.Len()
			newVal := int(data[i+1])
			oldVal := h.items[idx]

			// Update oracle.
			j := sort.SearchInts(oracle, oldVal)
			oracle = append(oracle[:j], oracle[j+1:]...)
			k := sort.SearchInts(oracle, newVal)
			oracle = append(oracle, 0)
			copy(oracle[k+1:], oracle[k:])
			oracle[k] = newVal

			h.items[idx] = newVal
			h.Fix(idx)
			if !h.Verify() {
				t.Fatalf("Verify() failed after Fix")
			}
		}

		for di, want := range oracle {
			got := h.Pop()
			if got != want {
				t.Fatalf("Pop[%d] = %d, want %d", di, got, want)
			}
		}
	})
}

// TestOrderedRemoveSingleElement verifies Remove(0) on a 1-element heap
// returns the correct value and leaves the heap empty.
func TestOrderedRemoveSingleElement(t *testing.T) {
	h := From(42)
	v := h.Remove(0)
	if v != 42 {
		t.Fatalf("Remove(0) = %d, want 42", v)
	}
	if h.Len() != 0 {
		t.Fatalf("Len after Remove = %d, want 0", h.Len())
	}
}

// TestOrderedRemoveTwoElements verifies Remove on a 2-element heap at both
// indices, including equal elements.
func TestOrderedRemoveTwoElements(t *testing.T) {
	// Remove index 0 (min)
	h := From(3, 7)
	v := h.Remove(0)
	if v != 3 {
		t.Fatalf("Remove(0) = %d, want 3", v)
	}
	if r := h.Pop(); r != 7 {
		t.Fatalf("remaining = %d, want 7", r)
	}

	// Remove index 1 (max)
	h = From(3, 7)
	v = h.Remove(1)
	if v != 7 {
		t.Fatalf("Remove(1) = %d, want 7", v)
	}
	if r := h.Pop(); r != 3 {
		t.Fatalf("remaining = %d, want 3", r)
	}

	// Equal elements
	h = From(5, 5)
	v = h.Remove(0)
	if v != 5 {
		t.Fatalf("Remove(0) equal = %d, want 5", v)
	}
	if r := h.Pop(); r != 5 {
		t.Fatalf("remaining equal = %d, want 5", r)
	}

	h = From(5, 5)
	v = h.Remove(1)
	if v != 5 {
		t.Fatalf("Remove(1) equal = %d, want 5", v)
	}
	if r := h.Pop(); r != 5 {
		t.Fatalf("remaining equal = %d, want 5", r)
	}
}

// TestOrderedFix verifies Fix restores the heap property after modifying
// a single element. Covers bubbledown, bubbleup, both level types,
// root/leaf/interior, and edge cases.
func TestOrderedFix(t *testing.T) {
	drain := func(h *Deheap[int]) []int {
		var out []int
		for h.Len() > 0 {
			out = append(out, h.Pop())
		}
		return out
	}

	// Increase root (min level) — needs bubbledown.
	h := From(1, 9, 5, 4, 6, 3, 2)
	h.items[0] = 100
	h.Fix(0)
	if !h.Verify() {
		t.Fatalf("Fix root increase: Verify() failed")
	}
	got := drain(h)
	want := []int{2, 3, 4, 5, 6, 9, 100}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Fix root increase: got %v, want %v", got, want)
	}

	// Decrease min-level node to new global min — bubbleup.
	h = From(1, 9, 5, 4, 6, 3, 2)
	h.items[3] = -1
	h.Fix(3)
	if !h.Verify() {
		t.Fatalf("Fix min node decrease: Verify() failed")
	}
	if h.Peek() != -1 {
		t.Fatalf("Fix min node decrease: Peek = %d, want -1", h.Peek())
	}

	// Increase min-level leaf — bubbleup through max chain.
	h = From(1, 9, 5, 4, 6, 3, 2)
	h.items[6] = 100
	h.Fix(6)
	if !h.Verify() {
		t.Fatalf("Fix min leaf increase: Verify() failed")
	}
	if h.PeekMax() != 100 {
		t.Fatalf("Fix min leaf increase: PeekMax = %d, want 100", h.PeekMax())
	}

	// Decrease max-level node.
	h = From(1, 9, 5, 4, 6, 3, 2)
	h.items[1] = 0
	h.Fix(1)
	if !h.Verify() {
		t.Fatalf("Fix max node decrease: Verify() failed")
	}
	if h.Peek() != 0 {
		t.Fatalf("Fix max node decrease: Peek = %d, want 0", h.Peek())
	}

	// No-op: value unchanged.
	h = From(1, 9, 5, 4, 6, 3, 2)
	h.Fix(3)
	if !h.Verify() {
		t.Fatalf("Fix no-op: Verify() failed")
	}
	got = drain(h)
	want = []int{1, 2, 3, 4, 5, 6, 9}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Fix no-op: got %v, want %v", got, want)
	}

	// Single element.
	h = From(42)
	h.items[0] = 99
	h.Fix(0)
	if !h.Verify() {
		t.Fatalf("Fix single: Verify() failed")
	}
	if h.Pop() != 99 {
		t.Fatalf("Fix single: unexpected value")
	}

	// Two elements: fix min to exceed max.
	h = From(3, 7)
	h.items[0] = 10
	h.Fix(0)
	if !h.Verify() {
		t.Fatalf("Fix two (min): Verify() failed")
	}
	if h.Peek() != 7 {
		t.Fatalf("Fix two (min): Peek = %d, want 7", h.Peek())
	}
}

// TestOrderedFixRandomized modifies random elements in random heaps and
// verifies Fix restores the heap property. Validates content by draining.
func TestOrderedFixRandomized(t *testing.T) {
	s := rand.New(rand.NewSource(time.Now().UnixNano()))

	for iter := 0; iter < 1000; iter++ {
		n := s.Intn(64) + 2
		h := New[int]()
		for i := 0; i < n; i++ {
			h.Push(s.Intn(n))
		}

		// Build sorted oracle from current heap state.
		oracle := make([]int, len(h.items))
		copy(oracle, h.items)
		sort.Ints(oracle)

		// Modify a random element and fix.
		idx := s.Intn(h.Len())
		oldVal := h.items[idx]
		newVal := s.Intn(n*2) - n
		h.items[idx] = newVal
		h.Fix(idx)
		if !h.Verify() {
			t.Fatalf("iter %d: Verify() failed after Fix(%d)", iter, idx)
		}

		// Update oracle: remove old, insert new.
		j := sort.SearchInts(oracle, oldVal)
		oracle = append(oracle[:j], oracle[j+1:]...)
		k := sort.SearchInts(oracle, newVal)
		oracle = append(oracle, 0)
		copy(oracle[k+1:], oracle[k:])
		oracle[k] = newVal

		for di, want := range oracle {
			got := h.Pop()
			if got != want {
				t.Fatalf("iter %d: Pop[%d] = %d, want %d", iter, di, got, want)
			}
		}
	}
}

// TestOrderedPeekSmall verifies Peek on 1-element and 2-element heaps.
func TestOrderedPeekSmall(t *testing.T) {
	h := From(42)
	if v := h.Peek(); v != 42 {
		t.Fatalf("Peek single = %d, want 42", v)
	}

	h = From(7, 3)
	if v := h.Peek(); v != 3 {
		t.Fatalf("Peek two = %d, want 3", v)
	}

	h = From(3, 7)
	if v := h.Peek(); v != 3 {
		t.Fatalf("Peek two ordered = %d, want 3", v)
	}
}

// TestOrderedPeekEmptyPanics verifies that Peek panics on an empty heap.
func TestOrderedPeekEmptyPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Peek on empty heap did not panic")
		}
	}()
	h := New[int]()
	h.Peek()
}

// TestOrderedPeekMaxEmptyPanics verifies that PeekMax panics on an empty heap.
func TestOrderedPeekMaxEmptyPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("PeekMax on empty heap did not panic")
		}
	}()
	h := New[int]()
	h.PeekMax()
}

// TestPushThenPeek pushes a single element into an empty heap and verifies
// both Peek and PeekMax return it.
func TestPushThenPeek(t *testing.T) {
	h := New[int]()
	h.Push(99)
	if v := h.Peek(); v != 99 {
		t.Fatalf("Peek = %d, want 99", v)
	}
	if v := h.PeekMax(); v != 99 {
		t.Fatalf("PeekMax = %d, want 99", v)
	}
	if h.Len() != 1 {
		t.Fatalf("Len = %d, want 1", h.Len())
	}
}

// TestOrderedString verifies the generic API works with string types.
func TestOrderedString(t *testing.T) {
	h := From("banana", "apple", "cherry", "date")
	if v := h.Pop(); v != "apple" {
		t.Fatalf("Pop = %q, want %q", v, "apple")
	}
	if v := h.PopMax(); v != "date" {
		t.Fatalf("PopMax = %q, want %q", v, "date")
	}
	if v := h.Pop(); v != "banana" {
		t.Fatalf("Pop = %q, want %q", v, "banana")
	}
	if v := h.Pop(); v != "cherry" {
		t.Fatalf("Pop = %q, want %q", v, "cherry")
	}
}

// TestOrderedNew verifies the generic API works with float64 types.
func TestOrderedNew(t *testing.T) {
	h := New[float64]()
	h.Push(3.14)
	h.Push(2.71)
	h.Push(1.41)
	if v := h.Pop(); v != 1.41 {
		t.Fatalf("Pop = %v, want 1.41", v)
	}
	if v := h.PopMax(); v != 3.14 {
		t.Fatalf("PopMax = %v, want 3.14", v)
	}
}

// TestNewBounded verifies NewBounded creates a heap with the correct max size.
func TestNewBounded(t *testing.T) {
	h := NewBounded[int](5)
	if h.MaxLen() != 5 {
		t.Fatalf("MaxLen = %d, want 5", h.MaxLen())
	}
	if h.Len() != 0 {
		t.Fatalf("Len = %d, want 0", h.Len())
	}
}

// TestNewBoundedPanicsOnZeroOrNeg verifies NewBounded panics on non-positive maxSize.
func TestNewBoundedPanicsOnZeroOrNeg(t *testing.T) {
	for _, v := range []int{0, -1, -100} {
		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Fatalf("NewBounded(%d) did not panic", v)
				}
			}()
			NewBounded[int](v)
		}()
	}
}

// TestMaxLenUnbounded verifies MaxLen returns 0 for unbounded heaps.
func TestMaxLenUnbounded(t *testing.T) {
	h := New[int]()
	if h.MaxLen() != 0 {
		t.Fatalf("MaxLen = %d, want 0", h.MaxLen())
	}
	h2 := From(1, 2, 3)
	if h2.MaxLen() != 0 {
		t.Fatalf("MaxLen = %d, want 0", h2.MaxLen())
	}
}

// TestOrderedPushPopEmpty verifies PushPop on an empty heap returns the
// pushed element without modifying the heap.
func TestOrderedPushPopEmpty(t *testing.T) {
	h := New[int]()
	if v := h.PushPop(42); v != 42 {
		t.Fatalf("PushPop empty = %d, want 42", v)
	}
	if h.Len() != 0 {
		t.Fatalf("Len = %d, want 0", h.Len())
	}
}

// TestOrderedPushPopNewMin verifies PushPop returns the pushed element
// when it is smaller than the current minimum (no heap modification).
func TestOrderedPushPopNewMin(t *testing.T) {
	h := From(3, 7, 5)
	if v := h.PushPop(1); v != 1 {
		t.Fatalf("PushPop new min = %d, want 1", v)
	}
	if h.Len() != 3 {
		t.Fatalf("Len = %d, want 3", h.Len())
	}
	if !h.Verify() {
		t.Fatal("Verify() failed")
	}
}

// TestOrderedPushPopEqual verifies PushPop when the pushed element
// equals the current minimum (returns the pushed element).
func TestOrderedPushPopEqual(t *testing.T) {
	h := From(3, 7, 5)
	if v := h.PushPop(3); v != 3 {
		t.Fatalf("PushPop equal = %d, want 3", v)
	}
	if h.Len() != 3 {
		t.Fatalf("Len = %d, want 3", h.Len())
	}
	if !h.Verify() {
		t.Fatal("Verify() failed")
	}
}

// TestOrderedPushPopReplaces verifies PushPop replaces the root when
// the pushed element is larger than the current minimum.
func TestOrderedPushPopReplaces(t *testing.T) {
	h := From(1, 9, 5, 4, 6, 3, 2)
	v := h.PushPop(4)
	if v != 1 {
		t.Fatalf("PushPop = %d, want 1", v)
	}
	if h.Len() != 7 {
		t.Fatalf("Len = %d, want 7", h.Len())
	}
	if !h.Verify() {
		t.Fatal("Verify() failed")
	}
	var got []int
	for h.Len() > 0 {
		got = append(got, h.Pop())
	}
	want := []int{2, 3, 4, 4, 5, 6, 9}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("drain = %v, want %v", got, want)
	}
}

// TestOrderedPushPopSingle verifies PushPop on a single-element heap.
func TestOrderedPushPopSingle(t *testing.T) {
	h := From(5)
	if v := h.PushPop(3); v != 3 {
		t.Fatalf("PushPop smaller = %d, want 3", v)
	}
	if h.Peek() != 5 {
		t.Fatalf("remaining = %d, want 5", h.Peek())
	}

	h = From(5)
	if v := h.PushPop(7); v != 5 {
		t.Fatalf("PushPop larger = %d, want 5", v)
	}
	if h.Peek() != 7 {
		t.Fatalf("remaining = %d, want 7", h.Peek())
	}
}

// TestOrderedPushPopRandomized verifies PushPop against a Push+Pop oracle.
func TestOrderedPushPopRandomized(t *testing.T) {
	s := rand.New(rand.NewSource(time.Now().UnixNano()))

	for iter := 0; iter < 1000; iter++ {
		n := s.Intn(64) + 1
		h := New[int]()
		oracle := New[int]()
		for i := 0; i < n; i++ {
			v := s.Intn(n)
			h.Push(v)
			oracle.Push(v)
		}

		o := s.Intn(n * 2)
		got := h.PushPop(o)
		oracle.Push(o)
		want := oracle.Pop()
		if got != want {
			t.Fatalf("iter %d: PushPop(%d) = %d, want %d", iter, o, got, want)
		}
		if !h.Verify() {
			t.Fatalf("iter %d: Verify() failed", iter)
		}
		for h.Len() > 0 {
			hv := h.Pop()
			ov := oracle.Pop()
			if hv != ov {
				t.Fatalf("iter %d: drain mismatch %d != %d", iter, hv, ov)
			}
		}
	}
}

// TestOrderedPushPopMaxEmpty verifies PushPopMax on an empty heap returns
// the pushed element without modifying the heap.
func TestOrderedPushPopMaxEmpty(t *testing.T) {
	h := New[int]()
	if v := h.PushPopMax(42); v != 42 {
		t.Fatalf("PushPopMax empty = %d, want 42", v)
	}
	if h.Len() != 0 {
		t.Fatalf("Len = %d, want 0", h.Len())
	}
}

// TestOrderedPushPopMaxNewMax verifies PushPopMax returns the pushed element
// when it is >= the current maximum (no heap modification).
func TestOrderedPushPopMaxNewMax(t *testing.T) {
	h := From(3, 7, 5)
	if v := h.PushPopMax(10); v != 10 {
		t.Fatalf("PushPopMax new max = %d, want 10", v)
	}
	if h.Len() != 3 {
		t.Fatalf("Len = %d, want 3", h.Len())
	}
	if !h.Verify() {
		t.Fatal("Verify() failed")
	}
}

// TestOrderedPushPopMaxEqual verifies PushPopMax when pushed element
// equals the current maximum (returns it without modifying heap).
func TestOrderedPushPopMaxEqual(t *testing.T) {
	h := From(3, 7, 5)
	if v := h.PushPopMax(7); v != 7 {
		t.Fatalf("PushPopMax equal = %d, want 7", v)
	}
	if h.Len() != 3 {
		t.Fatalf("Len = %d, want 3", h.Len())
	}
	if !h.Verify() {
		t.Fatal("Verify() failed")
	}
}

// TestOrderedPushPopMaxReplaces verifies PushPopMax evicts the current max
// when the pushed element is smaller.
func TestOrderedPushPopMaxReplaces(t *testing.T) {
	h := From(1, 9, 5, 4, 6, 3, 2)
	v := h.PushPopMax(4)
	if v != 9 {
		t.Fatalf("PushPopMax = %d, want 9", v)
	}
	if h.Len() != 7 {
		t.Fatalf("Len = %d, want 7", h.Len())
	}
	if !h.Verify() {
		t.Fatal("Verify() failed")
	}
}

// TestOrderedPushPopMaxSingle verifies PushPopMax on a single-element heap.
func TestOrderedPushPopMaxSingle(t *testing.T) {
	h := From(5)
	if v := h.PushPopMax(7); v != 7 {
		t.Fatalf("PushPopMax larger = %d, want 7", v)
	}
	if h.Peek() != 5 {
		t.Fatalf("remaining = %d, want 5", h.Peek())
	}

	h = From(5)
	if v := h.PushPopMax(3); v != 5 {
		t.Fatalf("PushPopMax smaller = %d, want 5", v)
	}
	if h.Peek() != 3 {
		t.Fatalf("remaining = %d, want 3", h.Peek())
	}
}

// TestOrderedPushPopMaxTwo verifies PushPopMax on a two-element heap.
func TestOrderedPushPopMaxTwo(t *testing.T) {
	h := From(3, 7)
	if v := h.PushPopMax(5); v != 7 {
		t.Fatalf("PushPopMax = %d, want 7", v)
	}
	if !h.Verify() {
		t.Fatal("Verify() failed")
	}

	h = From(3, 7)
	if v := h.PushPopMax(10); v != 10 {
		t.Fatalf("PushPopMax larger = %d, want 10", v)
	}
	if h.Len() != 2 {
		t.Fatalf("Len = %d, want 2", h.Len())
	}

	h = From(3, 7)
	if v := h.PushPopMax(1); v != 7 {
		t.Fatalf("PushPopMax smaller = %d, want 7", v)
	}
	if h.Peek() != 1 {
		t.Fatalf("Peek = %d, want 1", h.Peek())
	}
}

// TestOrderedPushPopMaxRandomized verifies PushPopMax against a Push+PopMax oracle.
func TestOrderedPushPopMaxRandomized(t *testing.T) {
	s := rand.New(rand.NewSource(time.Now().UnixNano()))

	for iter := 0; iter < 1000; iter++ {
		n := s.Intn(64) + 1
		h := New[int]()
		oracle := New[int]()
		for i := 0; i < n; i++ {
			v := s.Intn(n)
			h.Push(v)
			oracle.Push(v)
		}

		o := s.Intn(n * 2)
		got := h.PushPopMax(o)
		oracle.Push(o)
		want := oracle.PopMax()
		if got != want {
			t.Fatalf("iter %d: PushPopMax(%d) = %d, want %d", iter, o, got, want)
		}
		if !h.Verify() {
			t.Fatalf("iter %d: Verify() failed", iter)
		}
		for h.Len() > 0 {
			hv := h.Pop()
			ov := oracle.Pop()
			if hv != ov {
				t.Fatalf("iter %d: drain mismatch %d != %d", iter, hv, ov)
			}
		}
	}
}

// TestDrainAscEmpty verifies DrainAsc on an empty heap yields nothing.
func TestDrainAscEmpty(t *testing.T) {
	h := New[int]()
	count := 0
	for range h.DrainAsc() {
		count++
	}
	if count != 0 {
		t.Fatalf("DrainAsc empty yielded %d elements", count)
	}
}

// TestDrainAsc verifies DrainAsc yields elements in ascending order.
func TestDrainAsc(t *testing.T) {
	h := From(5, 1, 9, 3, 7)
	var got []int
	for v := range h.DrainAsc() {
		got = append(got, v)
	}
	want := []int{1, 3, 5, 7, 9}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("DrainAsc = %v, want %v", got, want)
	}
	if h.Len() != 0 {
		t.Fatalf("Len after drain = %d, want 0", h.Len())
	}
}

// TestDrainDescEmpty verifies DrainDesc on an empty heap yields nothing.
func TestDrainDescEmpty(t *testing.T) {
	h := New[int]()
	count := 0
	for range h.DrainDesc() {
		count++
	}
	if count != 0 {
		t.Fatalf("DrainDesc empty yielded %d elements", count)
	}
}

// TestDrainDesc verifies DrainDesc yields elements in descending order.
func TestDrainDesc(t *testing.T) {
	h := From(5, 1, 9, 3, 7)
	var got []int
	for v := range h.DrainDesc() {
		got = append(got, v)
	}
	want := []int{9, 7, 5, 3, 1}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("DrainDesc = %v, want %v", got, want)
	}
	if h.Len() != 0 {
		t.Fatalf("Len after drain = %d, want 0", h.Len())
	}
}

// TestDrainAscEarlyBreak verifies early termination leaves the heap valid.
func TestDrainAscEarlyBreak(t *testing.T) {
	h := From(5, 1, 9, 3, 7)
	var got []int
	for v := range h.DrainAsc() {
		if v > 3 {
			break
		}
		got = append(got, v)
	}
	want := []int{1, 3}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("DrainAsc early = %v, want %v", got, want)
	}
	if h.Len() != 2 {
		t.Fatalf("remaining Len = %d, want 2", h.Len())
	}
	if !h.Verify() {
		t.Fatal("Verify() failed after early break")
	}
}

// TestDrainDescEarlyBreak verifies early termination of DrainDesc.
func TestDrainDescEarlyBreak(t *testing.T) {
	h := From(5, 1, 9, 3, 7)
	var got []int
	for v := range h.DrainDesc() {
		if v < 7 {
			break
		}
		got = append(got, v)
	}
	want := []int{9, 7}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("DrainDesc early = %v, want %v", got, want)
	}
	if h.Len() != 2 {
		t.Fatalf("remaining Len = %d, want 2", h.Len())
	}
	if !h.Verify() {
		t.Fatal("Verify() failed after early break")
	}
}

// TestDrainAscRandomized compares DrainAsc output to sort.Ints oracle.
func TestDrainAscRandomized(t *testing.T) {
	s := rand.New(rand.NewSource(time.Now().UnixNano()))
	for iter := 0; iter < 1000; iter++ {
		n := s.Intn(64) + 1
		items := make([]int, n)
		for i := range items {
			items[i] = s.Intn(n)
		}
		h := From(items...)
		want := make([]int, n)
		copy(want, items)
		sort.Ints(want)
		var got []int
		for v := range h.DrainAsc() {
			got = append(got, v)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("iter %d: DrainAsc = %v, want %v", iter, got, want)
		}
	}
}

// TestDrainDescRandomized compares DrainDesc output to reverse-sorted oracle.
func TestDrainDescRandomized(t *testing.T) {
	s := rand.New(rand.NewSource(time.Now().UnixNano()))
	for iter := 0; iter < 1000; iter++ {
		n := s.Intn(64) + 1
		items := make([]int, n)
		for i := range items {
			items[i] = s.Intn(n)
		}
		h := From(items...)
		want := make([]int, n)
		copy(want, items)
		sort.Sort(sort.Reverse(sort.IntSlice(want)))
		var got []int
		for v := range h.DrainDesc() {
			got = append(got, v)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("iter %d: DrainDesc = %v, want %v", iter, got, want)
		}
	}
}

// TestOfferUnbounded verifies Offer on an unbounded heap behaves like Push.
func TestOfferUnbounded(t *testing.T) {
	h := New[int]()
	evicted, didEvict := h.Offer(5)
	if didEvict {
		t.Fatal("Offer on unbounded heap evicted")
	}
	var zero int
	if evicted != zero {
		t.Fatalf("evicted = %v, want zero", evicted)
	}
	if h.Len() != 1 {
		t.Fatalf("Len = %d, want 1", h.Len())
	}
}

// TestOfferUnderCapacity verifies Offer when bounded heap is not full.
func TestOfferUnderCapacity(t *testing.T) {
	h := NewBounded[int](5)
	for i := 0; i < 4; i++ {
		_, didEvict := h.Offer(i)
		if didEvict {
			t.Fatalf("Offer(%d) evicted when under capacity", i)
		}
	}
	if h.Len() != 4 {
		t.Fatalf("Len = %d, want 4", h.Len())
	}
	if !h.Verify() {
		t.Fatal("Verify() failed")
	}
}

// TestOfferAtCapacityReject verifies Offer rejects elements >= current max.
func TestOfferAtCapacityReject(t *testing.T) {
	h := NewBounded[int](3)
	h.Push(1)
	h.Push(5)
	h.Push(3)

	evicted, didEvict := h.Offer(5)
	if !didEvict {
		t.Fatal("Offer(5) should evict")
	}
	if evicted != 5 {
		t.Fatalf("evicted = %d, want 5 (rejected)", evicted)
	}
	if h.Len() != 3 {
		t.Fatalf("Len = %d, want 3", h.Len())
	}

	evicted, didEvict = h.Offer(100)
	if !didEvict {
		t.Fatal("Offer(100) should evict")
	}
	if evicted != 100 {
		t.Fatalf("evicted = %d, want 100 (rejected)", evicted)
	}
}

// TestOfferAtCapacityEvict verifies Offer evicts the max and inserts o.
func TestOfferAtCapacityEvict(t *testing.T) {
	h := NewBounded[int](3)
	h.Push(1)
	h.Push(5)
	h.Push(3)

	evicted, didEvict := h.Offer(2)
	if !didEvict {
		t.Fatal("Offer(2) should evict")
	}
	if evicted != 5 {
		t.Fatalf("evicted = %d, want 5", evicted)
	}
	if h.Len() != 3 {
		t.Fatalf("Len = %d, want 3", h.Len())
	}
	if !h.Verify() {
		t.Fatal("Verify() failed")
	}
	if h.PeekMax() != 3 {
		t.Fatalf("PeekMax = %d, want 3", h.PeekMax())
	}
}

// TestOfferSingleCapacity verifies Offer on a bounded heap of size 1.
func TestOfferSingleCapacity(t *testing.T) {
	h := NewBounded[int](1)
	_, didEvict := h.Offer(5)
	if didEvict {
		t.Fatal("first Offer should not evict")
	}

	evicted, didEvict := h.Offer(3)
	if !didEvict {
		t.Fatal("second Offer should evict")
	}
	if evicted != 5 {
		t.Fatalf("evicted = %d, want 5", evicted)
	}
	if h.Peek() != 3 {
		t.Fatalf("Peek = %d, want 3", h.Peek())
	}

	evicted, didEvict = h.Offer(10)
	if !didEvict {
		t.Fatal("Offer(10) should evict")
	}
	if evicted != 10 {
		t.Fatalf("evicted = %d, want 10 (rejected)", evicted)
	}
}

// TestFromBoundedUnderCapacity verifies FromBounded with fewer items than maxSize.
func TestFromBoundedUnderCapacity(t *testing.T) {
	h := FromBounded[int](5, 1, 3, 2)
	if h.Len() != 3 {
		t.Fatalf("Len = %d, want 3", h.Len())
	}
	if h.MaxLen() != 5 {
		t.Fatalf("MaxLen = %d, want 5", h.MaxLen())
	}
	if !h.Verify() {
		t.Fatal("Verify() failed")
	}
}

// TestFromBoundedExactCapacity verifies FromBounded with exactly maxSize items.
func TestFromBoundedExactCapacity(t *testing.T) {
	h := FromBounded[int](3, 5, 1, 3)
	if h.Len() != 3 {
		t.Fatalf("Len = %d, want 3", h.Len())
	}
	if !h.Verify() {
		t.Fatal("Verify() failed")
	}
}

// TestFromBoundedOverCapacity verifies FromBounded keeps the N smallest elements.
func TestFromBoundedOverCapacity(t *testing.T) {
	h := FromBounded[int](3, 9, 1, 7, 3, 5)
	if h.Len() != 3 {
		t.Fatalf("Len = %d, want 3", h.Len())
	}
	if !h.Verify() {
		t.Fatal("Verify() failed")
	}
	// Should contain the 3 smallest: 1, 3, 5
	var got []int
	for v := range h.DrainAsc() {
		got = append(got, v)
	}
	want := []int{1, 3, 5}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("FromBounded drain = %v, want %v", got, want)
	}
}

// TestFromBoundedEmpty verifies FromBounded with no items.
func TestFromBoundedEmpty(t *testing.T) {
	h := FromBounded[int](5)
	if h.Len() != 0 {
		t.Fatalf("Len = %d, want 0", h.Len())
	}
	if h.MaxLen() != 5 {
		t.Fatalf("MaxLen = %d, want 5", h.MaxLen())
	}
}

// TestFromBoundedPanicsOnNonPositive verifies FromBounded panics on maxSize <= 0.
func TestFromBoundedPanicsOnNonPositive(t *testing.T) {
	for _, v := range []int{0, -1, -100} {
		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Fatalf("FromBounded(%d) did not panic", v)
				}
			}()
			FromBounded[int](v)
		}()
	}
}

// TestOfferRandomized stress-tests Offer with a sorted-oracle.
// Builds a bounded heap and streams n values through it, maintaining
// a sorted-slice oracle of the maxSize smallest values seen.
func TestOfferRandomized(t *testing.T) {
	s := rand.New(rand.NewSource(time.Now().UnixNano()))

	for iter := 0; iter < 1000; iter++ {
		maxSize := s.Intn(30) + 1
		n := s.Intn(100) + maxSize
		h := NewBounded[int](maxSize)
		var oracle []int

		for i := 0; i < n; i++ {
			v := s.Intn(n)
			h.Offer(v)
			// Insert v into sorted oracle, keep only maxSize smallest.
			j := sort.SearchInts(oracle, v)
			oracle = append(oracle, 0)
			copy(oracle[j+1:], oracle[j:])
			oracle[j] = v
			if len(oracle) > maxSize {
				oracle = oracle[:maxSize]
			}
		}

		if h.Len() != len(oracle) {
			t.Fatalf("iter %d: Len = %d, want %d", iter, h.Len(), len(oracle))
		}
		if !h.Verify() {
			t.Fatalf("iter %d: Verify() failed", iter)
		}
		for i, want := range oracle {
			got := h.Pop()
			if got != want {
				t.Fatalf("iter %d: Pop[%d] = %d, want %d", iter, i, got, want)
			}
		}
	}
}

func BenchmarkOrderedPush(b *testing.B) {
	r := make([]int, b.N)
	for i := range r {
		r[i] = i
	}
	s := rand.New(rand.NewSource(time.Now().UnixNano()))
	s.Shuffle(len(r), func(i, j int) { r[i], r[j] = r[j], r[i] })

	b.ResetTimer()

	h := New[int]()
	for _, q := range r {
		h.Push(q)
	}
}

func BenchmarkOrderedPop(b *testing.B) {
	r := make([]int, b.N)
	for i := range r {
		r[i] = i
	}
	s := rand.New(rand.NewSource(time.Now().UnixNano()))
	s.Shuffle(len(r), func(i, j int) { r[i], r[j] = r[j], r[i] })

	h := New[int]()
	for _, q := range r {
		h.Push(q)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		h.Pop()
	}
}

func BenchmarkOrderedPopMax(b *testing.B) {
	r := make([]int, b.N)
	for i := range r {
		r[i] = i
	}
	s := rand.New(rand.NewSource(time.Now().UnixNano()))
	s.Shuffle(len(r), func(i, j int) { r[i], r[j] = r[j], r[i] })

	h := New[int]()
	for _, q := range r {
		h.Push(q)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		h.PopMax()
	}
}

func BenchmarkOrderedPushPop(b *testing.B) {
	r := make([]int, b.N)
	for i := range r {
		r[i] = i
	}
	s := rand.New(rand.NewSource(time.Now().UnixNano()))
	s.Shuffle(len(r), func(i, j int) { r[i], r[j] = r[j], r[i] })

	b.ResetTimer()

	h := New[int]()
	for _, q := range r {
		h.Push(q)
	}
	for i := 0; i < b.N; i++ {
		h.Pop()
	}
}
