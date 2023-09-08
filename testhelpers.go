package heap

import (
	"fmt"
	"sort"
	"strings"
)

// Pretty prints the heap as a tree. Used in tests.
func debugPrintHeap[T any, MOM MinOrMax](heap *Heap[T, MOM]) string {
	a := heap.sl
	var sb strings.Builder

	if len(a) == 0 {
		return ""
	}

	bhl := 1
	for bhl < len(a) {
		bhl *= 2
	}

	formatted := make([]string, bhl)

	maxLen := 0
	for i, h := range a {
		formatted[i] = fmt.Sprintf("%v", h)
		if len(formatted[i]) > maxLen {
			maxLen = len(formatted[i])
		}
	}
	maxLen += 2

	offsets := make([]int, bhl)
	var f func(i int)
	f = func(i int) {
		if i >= bhl {
			return
		}

		li := leftChildIndex(i)
		ri := rightChildIndex(i)

		if ri >= bhl {
			if i*2 >= bhl {
				offsets[i] = maxLen + offsets[i-1]
			} else {
				offsets[i] = maxLen / 2
			}
		} else {
			f(li)
			f(ri)
			offsets[i] = (offsets[li] + offsets[ri]) / 2
		}
	}

	f(0)

	startingOffset := (maxLen - len(formatted[bhl/2])) / 2
	level := 0
	off := 0
	for {
		currentOff := 0
		wspace := func() {
			if currentOff >= startingOffset {
				sb.WriteByte(' ')
			}
			currentOff++
		}

		for i := off; i < off+(1<<level); i++ {
			if i >= bhl {
				return sb.String()
			}

			for currentOff+maxLen/2 < offsets[i] {
				wspace()
			}
			lpad := (maxLen - len(formatted[i])) / 2
			for j := 0; j < lpad; j++ {
				wspace()
			}
			sb.WriteString(formatted[i])
			currentOff += len(formatted[i])
			for j := 0; i+1 < off+(1<<level) && j < maxLen-len(formatted[i])-lpad; j++ {
				wspace()
			}
		}
		sb.WriteByte('\n')

		off += (1 << level)
		level++
	}
}

// Push an element onto a slice then sort the slice in ascending order
func naiveMinHeapPush(heap *[]int, v int) {
	*heap = append(*heap, v)
	sort.Slice(*heap, func(i, j int) bool {
		return (*heap)[i] < (*heap)[j]
	})
}

// Push an element onto a slice then sort the slice in descending order
func naiveMaxHeapPush(heap *[]int, v int) {
	*heap = append(*heap, v)
	sort.Slice(*heap, func(i, j int) bool {
		return (*heap)[j] < (*heap)[i]
	})
}

// Remove the last element from the slice. As the slice is already sorted, it's
// not necessary to do anything else.
func naiveHeapPop(heap *[]int) (v int, ok bool) {
	if len(*heap) == 0 {
		return
	}
	ok = true
	v = (*heap)[0]
	*heap = (*heap)[1:]
	return
}

func naiveHeapFilter(heap *[]int, f func(*int) (keepElement bool, breakOrContinue BreakOrContinue)) {
	newHeap := make([]int, 0)
	for i, elem := range *heap {
		keep, boc := f(&(*heap)[i])
		if keep {
			newHeap = append(newHeap, elem)
		}
		if boc == Break {
			break
		}
	}
	*heap = newHeap
}

func slicesHaveSameElems(sl1 []int, sl2 []int) bool {
	counts1 := make(map[int]int)
	counts2 := make(map[int]int)
	for _, elem := range sl1 {
		counts1[elem]++
	}
	for _, elem := range sl2 {
		counts2[elem]++
	}
	for i, n := range counts1 {
		if n2, ok := counts2[i]; !ok || n2 != n {
			return false
		}
	}
	for i, n := range counts2 {
		if n2, ok := counts1[i]; !ok || n2 != n {
			return false
		}
	}
	return true
}

func checkMinHeapProperty(heap *Heap[int, Min], i int) bool {
	if i >= len(heap.sl) {
		return true
	}
	lci := leftChildIndex(i)
	rci := rightChildIndex(i)
	if (lci < len(heap.sl) && heap.sl[lci] < heap.sl[i]) || (rci < len(heap.sl) && heap.sl[rci] < heap.sl[i]) {
		return false
	}
	return checkMinHeapProperty(heap, lci) && checkMinHeapProperty(heap, rci)
}

func checkMaxHeapProperty(heap *Heap[int, Max], i int) bool {
	if i >= len(heap.sl) {
		return true
	}
	lci := leftChildIndex(i)
	rci := rightChildIndex(i)
	if (lci < len(heap.sl) && heap.sl[lci] > heap.sl[i]) || (rci < len(heap.sl) && heap.sl[rci] > heap.sl[i]) {
		return false
	}
	return checkMaxHeapProperty(heap, lci) && checkMaxHeapProperty(heap, rci)
}
