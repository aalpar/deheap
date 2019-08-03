//
// Copyright 2019 Aaron H. Alpar
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

//
// Package deheap provides the implementation of a doubly ended heap.
// Doubly ended heaps are heaps with two sides, a min side and a max side.
// Like normal single-sided heaps, elements can be pushed onto and pulled
// off of a deheap.  deheaps have an additional Pop function, PopMax, that
// returns elements from the opposite side of the ordering.
//
// This implementation has emphasized compatibility with existing libraries
// in the sort and heap packages.
//
// Performace of the deheap functions should be very close to the
// performance of the functions of the heap library
//
package deheap

import (
	"container/heap"
	"math/bits"
)

func hparent(i int) int {
	return (i - 1) / 2
}

func hlchild(i int) int {
	return (i * 2) + 1
}

func parent(i int) int {
	return ((i + 1) / 4) - 1
}

func lchild(i int) int {
	return ((i + 1) * 4) - 1
}

func level(i int) int {
	return bits.Len(uint(i)+1) - 1
}

func isMinHeap(i int) bool {
	return level(i)%2 == 0
}

func min4(h heap.Interface, min bool, i int) (k int) {
	l := h.Len()
	k = i
	if i+1 >= l {
		return k
	}
	if min == h.Less(i+1, k) {
		k = i + 1
	}
	if i+2 >= l {
		return k
	}
	if min == h.Less(i+2, k) {
		k = i + 2
	}
	if i+3 >= l {
		return k
	}
	if min == h.Less(i+3, k) {
		k = i + 3
	}
	return k
}

// min2
func min2(h heap.Interface, min bool, i int) int {
	if i+1 >= h.Len() {
		return i
	}
	if min != h.Less(i+1, i) {
		return i
	}
	return i + 1
}

// min3
func min3(h heap.Interface, min bool, i, j, k int) int {
	q := i
	if j < h.Len() && h.Less(j, q) == min {
		q = j
	}
	if k < h.Len() && h.Less(k, q) == min {
		q = k
	}
	return q
}

// bubbledown
func bubbledown(h heap.Interface, min bool, i int) (q int, r int) {
	l := h.Len()
	q = i
	r = i
	for {
		// find min of children
		j := min2(h, min, hlchild(i))
		if j >= l {
			break
		}
		// find min of grandchildren
		k := min4(h, min, lchild(i))
		// swap of less than the element at i
		v := min3(h, min, i, j, k)
		if v == i || v >= l {
			break
		}
		if v == j {
			h.Swap(v, i)
			q = v
			break
		}
		// v == k
		q = v
		h.Swap(v, i)
		p := hparent(v)
		if h.Less(p, v) == min {
			h.Swap(p, v)
			r = p
		}
		i = v
	}
	return q, r
}

// bubbleup
func bubbleup(h heap.Interface, min bool, i int) (q bool) {
	if i < 0 {
		return false
	}
	j := parent(i)
	for j >= 0 && min == h.Less(i, j) {
		q = true
		h.Swap(i, j)
		i = j
		j = parent(i)
	}
	min = !min
	j = hparent(i)
	for j >= 0 && min == h.Less(i, j) {
		q = true
		h.Swap(i, j)
		i = j
		j = parent(i)
	}
	return q
}

// Pop the smallest value off the heap.  See heap.Pop().
// Time complexity is O(log n), where n = h.Len()
func Pop(h heap.Interface) interface{} {
	h.Swap(0, h.Len()-1)
	q := h.Pop()
	bubbledown(h, true, 0)
	return q
}

// Pop the largest value off the heap.  See heap.Pop().
// Time complexity is O(log n), where n = h.Len()
func PopMax(h heap.Interface) interface{} {
	l := h.Len() - 1
	j := 0
	if l > 0 {
		j = min2(h, false, 1)
	}
	h.Swap(j, l)
	q := h.Pop()
	bubbledown(h, false, j)
	return q
}

// Remove element at index i.  See heap.Remove().
// The complexity is O(log n) where n = h.Len().
func Remove(h heap.Interface, i int) (q interface{}) {
	n := h.Len() - 1
	h.Swap(i, n)
	q = h.Pop()
	if n != i {
		q, r := bubbledown(h, isMinHeap(i), i)
		bubbleup(h, isMinHeap(q), q)
		bubbleup(h, isMinHeap(r), r)
	}
	return q
}

// Push an element onto the heap.  See heap.Push()
// Time complexity is O(log n), where n = h.Len()
func Push(h heap.Interface, o interface{}) {
	h.Push(o)
	l := h.Len() - 1
	bubbleup(h, isMinHeap(l), l)
}

// Init initializes the heap.
// This should be called once on non-empty heaps before calling Pop(), PopMax() or Push().  See heap.Init()
func Init(h heap.Interface) {
	for i := 0; i < h.Len(); i++ {
		bubbleup(h, isMinHeap(i), i)
	}
}
