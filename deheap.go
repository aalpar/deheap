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
//The above copyright notice and this permission notice shall be included
// in all copies or substantial portions of the Software.
//
//	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
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

func sort3(h heap.Interface, min bool, i, j, k int) {
	l := h.Len()
	if j < l && k < l && min == h.Less(k, j) {
		h.Swap(k, j)
	}
	if j < l && i < l && min == h.Less(j, i) {
		h.Swap(j, i)
	}
	if j < l && k < l && min == h.Less(k, j) {
		h.Swap(k, j)
	}
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

func min2(h heap.Interface, min bool, i int) int {
	if i+1 >= h.Len() {
		return i
	}
	if min != h.Less(i+1, i) {
		return i
	}
	return i + 1
}

func bubbledown(h heap.Interface, min bool, i int) {
	l := h.Len()
	for {
		j := min2(h, min, hlchild(i))
		if j >= l {
			break
		}
		if min == h.Less(j, i) {
			h.Swap(j, i)
		}
		j = min4(h, min, lchild(i))
		if j >= l {
			break
		}
		if min == h.Less(j, i) {
			h.Swap(j, i)
		}
		i = j
	}
}

func bubbleup(h heap.Interface, min bool, i int) {
	if i < 0 {
		return
	}
	j := parent(i)
	for j >= 0 && min == h.Less(i, j) {
		h.Swap(i, j)
		i = j
		j = parent(i)
	}
	min = !min
	j = hparent(i)
	for j >= 0 && min == h.Less(i, j) {
		h.Swap(i, j)
		i = j
		j = parent(i)
	}
}

// Pop the smallest value off the heap.  See heap.Pop().
// Time complexity is O(log_2 n), where n = h.Len()
func Pop(h heap.Interface) interface{} {
	h.Swap(0, h.Len()-1)
	q := h.Pop()
	bubbledown(h, true, 0)
	return q
}

// Pop the largest value off the heap.  See heap.Pop().
// Time complexity is O(log_2 n), where n = h.Len()
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

// Push an element onto the heap.  See heap.Push()
// Time complexity is O(log_2 n), where n = h.Len()
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
