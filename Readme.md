# Heap

A generic implementation of min and max heaps in Go with an interface suitable
for use as a priority queue.

* Use with types that satisfy `constraints.Comparable`, or define a
  `Cmp` method for your type.
* Choose min or max heap property via a type parameter.
* Extensive tests (including fuzz tests).
* Benchmarks confirm O(1) push and O(log n) pop.

## What makes this heap implementation different?

1. It uses generics.

2. Unlike other generic Go heaps that I've seen, the ordering function is
  obtained from an interface implementation rather than a constructor argument.
  This has advantages and disadvantages. The main advantage is to guarantee
  that a heap of Ts always uses the 'correct' ordering function for T. The main
  disadvantages are (i) more complex types (though you don't really have to
  think about these as a consumer of the library) and (ii) the need to define
  dummy wrapper types if you want to use different ordering functions for the
  same underlying type.

## Example

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
	// maxPriorityTask == Task{Priority: 10, Payload:  "A priority 10 task"}
	fmt.Printf("%v: %+v\n", ok, maxPriorityTask)
}
```

## Docs

https://pkg.go.dev/github.com/addrummond/heap
