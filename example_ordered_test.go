//
// Copyright 2024 Aaron H. Alpar
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

package deheap_test

import (
	"fmt"

	"github.com/aalpar/deheap"
)

func Example_ordered() {
	h := deheap.From(2, 1, 5, 6)
	h.Push(3)
	fmt.Printf("minimum: %d\n", h.Peek())
	fmt.Printf("maximum: %d\n", h.PeekMax())
	for h.Len() > 3 {
		fmt.Printf("%d ", h.PopMax())
	}
	for h.Len() > 1 {
		fmt.Printf("%d ", h.Pop())
	}
	fmt.Printf("middle value: %d\n", h.Peek())
	// Output:
	// minimum: 1
	// maximum: 6
	// 6 5 1 2 middle value: 3
}

func ExampleDeheap_Verify() {
	h := deheap.From(2, 1, 5, 6)
	fmt.Println(h.Verify())

	h.Push(3)
	fmt.Println(h.Verify())
	// Output:
	// true
	// true
}
