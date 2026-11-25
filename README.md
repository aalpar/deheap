# deheap

[![Go Reference](https://pkg.go.dev/badge/github.com/aalpar/deheap.svg)](https://pkg.go.dev/github.com/aalpar/deheap)
[![Go Report Card](https://goreportcard.com/badge/github.com/aalpar/deheap)](https://goreportcard.com/report/github.com/aalpar/deheap)
[![codecov](https://codecov.io/gh/aalpar/deheap/graph/badge.svg)](https://codecov.io/gh/aalpar/deheap)

A doubly-ended heap (min-max heap) for Go. Provides O(log n) access to both
the minimum and maximum elements of a collection through a single data
structure, with zero external dependencies.

```
go get github.com/aalpar/deheap
```

## Why a doubly-ended heap?

A standard binary heap gives you efficient access to one extremum — the
smallest or the largest element — but not both. Many practical problems need
both ends simultaneously:

**Scheduling and resource allocation.** Operating system schedulers and job
queues routinely need the highest-priority task (to run next) and the
lowest-priority task (to evict or age). A doubly-ended priority queue avoids
maintaining two separate heaps and the bookkeeping to keep them synchronized.

**Bounded-size caches and buffers.** When a priority queue has a capacity
limit, insertions must discard the least valuable element. With a min-max
heap, both the insertion (against the max) and the eviction (from the min, or
vice versa) are logarithmic — no linear scan required.

**Median maintenance and order statistics.** Streaming median algorithms
typically partition data into a max-heap of the lower half and a min-heap of
the upper half. A single min-max heap can serve double duty, simplifying the
implementation.

**Network packet scheduling.** Rate-controlled and deadline-aware packet
schedulers (e.g., in QoS systems) need to dequeue by earliest deadline and
drop by lowest priority, both efficiently.

**Search algorithms.** Algorithms like SMA\* (Simplified Memory-Bounded A\*)
maintain an open set where the node with the lowest f-cost is expanded next
and the node with the highest f-cost is pruned when memory is exhausted.

## API

`deheap` provides two API surfaces.

### Generic API (Go 1.23+)

For `cmp.Ordered` types — `int`, `float64`, `string`, and friends — use the
type-safe generic API directly:

```go
import "github.com/aalpar/deheap"

// Construct from existing elements.
h := deheap.From(5, 3, 8, 1, 9)

// Or build incrementally.
h := deheap.New[int]()
h.Push(5)
h.Push(3)

// O(1) access to both extrema.
fmt.Println(h.Peek())    // smallest
fmt.Println(h.PeekMax()) // largest

// O(log n) removal from either end.
min := h.Pop()
max := h.PopMax()

// Remove by index.
val := h.Remove(2)

// Update an element and restore heap order.
h.Push(10)
h.Push(20)
// ... mutate the element at index 0 ...
h.Fix(0)

// Check if the heap is valid.
fmt.Println(h.Verify()) // true

// Push and pop in one step; cheaper than Push followed by Pop or PopMax.
smallest := h.PushPop(4)   // smaller of 4 and the current minimum
largest := h.PushPopMax(7) // larger of 7 and the current maximum

// Bounded heap: Offer evicts the largest element when full.
// Push ignores the bound.
b := deheap.NewBounded[int](3)
for _, v := range []int{5, 1, 9, 3} {
    if evicted, ok := b.Offer(v); ok {
        fmt.Println("evicted", evicted) // evicted 9
    }
}

// Drain in sorted order (or DrainDesc). Breaking out early leaves the
// remaining elements in a valid heap.
for v := range h.DrainAsc() {
    fmt.Println(v)
}
```

### Interface API

For custom types, implement `heap.Interface` and use the package-level
functions. This is the original v1 API and remains stable.

```go
import "github.com/aalpar/deheap"

type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *IntHeap) Push(x interface{}) { *h = append(*h, x.(int)) }

func (h *IntHeap) Pop() interface{} {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[:n-1]
    return x
}

h := &IntHeap{2, 1, 5, 6}
deheap.Init(h)
deheap.Push(h, 3)
min := deheap.Pop(h)
max := deheap.PopMax(h)

// Push and pop in one step; cheaper than Push followed by Pop or PopMax.
min = deheap.PushPop(h, 4)
max = deheap.PushPopMax(h, 7)

// Update an element and restore heap order.
(*h)[0] = 42
deheap.Fix(h, 0)

// Check if the heap is valid.
fmt.Println(deheap.Verify(h)) // true
```

## Implementation

The underlying data structure is a **min-max heap** stored in a flat slice
with no pointers and no additional bookkeeping. Nodes on even levels
(0, 2, 4, ...) satisfy the min-heap property with respect to their
descendants, and nodes on odd levels (1, 3, 5, ...) satisfy the max-heap
property. The root is always the minimum; the maximum is one of its two
children.

Level parity is determined by bit-length of the 1-based index, computed via
`math/bits.Len` — a single CPU instruction on most architectures. Insertions
bubble up through grandparent links; deletions bubble down through
grandchild links, with a secondary swap against the binary-tree parent when
the element crosses a level boundary.

The generic API implements the same algorithms directly on `[]T` with
native `<` comparisons, eliminating interface dispatch, adapter allocation,
and boxing overhead.

### Complexity

| Operation | Time     | Space |
|-----------|----------|-------|
| `Push`    | O(log n) | O(1) amortized |
| `Pop`     | O(log n) | O(1) |
| `PopMax`  | O(log n) | O(1) |
| `PushPop`, `PushPopMax` | O(log n) | O(1) |
| `Offer`   | O(log n) | O(1) amortized |
| `Remove`  | O(log n) | O(1) |
| `Fix`     | O(log n) | O(1) |
| `Peek`    | O(1)     | O(1) |
| `PeekMax` | O(1)     | O(1) |
| `DrainAsc`, `DrainDesc` | O(n log n) | O(1) |
| `Init`    | O(n)     | O(1) |
| `Verify`  | O(n)     | O(1) |

Storage is a single contiguous slice — one element per slot, no child
pointers, no color bits, no auxiliary arrays. Memory overhead beyond the
elements themselves is the slice header (24 bytes on 64-bit systems).

### Benchmarks

Measured on Apple M4 Max (arm64), Go 1.27; each figure is the median of
five runs.

**Generic API** (`Deheap[T]`, direct `<` comparisons, no interface dispatch):

| Operation   | ns/op  | B/op | allocs/op |
|-------------|--------|------|-----------|
| Push        | 12     | 46   | 0         |
| Pop         | 190    | 0    | 0         |
| PopMax      | 187    | 0    | 0         |

**Interface API** (v1, `heap.Interface`):

| Operation   | ns/op  | B/op | allocs/op |
|-------------|--------|------|-----------|
| Push        | 20     | 51   | 0         |
| Pop         | 260    | 7    | 0         |
| PopMax      | 255    | 7    | 0         |

**Standard library** (`container/heap`) for comparison:

| Operation   | ns/op  | B/op | allocs/op |
|-------------|--------|------|-----------|
| Push        | 21     | 53   | 0         |
| Pop         | 189    | 7    | 0         |

The generic API is faster than `container/heap` on Push and on par with it on
Pop, despite examining up to four grandchildren per node (vs two children in a
binary heap), while also providing O(log n) access to the maximum. The v1
interface API pays the same `heap.Interface` dispatch cost as
`container/heap`. Generic API operations are zero-allocation in steady state.

#### Benchmark descriptions

The benchmarks compare deheap's two API surfaces against each other and
against `container/heap`. Baseline measurements isolate the cost of slice
operations from heap logic. Unless noted, benchmarks use a heap of 10,000
`int` elements to ensure the working set fits in L1 cache while exercising
the full tree depth.

Sample run (median run of five per benchmark):

```
$ go test -bench . -benchmem
goos: darwin
goarch: arm64
pkg: github.com/aalpar/deheap
cpu: Apple M4 Max
BenchmarkMin4-16                        	262670930	         4.548 ns/op	       0 B/op	       0 allocs/op
BenchmarkBaselinePush-16                	796840660	         2.456 ns/op	      42 B/op	       0 allocs/op
BenchmarkPush-16                        	54124554	        19.70 ns/op	      51 B/op	       0 allocs/op
BenchmarkPop-16                         	 6313336	       260.3 ns/op	       7 B/op	       0 allocs/op
BenchmarkPopMax-16                      	 6505016	       255.0 ns/op	       7 B/op	       0 allocs/op
BenchmarkPushPop-16                     	 5705757	       270.3 ns/op	      60 B/op	       1 allocs/op
BenchmarkHeapPushPop-16                 	 6924955	       217.0 ns/op	      61 B/op	       1 allocs/op
BenchmarkHeapPop-16                     	 9926960	       188.6 ns/op	       7 B/op	       0 allocs/op
BenchmarkHeapPush-16                    	51842851	        21.10 ns/op	      53 B/op	       0 allocs/op
BenchmarkV1PushPop-16                   	59951672	        20.36 ns/op	      15 B/op	       1 allocs/op
BenchmarkV1PushPopMax-16                	55412994	        21.38 ns/op	      15 B/op	       1 allocs/op
BenchmarkOrderedPush-16                 	99061327	        11.58 ns/op	      46 B/op	       0 allocs/op
BenchmarkOrderedPop-16                  	 9352300	       190.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkOrderedPopMax-16               	 9542306	       187.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkOrderedPushPop-16              	 8619138	       205.0 ns/op	      45 B/op	       0 allocs/op
BenchmarkOrderedGenericPushPop-16       	359725708	         3.294 ns/op	       0 B/op	       0 allocs/op
BenchmarkOrderedGenericPushPopMax-16    	297904705	         3.889 ns/op	       0 B/op	       0 allocs/op
BenchmarkOrderedDrainAsc-16             	    2322	    536634 ns/op	       0 B/op	       0 allocs/op
BenchmarkOrderedDrainDesc-16            	    2605	    443729 ns/op	       0 B/op	       0 allocs/op
BenchmarkOrderedOffer-16                	265210098	         4.496 ns/op	       0 B/op	       0 allocs/op
PASS
```

| Benchmark | What it measures |
|-----------|-----------------|
| `OrderedPush` | Generic `Deheap[int].Push`: append + bubble-up with direct `<` comparisons. |
| `OrderedPop` | Generic `Deheap[int].Pop`: remove minimum and bubble-down with direct `<`. |
| `OrderedPopMax` | Generic `Deheap[int].PopMax`: remove maximum and bubble-down with direct `<`. |
| `OrderedPushPop` | Generic push-then-drain throughput. |
| `OrderedGenericPushPop` | Generic `PushPop`. Each call removes the minimum, so the heap converges to the largest values seen and nearly every input takes the early return (input ≤ minimum). Measures that fast path, not the bubble-down case. |
| `OrderedGenericPushPopMax` | Generic `PushPopMax`. Mirror of the above: the heap converges to the smallest values, so nearly every input takes the early return (input ≥ maximum). |
| `OrderedDrainAsc` | Full `DrainAsc` of a 10,000-element heap; ns/op is per drain (about 54 ns per element). |
| `OrderedDrainDesc` | Full `DrainDesc` of a 10,000-element heap; ns/op is per drain (about 44 ns per element). |
| `OrderedOffer` | `Offer` on a full bounded heap of 1,000. The heap converges to the smallest values seen, so nearly every input is rejected without modifying the heap. |
| `Min4` | Cost of `min4`, the internal function that finds the extremum among up to 4 grandchildren during bubble-down. |
| `BaselinePush` | Raw slice append with no heap ordering; establishes the floor cost of memory allocation and copying. |
| `Push` | `deheap.Push` (v1): append + bubble-up via `heap.Interface`. |
| `Pop` | `deheap.Pop` (v1): remove the minimum element and bubble-down via `heap.Interface`. |
| `PopMax` | `deheap.PopMax` (v1): remove the maximum element and bubble-down via `heap.Interface`. |
| `PushPop` | v1 push-then-drain throughput. |
| `V1PushPop` | `deheap.PushPop` (v1); same early-return fast path as `OrderedGenericPushPop`. |
| `V1PushPopMax` | `deheap.PushPopMax` (v1); same early-return fast path as `OrderedGenericPushPopMax`. |
| `HeapPushPop` | Same push-then-drain pattern using `container/heap` for direct comparison. |
| `HeapPop` | `container/heap.Pop` in isolation, comparable to `Pop` above. |
| `HeapPush` | `container/heap.Push` in isolation, comparable to `Push` above. |

## Testing

The test suite includes 90 test functions covering internal helpers,
algorithmic correctness, edge cases (empty, single-element, two-element
heaps), and large-scale randomized validation. Six native Go fuzz targets
(`testing.F`) exercise both API surfaces under arbitrary input. Tests are run
against Go 1.23 and 1.24 in CI.

```bash
go test ./...          # run all tests
go test -bench . ./... # run benchmarks
go test -fuzz .        # run fuzz tests
```

## Requirements

- Go 1.23 or later
- Zero external dependencies

## References

1. M.D. Atkinson, J.-R. Sack, N. Santoro, and T. Strothotte. "Min-Max
   Heaps and Generalized Priority Queues." *Communications of the ACM*,
   29(10):996–1000, October 1986.
   https://doi.org/10.1145/6617.6621

2. J. van Leeuwen and D. Wood. "Interval Heaps." *The Computer Journal*,
   36(3):209–216, 1993.

3. P. Brass. *Advanced Data Structures*. Cambridge University Press, 2008.

## License

MIT — see [LICENSE](LICENSE).
