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
package deheap

import (
	"container/heap"
	"math/rand"
	"reflect"
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

func TestLChild(t *testing.T) {
	x := lchild(0)
	if x != 3 {
		t.Fatalf("unexpected value")
	}
	x = lchild(3)
	if x != 15 {
		t.Fatalf("unexpected value")
	}
	x = parent(15)
	if x != 3 {
		t.Fatalf("unexpected value")
	}
	x = parent(5)
	if x != 0 {
		t.Fatalf("unexpected value: %d", x)
	}
}

func TestHLChild(t *testing.T) {
	x := hlchild(0)
	if x != 1 {
		t.Fatalf("unexpected value")
	}
	x = hlchild(1)
	if x != 3 {
		t.Fatalf("unexpected value")
	}
	x = hparent(3)
	if x != 1 {
		t.Fatalf("unexpected value")
	}
}

func TestMin2(t *testing.T) {
	h := &IntHeap{10, 1}
	x := min2(h, true, 0)
	if x != 1 {
		t.Fatalf("unexpected value")
	}

	h = &IntHeap{1, 10}
	x = min2(h, true, 0)
	if x != 0 {
		t.Fatalf("unexpected value")
	}

	x = min2(h, false, 9)
	if x != 9 {
		t.Fatalf("unexpected value")
	}

	h = &IntHeap{10, 10}
	x = min2(h, true, 0)
	if x != 0 {
		t.Fatalf("unexpected value")
	}

}

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
	if _, _, ok := isHeap(h); !ok {
		t.Fatalf("unexpected value: %v", h)
	}
}

func TestBubbleDown(t *testing.T) {
	h := &IntHeap{15, 1, 2}
	bubbledown(h, isMinHeap(0), 0)
	if !reflect.DeepEqual(h, &IntHeap{1, 15, 2}) {
		t.Fatalf("unexpected value: %v", h)
	}

	h = &IntHeap{14, 15, 12, 4, 2, 3, 5, 13}
	bubbledown(h, isMinHeap(0), 0)
	if _, _, ok := isHeap(h); !ok {
		t.Fatalf("unexpected value: %v", h)
	}

	h = &IntHeap{13, 14, 15, 3, 4, 5, 6, 7, 8, 9, 10}
	bubbledown(h, isMinHeap(0), 0)
	if _, _, ok := isHeap(h); !ok {
		t.Fatalf("unexpected value: %v", h)
	}
}

func TestMin4(t *testing.T) {
	h := &IntHeap{3, 1, 2, 4}
	x := min4(h, true, 0)
	if x != 1 {
		t.Fatalf("unexpected value")
	}

	h = &IntHeap{1, 3, 2, 4}
	x = min4(h, true, 0)
	if x != 0 {
		t.Fatalf("unexpected value")
	}

	h = &IntHeap{2, 3, 1, 4}
	x = min4(h, true, 0)
	if x != 2 {
		t.Fatalf("unexpected value")
	}

	h = &IntHeap{2, 3, 4, 1}
	x = min4(h, true, 0)
	if x != 3 {
		t.Fatalf("unexpected value")
	}

	h = &IntHeap{1, 1, 2, 2}
	x = min4(h, true, 0)
	if x != 0 {
		t.Fatalf("unexpected value")
	}

	h = &IntHeap{2, 2, 1, 1}
	x = min4(h, true, 0)
	if x != 2 {
		t.Fatalf("unexpected value")
	}

	h = &IntHeap{2}
	x = min4(h, true, 0)
	if x != 0 {
		t.Fatalf("unexpected value")
	}
}

func TestInit(t *testing.T) {
	h := &IntHeap{15, 1, 2, 14, 13, 12, 11, 3, 4, 5, 6, 7, 8, 9, 10}
	Init(h)
	if x, y, ok := isHeap(h); !ok {
		t.Fatalf("unexpected value: %v %v %v", x, y, h)
	}
}

func TestSort3(t *testing.T) {
	good := &IntHeap{1, 2, 15}

	h := &IntHeap{1, 2, 15}
	sort3(h, true, 0, 1, 2)
	if !reflect.DeepEqual(h, good) {
		t.Fatalf("unexpected value: %v", h)
	}

	h = &IntHeap{1, 15, 2}
	sort3(h, true, 0, 1, 2)
	if !reflect.DeepEqual(h, good) {
		t.Fatalf("unexpected value: %v", h)
	}

	h = &IntHeap{2, 15, 1}
	sort3(h, true, 0, 1, 2)
	if !reflect.DeepEqual(h, good) {
		t.Fatalf("unexpected value: %v", h)
	}

	h = &IntHeap{2, 1, 15}
	sort3(h, true, 0, 1, 2)
	if !reflect.DeepEqual(h, good) {
		t.Fatalf("unexpected value: %v", h)
	}

	h = &IntHeap{15, 1, 2}
	sort3(h, true, 0, 1, 2)
	if !reflect.DeepEqual(h, good) {
		t.Fatalf("unexpected value: %v", h)
	}

	h = &IntHeap{15, 2, 1}
	sort3(h, true, 0, 1, 2)
	if !reflect.DeepEqual(h, good) {
		t.Fatalf("unexpected value: %v", h)
	}
}

func TestPush(t *testing.T) {

	h := &IntHeap{}
	for i := 0; i < 32; i++ {
		Push(h, i)
	}
	if x, y, ok := isHeap(h); !ok {
		t.Fatalf("unexpected value: %v %v %v", x, y, h)
	}

	h = &IntHeap{}
	for i := 3; i >= 0; i-- {
		Push(h, i)
		if x, y, ok := isHeap(h); !ok {
			t.Fatalf("unexpected value: %v %v %v", x, y, h)
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
	if x, y, ok := isHeap(h); !ok {
		t.Fatalf("unexpected value: %v %v %v", x, y, h)
	}

}

func TestPop(t *testing.T) {

	h := &IntHeap{1, 31, 30, 4, 3, 2, 5, 22, 17, 19, 21, 23, 25, 27, 29, 15, 10, 7, 16, 8, 18, 9, 20, 6, 14, 11, 24, 12, 26, 13, 28}
	x0 := Pop(h).(int)
	if x0, y, ok := isHeap(h); !ok {
		t.Fatalf("unexpected value: %d %d %2d", x0, y, h)
	}
	x1 := Pop(h).(int)
	if x1 < x0 {
		t.Fatalf("unexpected value: %2d", h)
	}

	h = &IntHeap{17, 31, 30, 20, 19, 22, 18, 24, 26, 23, 21, 28, 25, 27, 29}
	x0 = Pop(h).(int)
	if x0, y, ok := isHeap(h); !ok {
		t.Fatalf("unexpected value: %d %d %2d", x0, y, h)
	}
	x1 = Pop(h).(int)
	if x1 < x0 {
		t.Fatalf("unexpected value: %2d", h)
	}

	h = &IntHeap{7, 99, 98, 12, 8, 10, 9, 95, 96, 87, 69, 80, 78, 97, 81, 13, 19, 20, 22, 25, 41, 17, 16, 14, 11, 15, 51, 28, 46, 44, 38, 89, 91, 94, 86, 85, 93, 88, 68, 53, 79, 83, 72, 56, 43, 60, 29, 63, 58, 64, 77, 75, 71, 66, 73, 84, 90, 92, 50, 45, 57, 74, 62, 40, 34, 18, 59, 42, 70, 47, 55, 23, 24, 39, 67, 65, 27, 61, 52, 37, 32, 33, 30, 76, 82, 54, 48, 35, 31, 21, 26, 36, 49}

	x0 = Pop(h).(int)
	if x0, y, ok := isHeap(h); !ok {
		t.Fatalf("unexpected value: %d %d %3d", x0, y, h)
	}
	x1 = Pop(h).(int)
	if x1 < x0 {
		t.Fatalf("unexpected value: %2d", h)
	}

}

func BenchmarkMin4(b *testing.B) {

	r := &[]int{}
	for i := 0; i < b.N; i++ {
		*r = append(*r, i)
	}

	s := rand.New(rand.NewSource(time.Now().Unix()))
	s.Shuffle(len(*r), func(i, j int) { (*r)[i], (*r)[j] = (*r)[j], (*r)[i] })

	h := (*IntHeap)(r)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		min4(h, true, i)
	}

}

func BenchmarkPush(b *testing.B) {

	r := []int{}
	for i := 0; i < b.N; i++ {
		r = append(r, i)
	}

	s := rand.New(rand.NewSource(time.Now().Unix()))
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

	s := rand.New(rand.NewSource(time.Now().Unix()))
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

	s := rand.New(rand.NewSource(time.Now().Unix()))
	s.Shuffle(len(r), func(i, j int) { r[i], r[j] = r[j], r[i] })

	b.ResetTimer()

	h := &IntHeap{}
	for _, q := range r {
		Push(h, q)
	}
	for i := 0; i < b.N; i++ {
		PopMax(h)
	}

}

func BenchmarkPushPop(b *testing.B) {

	r := []int{}
	for i := 0; i < b.N; i++ {
		r = append(r, i)
	}

	s := rand.New(rand.NewSource(time.Now().Unix()))
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

	s := rand.New(rand.NewSource(time.Now().Unix()))
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

	s := rand.New(rand.NewSource(time.Now().Unix()))
	s.Shuffle(len(r), func(i, j int) { r[i], r[j] = r[j], r[i] })

	h := &IntHeap{}
	for _, q := range r {
		heap.Push(h, q)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		heap.Pop(h)
	}

}

func BenchmarkHeapPush(b *testing.B) {

	r := []int{}
	for i := 0; i < b.N; i++ {
		r = append(r, i)
	}

	s := rand.New(rand.NewSource(time.Now().Unix()))
	s.Shuffle(len(r), func(i, j int) { r[i], r[j] = r[j], r[i] })

	b.ResetTimer()

	h := &IntHeap{}
	for _, q := range r {
		heap.Push(h, q)
	}

}

func isHeap(h heap.Interface) (int, int, bool) {
	l := h.Len()
	for i := l - 1; i >= 0; i-- {
		min := isMinHeap(i)
		p0 := parent(i)
		p1 := hparent(i)
		if p0 < 0 {
			continue
		}
		if min != h.Less(p0, i) {
			return p0, i, false
		}
		if min == h.Less(p1, i) {
			return p1, i, false
		}
	}
	if l > 1 && h.Less(1, 0) {
		return 1, 0, false
	}
	if l > 2 && h.Less(2, 0) {
		return 2, 0, false
	}
	return 0, 0, true
}
