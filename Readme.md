# Heap

A generic implementation of min and max binary heaps in Go with an interface
suitable for use as a priority queue.

* Use with types that satisfy `constraints.Comparable`, or define a
  `Cmp` method for your type.
* Choose min or max heap property via a type parameter.
* Extensive tests (including fuzz tests).
* Benchmarks confirm O(1) push and O(log n) pop.

## Docs

https://pkg.go.dev/github.com/addrummond/heap

## What makes this heap implementation different?

* It uses generics.
* Unlike other generic Go heaps that I've seen, the ordering function is
  obtained from an interface implementation rather than a constructor argument.
  This has advantages and disadvantages.

  Advantages:
    - All heaps of Ts are guaranteed to use the same ordering function.
    - Empty heaps consume only the space required by a `nil` slice (as it's not
      necessary to store the ordering function as a field of the `Heap` struct).

  Disadvantages:
    - The types are more complex (though you don't really have to
      think about these as a consumer of the library).
    - You need to define dummy wrapper types if you want different heaps to use
      different ordering functions for the same underlying type.

## Example with a built-in type that can be compared using <

```go
package main

import (
	"fmt"

	"github.com/addrummond/heap"
)

func main() {
	var h heap.Heap[int, heap.Max]

	heap.Push(&h, 5)
	heap.Push(&h, 10)
	heap.Push(&h, 1)

	maxVal, ok := heap.Pop(&h)
	// ok == true
	// maxVal == 10
	fmt.Printf("%v: %+v\n", ok, maxVal)
}
```

## Example with a custom data type

```go
package main

import (
	"fmt"

	"github.com/addrummond/heap"
)

type Task struct {
	Priority int
	Payload  any
}

// As Task is a user-defined datatype that doesn't satisfy constraints.Ordered,
// we need to implement the heap.Orderable interface, which has a single method,
// Cmp.
func (t1 *Task) Cmp(t2 *Task) int {
	return t1.Priority - t2.Priority
}

func main() {
	var h heap.Heap[Task, heap.Max]

	heap.PushOrderable(&h, Task{
		Priority: 5,
		Payload:  "A priority 5 task",
	})

	heap.PushOrderable(&h, Task{
		Priority: 10,
		Payload:  "A priority 10 task",
	})

	heap.PushOrderable(&h, Task{
		Priority: 1,
		Payload:  "A priority 1 task",
	})

	maxPriorityTask, ok := heap.PopOrderable(&h)
	// ok == true
	// maxPriorityTask == Task{Priority: 10, Payload: "A priority 10 task"}
	fmt.Printf("%v: %+v\n", ok, maxPriorityTask)
}
```
