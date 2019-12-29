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
	x := min2(h, h.Len(), true, 0)
	if x != 1 {
		t.Fatalf("unexpected value")
	}

	h = &IntHeap{1, 10}
	x = min2(h,  h.Len(),true, 0)
	if x != 0 {
		t.Fatalf("unexpected value")
	}

	x = min2(h,  h.Len(),false, 9)
	if x != 9 {
		t.Fatalf("unexpected value")
	}

	h = &IntHeap{10, 10}
	x = min2(h,  h.Len(),true, 0)
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
	if _, _, ok := isHeap(t, h); !ok {
		t.Fatalf("unexpected value: %v", h)
	}

	h = &IntHeap{1, 15, 14, 2, 3, 4, 5, 13, 12, 11, 10, 6, 7, 8, 9}
	bubbleup(h, isMinHeap(14), 14)
	if _, _, ok := isHeap(t, h); !ok {
		t.Fatalf("unexpected value: %v", h)
	}

}

func TestBubbleDown(t *testing.T) {

	h := &IntHeap{15, 1, 2}
	bubbledown(h,  h.Len(),isMinHeap(0), 0)
	if !reflect.DeepEqual(h, &IntHeap{1, 15, 2}) {
		t.Fatalf("unexpected value: %v", h)
	}

	h = &IntHeap{5, 7, 4, 6, 1, 3, 2}
	bubbledown(h,  h.Len(),isMinHeap(0), 0)
	if !reflect.DeepEqual(h, &IntHeap{1, 7, 4, 6, 5, 3, 2}) {
		t.Fatalf("unexpected value: %v", h)
	}

	h = &IntHeap{10, 8, 12, 1, 2, 9, 10, 5, 3, 4, 6, 11}
	bubbledown(h,  h.Len(),isMinHeap(0), 0)
	if !reflect.DeepEqual(h, &IntHeap{1, 10, 12, 3, 2, 9, 10, 5, 8, 4, 6, 11}) {
		t.Fatalf("unexpected value: %v", h)
	}

	h = &IntHeap{14, 15, 12, 4, 2, 3, 5, 13}
	bubbledown(h,  h.Len(),isMinHeap(0), 0)
	if _, _, ok := isHeap(t, h); !ok {
		t.Fatalf("unexpected value: %v", h)
	}

	h = &IntHeap{13, 14, 15, 3, 4, 5, 6, 7, 8, 9, 10}
	bubbledown(h,  h.Len(),isMinHeap(0), 0)
	if _, _, ok := isHeap(t, h); !ok {
		t.Fatalf("unexpected value: %v", h)
	}

}

func TestMin4(t *testing.T) {
	h := &IntHeap{3, 1, 2, 4}
	x := min4(h,  h.Len(),true, 0)
	if x != 1 {
		t.Fatalf("unexpected value")
	}

	h = &IntHeap{1, 3, 2, 4}
	x = min4(h,  h.Len(),true, 0)
	if x != 0 {
		t.Fatalf("unexpected value")
	}

	h = &IntHeap{2, 3, 1, 4}
	x = min4(h, h.Len(), true, 0)
	if x != 2 {
		t.Fatalf("unexpected value")
	}

	h = &IntHeap{2, 3, 4, 1}
	x = min4(h,  h.Len(),true, 0)
	if x != 3 {
		t.Fatalf("unexpected value")
	}

	h = &IntHeap{1, 1, 2, 2}
	x = min4(h, h.Len(), true, 0)
	if x != 0 {
		t.Fatalf("unexpected value")
	}

	h = &IntHeap{2, 2, 1, 1}
	x = min4(h,  h.Len(),true, 0)
	if x != 2 {
		t.Fatalf("unexpected value")
	}

	h = &IntHeap{2}
	x = min4(h,  h.Len(),true, 0)
	if x != 0 {
		t.Fatalf("unexpected value")
	}
}

func TestMin3(t *testing.T) {
	h := &IntHeap{3, 1, 2}
	x := min3(h,  h.Len(),true, 0, 1, 2)
	if x != 1 {
		t.Fatalf("unexpected value")
	}

	h = &IntHeap{1, 3, 2}
	x = min3(h,  h.Len(),true, 0, 1, 2)
	if x != 0 {
		t.Fatalf("unexpected value")
	}

	h = &IntHeap{2, 3, 1}
	x = min3(h,  h.Len(),true, 0, 1, 2)
	if x != 2 {
		t.Fatalf("unexpected value")
	}

	h = &IntHeap{1, 1, 2}
	x = min3(h,  h.Len(),true, 0, 1, 2)
	if x != 0 && x != 1 {
		t.Fatalf("unexpected value")
	}

	h = &IntHeap{2, 1, 1}
	x = min3(h,  h.Len(),true, 0, 1, 2)
	if x != 2 && x != 1 {
		t.Fatalf("unexpected value")
	}

	h = &IntHeap{2}
	x = min3(h,  h.Len(),true, 0, 1, 2)
	if x != 0 {
		t.Fatalf("unexpected value")
	}
}

func TestInit(t *testing.T) {
	h := &IntHeap{15, 1, 2, 14, 13, 12, 11, 3, 4, 5, 6, 7, 8, 9, 10}
	Init(h)
	if x, y, ok := isHeap(t, h); !ok {
		t.Fatalf("unexpected value: %v %v %v", x, y, h)
	}
}

func TestPush(t *testing.T) {

	h := &IntHeap{}
	for i := 0; i < 32; i++ {
		Push(h, i)
	}
	if x, y, ok := isHeap(t, h); !ok {
		t.Fatalf("unexpected value: %v %v %v", x, y, h)
	}

	h = &IntHeap{}
	for i := 3; i >= 0; i-- {
		Push(h, i)
		if x, y, ok := isHeap(t, h); !ok {
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
	if x, y, ok := isHeap(t, h); !ok {
		t.Fatalf("unexpected value: %v %v %v", x, y, h)
	}

}

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

		if a, b, ok := isHeap(t, &tv.h); !ok {
			t.Fatalf("unexpected value: %d %d %v", a, b, tv.h)
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

func TestRemove5(t *testing.T) {

	h := &IntHeap{1, 4, 4, 2, 3, 3}
	x0 := Pop(h).(int)
	if x0 != 1 {
		t.Fatalf("unexpected value")
	}
	if x0, y, ok := isHeap(t, h); !ok {
		t.Fatalf("unexpected value: %d %d %2d", x0, y, h)
	}
	x1 := Remove(h, 1).(int)
	if x1 != 4 {
		t.Fatalf("unexpected value: %d", x1)
	}
	x3 := PopMax(h).(int)
	if x3 != 4 {
		t.Fatalf("unexpected value: %d", x3)
	}
	x4 := PopMax(h).(int)
	if x4 != 3 {
		t.Fatalf("unexpected value: %d", x4)
	}
	x5 := PopMax(h).(int)
	if x5 != 3 {
		t.Fatalf("unexpected value: %d", x5)
	}
	x6 := PopMax(h).(int)
	if x6 != 2 {
		t.Fatalf("unexpected value: %d", x6)
	}
}

func TestPop3(t *testing.T) {
	h := &IntHeap{1, 6, 4, 5, 2, 3}
	x0 := Pop(h).(int)
	if x0 != 1 {
		t.Fatalf("unexpected value")
	}
	if x0, y, ok := isHeap(t, h); !ok {
		t.Fatalf("unexpected value: %d %d %2d", x0, y, h)
	}
	x1 := Pop(h).(int)
	if x1 != 2 {
		t.Fatalf("unexpected value")
	}
	x3 := Pop(h).(int)
	if x3 != 3 {
		t.Fatalf("unexpected value")
	}
	x4 := Pop(h).(int)
	if x4 != 4 {
		t.Fatalf("unexpected value")
	}
	x5 := Pop(h).(int)
	if x5 != 5 {
		t.Fatalf("unexpected value")
	}
	x6 := Pop(h).(int)
	if x6 != 6 {
		t.Fatalf("unexpected value")
	}
}

func TestPop(t *testing.T) {

	h := &IntHeap{1, 31, 30, 4, 3, 2, 5, 22, 17, 19, 21, 23, 25, 27, 29, 15, 10, 7, 16, 8, 18, 9, 20, 6, 14, 11, 24, 12, 26, 13, 28}
	x0 := Pop(h).(int)
	if x0, y, ok := isHeap(t, h); !ok {
		t.Fatalf("unexpected value: %d %d %2d", x0, y, h)
	}
	x1 := Pop(h).(int)
	if x1 < x0 {
		t.Fatalf("unexpected value: %2d", h)
	}

	h = &IntHeap{17, 31, 30, 20, 19, 22, 18, 24, 26, 23, 21, 28, 25, 27, 29}
	x0 = Pop(h).(int)
	if x0, y, ok := isHeap(t, h); !ok {
		t.Fatalf("unexpected value: %d %d %2d", x0, y, h)
	}
	x1 = Pop(h).(int)
	if x1 < x0 {
		t.Fatalf("unexpected value: %2d", h)
	}

}

func TestRandomPop(t *testing.T) {

	N := 10000
	h := randIntHeap(t, N)

	i0 := 0
	i1 := N + 1

	s := _newRand()

	for h.Len() > 0 {
		if s.Intn(2) == 0 {
			x := Pop(h).(int)
			if x < i0 {
				t.Fatalf("unexpected value")
			}
			i0 = x
		} else {
			x := PopMax(h).(int)
			if x > i1 {
				t.Fatalf("unexpected value")
			}
			i1 = x
		}
	}

}

func TestRemove(t *testing.T) {

	h := &IntHeap{0, 9, 5, 6, 1, 2, 4, 8, 7, 3}

	x := Remove(h, 9).(int)
	if x != 3 {
		t.Fatalf("unexpected value")
	}
	if _, _, ok := isHeap(t, h); !ok {
		t.Fatalf("unexpected value")
	}
	if !reflect.DeepEqual(h, &IntHeap{0, 9, 5, 6, 1, 2, 4, 8, 7}) {
		t.Fatalf("unexpected value")
	}

	x = Remove(h, 2).(int)
	if x != 5 {
		t.Fatalf("unexpected value")
	}
	if _, _, ok := isHeap(t, h); !ok {
		t.Fatalf("unexpected value")
	}
	if !reflect.DeepEqual(h, &IntHeap{0, 9, 7, 6, 1, 2, 4, 8}) {
		t.Fatalf("unexpected value")
	}

	x = Remove(h, 0).(int)
	if x != 0 {
		t.Fatalf("unexpected value")
	}
	if _, _, ok := isHeap(t, h); !ok {
		t.Fatalf("unexpected value")
	}
	if !reflect.DeepEqual(h, &IntHeap{1, 9, 7, 6, 8, 2, 4}) && !reflect.DeepEqual(h, &IntHeap{1, 9, 8, 6, 7, 2, 4}) {
		t.Fatalf("unexpected value: %v", h)
	}

}

func TestDups(t *testing.T) {

	h := &IntHeap{}

	in := []int{0, 10, 9, 8, 7, 6, 5, 5, 4, 7, 3, 2, 1}
	out := []int{0, 1, 2, 3, 4, 5, 5, 6, 7, 7, 8, 9, 10}

	for _, v := range in {
		Push(h, v)
	}

	l := h.Len()
	if l != 13 {
		t.Fatalf("unexpected value")
	}

	for i, v := range out {
		x := Pop(h).(int)
		if h.Len() != l-(i+1) {
			t.Fatalf("unexpected value")
		}
		if x != v {
			t.Fatalf("unexpected value")
		}
	}

	h = &IntHeap{}

	out = []int{10, 9, 8, 7, 7, 6, 5, 5, 4, 3, 2, 1, 0}

	for _, v := range in {
		Push(h, v)
	}

	l = h.Len()
	if l != 13 {
		t.Fatalf("unexpected value")
	}

	for i, v := range out {
		x := PopMax(h).(int)
		if h.Len() != l-(i+1) {
			t.Fatalf("unexpected value")
		}
		if x != v {
			t.Fatalf("unexpected value")
		}
	}

}

func TestOps(t *testing.T) {

	for k := 0; k < 1000; k++ {

		s := rand.New(rand.NewSource(time.Now().Unix()))

		N := s.Intn(64) + 2

		h := randIntHeapWithDups(t, N, 1/10)
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
			if _, _, ok := isHeap(t, h); !ok {
				t.Fatalf("unexpected value")
			}
		}

		h = randIntHeapWithDups(t, N, 1/10)
		for h.Len() > 0 {
			h0 := make([]int, h.Len())
			copy(h0, []int(*h))
			x := s.Intn(h.Len())
			Remove(h, x)
			if i, j, ok := isHeap(t, h); !ok {
				t.Fatalf("unexpected value: %d %d %d\n%v\n%v", x, i, j, h0, h)
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

func isHeap(t *testing.T, h heap.Interface) (int, int, bool) {
	t.Helper()
	l := h.Len()
	for i := l - 1; i >= 0; i-- {
		min := isMinHeap(i)
		p0 := parent(i)
		p1 := hparent(i)
		if p0 >= 0 && min != h.Less(p0, i) {
			return p0, i, false
		}
		if p1 >= 0 && min == h.Less(p1, i) {
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

func randIntHeap(t *testing.T, n int) *IntHeap {
	t.Helper()
	r := []int{}
	for i := 0; i < n; i++ {
		r = append(r, i+1)
	}

	s := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
	s.Shuffle(len(r), func(i, j int) { r[i], r[j] = r[j], r[i] })

	h := &IntHeap{}
	for _, q := range r {
		Push(h, q)
	}

	if x, y, ok := isHeap(t, h); !ok {
		panic(fmt.Sprintf("not a heap!: %d %d", x, y))
	}

	return h
}

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

	if x, y, ok := isHeap(t, h); !ok {
		panic(fmt.Sprintf("not a heap!: %d %d", x, y))
	}

	return h
}
