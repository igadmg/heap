// Package heap provides a generic min or max heap that can be used as a
// priority queue. The backing data structure is a slice denoting an implicit
// binary heap. The default value of Heap is a valid empty heap.
//
// The heap grows in the usual way by appending elements to the backing slice
// and then making adjustments to preserve the heap property. As the heap
// shrinks, smaller backing slices are periodically allocated and elements
// copied over. Empty heaps are guaranteed to be backed by nil slices.
//
// If your type can be compared using the < operator then you can use the Push
// and Pop functions to manipulate the heap, e.g.:
//
//     var myMaxIntHeap heap.Heap[int, heap.Max]
//     heap.Push(&myMaxIntHeap, 17)
//     heap.Pop(&myMaxIntHeap)
//
// If your type can't be compared using < then you can implement the Cmp method
// of the Orderable interface for your type (with a pointer receiver) and use
// the PushOrderable and PopOrderable functions:
//
//     var heap heap.Heap[myCustomType, heap.Min]
//     heap.PushOrderable(&heap, myCustomType{Key: 1, Foo: "foo"})
//     heap.PopOrderable(&heap)
//
//     type myCustomType struct {
//       Key int
//       Foo string
//     }
//
//     func (a *myCustomType) Cmp(b *myCustomType) int {
//       return x.Key - y.Key
//     }
package heap

import (
	"github.com/savsgio/gotils/nocopy"
	c "golang.org/x/exp/constraints"
)

// Heap is a min or max heap backed by a slice denoting an implicit binary heap.
// Heap is marked as noCopy because the built-in copying operation creates a
// shallow copy of the underlying slice, which is likely to give rise to
// confusing and undesired behavior.
type Heap[T any, MOM MinOrMax] struct {
	sl []T
	nocopy.NoCopy
}

// The MinOrMax interface has two implementations (Min and Max) that can be
// passed as type parameters to Heap to choose between a min and max heap.
type MinOrMax interface {
	mul() int
}

// Pass this type as the second parameter of Heap to specify a min heap
type Min struct{}

// Pass this type as the second parameter of Heap to specify a max heap
type Max struct{}

func (Min) mul() int {
	return 1
}

func (Max) mul() int {
	return -1
}

// If your type doesn't satisfy constraints.Ordered, define a Cmp method with
// a pointer receiver for your type. This method should return 0 if the two
// values compare equal, an int < 0 if the first value is less than the second,
// and an int > 0 otherwise.
//
// Example:
//
//     type Date struct {
//       Year  int
//       Month int
//       Day   int
//     }
//
//     func (m1 *Month) Cmp(m2 *Month) string {
//       if m1.Year != m2.Year {
//         return m1.Year - m2.Year
//       }
//       if m1.Month != m2.Month {
//         return m1.Month - m2.Month
//       }
//       return m1.Day - m2.Day
//     }
//
type Orderable[R any] interface {
	Cmp(*R) int
	*R
}

// Len returns the number of elements in the heap.
func Len[T any, MOM MinOrMax](heap *Heap[T, MOM]) int {
	return len(heap.sl)
}

// Push adds an element to the heap for a T that satisfies constraints.Ordered.
func Push[T c.Ordered, MOM MinOrMax](heap *Heap[T, MOM], elem T) {
	push(heap, elem, func(i, j int) int { return cmpOrdered(heap.sl[i], heap.sl[j]) })
}

// PushOrderable adds an element to the heap for a T that implements Orderable.
func PushOrderable[T any, MOM MinOrMax, PT Orderable[T]](heap *Heap[T, MOM], elem T) {
	push(heap, elem, func(i, j int) int {
		return PT(&heap.sl[i]).Cmp(&heap.sl[j])
	})
}

func push[T any, MOM MinOrMax](heap *Heap[T, MOM], elem T, cmp func(i, j int) int) {
	heap.sl = append(heap.sl, elem)
	bubble(heap, len(heap.sl)-1, cmp)
}

// Pop removes the min/max element from the heap for a T that satisfies
// constraints.Ordered.
func Pop[T c.Ordered, MOM MinOrMax](heap *Heap[T, MOM]) (T, bool) {
	return pop(heap, func(i, j int) int { return cmpOrdered(heap.sl[i], heap.sl[j]) })
}

// Pop removes the min/max element from the heap for a T that implements
// Orderable.
func PopOrderable[T any, MOM MinOrMax, PT Orderable[T]](heap *Heap[T, MOM]) (T, bool) {
	return pop(heap, func(i, j int) int {
		return PT(&heap.sl[i]).Cmp(&heap.sl[j])
	})
}

func pop[T any, MOM MinOrMax](heap *Heap[T, MOM], cmp func(i, j int) int) (val T, ok bool) {
	// This differs from (and should be superior to) the classical implementation
	// which begins by swapping the last item with the root.
	// https://www.cs.princeton.edu/courses/archive/spr09/cos423/Lectures/i-heaps.pdf

	if len(heap.sl) == 0 {
		return
	}

	ok = true
	val = heap.sl[0]

	i := pushRootHoleDownToLeaf(heap, cmp)

	if i+1 == len(heap.sl) {
		heap.sl = shrink(heap.sl)
		return
	}

	displaced := heap.sl[len(heap.sl)-1]
	heap.sl = shrink(heap.sl)
	heap.sl[i] = displaced
	bubble(heap, i, cmp)

	return
}

// Peek returns the min/max element from the min/max heap without removing it.
func Peek[T any, MOM MinOrMax](heap *Heap[T, MOM]) (val T, ok bool) {
	if len(heap.sl) == 0 {
		return
	}
	ok = true
	val = heap.sl[0]
	return
}

// Clear empties the heap.
func Clear[T any, MOM MinOrMax](heap *Heap[T, MOM]) {
	heap.sl = nil
}

// Copy performs a deep copy of the heap
func Copy[T any, MOM MinOrMax](heap *Heap[T, MOM]) Heap[T, MOM] {
	a := make([]T, len(heap.sl))
	copy(a, heap.sl)
	return Heap[T, MOM]{sl: a}
}

// A BreakOrContinue value can be returned by an iteration callback to indicate
// whether or not iteration should continue.
type BreakOrContinue int

const (
	Break    BreakOrContinue = iota
	Continue BreakOrContinue = iota
)

// Filter iterates through the elements of the heap in the order given by the
// underlying slice. If the first return value of f is false then the relevant
// element is removed from the heap. If the second return value of f is Break
// then the iteration stops without visiting any subsequent items.
func Filter[T c.Ordered, MOM MinOrMax](heap *Heap[T, MOM], f func(*T) (keepElement bool, breakOrContinue BreakOrContinue)) {
	filter(heap, f, func(i, j int) int { return cmpOrdered(heap.sl[i], heap.sl[j]) })
}

// As for Filter, but for the case where T cannot be compared using < and there
// is an implementation of Orderable[T].
func FilterOrderable[T any, MOM MinOrMax, PT Orderable[T]](heap *Heap[T, MOM], f func(*T) (keepElement bool, breakOrContinue BreakOrContinue)) {
	filter(heap, f, func(i, j int) int {
		return PT(&heap.sl[i]).Cmp(&heap.sl[j])
	})
}

func filter[T any, MOM MinOrMax](heap *Heap[T, MOM], f func(*T) (bool, BreakOrContinue), cmp func(int, int) int) {
	i := 0
	first := -1
	for j := 0; j < len(heap.sl); j++ {
		keep, boc := f(&heap.sl[j])
		if keep {
			heap.sl[i] = heap.sl[j]
			if first == -1 {
				first = i
			}
			i++
		}
		if boc == Break {
			break
		}
	}

	heap.sl = heap.sl[:i]

	if first != -1 {
		heap.sl = heap.sl[:first]
		for j := first; j < i; j++ {
			// x[i] is valid only if i < len(x), but x[i:i+1] is valid if i < cap(x),
			// so we can use this syntax to index elements outside the slice's length
			push(heap, heap.sl[j : j+1][0], cmp)
		}
	}
}

func shrink[T any](a []T) []T {
	if len(a) == 1 {
		// when the heap becomes empty again, ensure that it reverts to a nil
		// backing slice without an associated heap allocation
		return nil
	}
	a = a[0 : len(a)-1]
	if cap(a)/2 >= len(a) {
		na := make([]T, len(a))
		copy(na, a)
		return na
	}
	return a
}

func pushRootHoleDownToLeaf[T any, MOM MinOrMax](heap *Heap[T, MOM], cmp func(i, j int) int) int {
	var mom MOM

	i := 0
	for {
		lci := leftChildIndex(i)
		rci := rightChildIndex(i)
		if lci >= len(heap.sl) {
			break
		}

		// prefer to go down to the right if we can, as the tree may be shallower
		// there
		if rci >= len(heap.sl) || mom.mul()*cmp(rci, lci) > 0 {
			heap.sl[i] = heap.sl[lci]
			i = lci
		} else {
			heap.sl[i] = heap.sl[rci]
			i = rci
		}
	}
	return i
}

func bubble[T any, MOM MinOrMax](heap *Heap[T, MOM], i int, cmp func(i, j int) int) {
	var mom MOM

	for i > 0 {
		pi := parentIndex(i)
		if mom.mul()*cmp(i, pi) >= 0 {
			break
		}
		heap.sl[i], heap.sl[pi] = heap.sl[pi], heap.sl[i]
		i = pi
	}
}

func cmpOrdered[T c.Ordered](a, b T) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

func parentIndex(i int) int {
	return (i - 1) / 2
}

func leftChildIndex(i int) int {
	return (i * 2) + 1
}

func rightChildIndex(i int) int {
	return (i * 2) + 2
}
