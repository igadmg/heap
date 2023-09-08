package heap

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"testing"
)

func TestPushAndPop1(t *testing.T) {
	elems := []int{1, 5, 2, 9, -3, 17, 18, 19, 14}
	var heap Heap[int, Min]
	for _, elem := range elems {
		Push(&heap, elem)
	}
	sort.Ints(elems)
	for i := 0; i < len(elems); i++ {
		v, ok := Pop(&heap)
		if !ok {
			t.Errorf("Expecting ok")
		}
		if v != elems[i] {
			t.Errorf("Unexpected value")
		}
	}
	if heap.sl != nil {
		t.Errorf("Expecting empty heap to have nil backing slice")
	}
}

func TestPushAndPop2(t *testing.T) {
	elems := []int{1, -3}
	var heap Heap[int, Min]
	for _, elem := range elems {
		Push(&heap, elem)
	}
	sort.Ints(elems)
	for i := 0; i < len(elems); i++ {
		v, ok := Pop(&heap)
		if !ok {
			t.Errorf("Expecting ok")
		}
		if v != elems[i] {
			t.Errorf("Unexpected value")
		}
	}
	if heap.sl != nil {
		t.Errorf("Expecting empty heap to have nil backing slice")
	}
}

func TestExpectedMinHeapLayout1(t *testing.T) {
	elems := []int{1, 5, 2, 9, -3, 17, 18, 19, 14}
	var heap Heap[int, Min]
	for _, elem := range elems {
		Push(&heap, elem)
	}
	const expected = `
              -3
      1               2
  9       5       17      18
19  14
`
	if layout := debugPrintHeap(&heap); strings.TrimSpace(expected) != strings.TrimSpace(layout) {
		t.Errorf("Unexpected heap layout:\n%v\n", layout)
	}
	l := Len(&heap)
	if l != len(elems) {
		t.Errorf("Expected heap to have length %v, got %v\n", len(elems), l)
	}
}

func TestExpectedMinHeapLayout2(t *testing.T) {
	elems := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	var heap Heap[int, Min]
	for _, elem := range elems {
		Push(&heap, elem)
	}
	const expected = `
              0
      1               2
  3       4       5       6
7   8   9   10
`
	if layout := debugPrintHeap(&heap); strings.TrimSpace(expected) != strings.TrimSpace(layout) {
		t.Errorf("Unexpected heap layout:\n%v\n", layout)
	}
	l := Len(&heap)
	if l != len(elems) {
		t.Errorf("Expected heap to have length %v, got %v\n", len(elems), l)
	}
}

func TestExpectedMinHeapLayout3(t *testing.T) {
	const len = 30
	var heap Heap[myCustomType, Min]
	for i := 0; i < len; i++ {
		PushOrderable(&heap, myCustomType{Key: 10, Content: fmt.Sprintf("%v", i)})
	}
	const expected = `
                                                                   {10 0}
                               {10 1}                                                                  {10 2}
             {10 3}                              {10 4}                              {10 5}                              {10 6}
    {10 7}            {10 8}            {10 9}            {10 10}           {10 11}           {10 12}           {10 13}           {10 14}
{10 15}  {10 16}  {10 17}  {10 18}  {10 19}  {10 20}  {10 21}  {10 22}  {10 23}  {10 24}  {10 25}  {10 26}  {10 27}  {10 28}  {10 29}
`
	if layout := debugPrintHeap(&heap); strings.TrimSpace(expected) != strings.TrimSpace(layout) {
		t.Errorf("Unexpected heap layout:\n%v\n", layout)
	}
	l := Len(&heap)
	if l != len {
		t.Errorf("Expected heap to have length %v, got %v\n", len, l)
	}
}

func TestExpectedMaxHeapLayout1(t *testing.T) {
	elems := []int{1, 5, 2, 9, -3, 17, 18, 19, 14}
	var heap Heap[int, Max]
	for _, elem := range elems {
		Push(&heap, elem)
	}
	const expected = `
              19
      18              17
  14      -3      2       9
1   5  
`
	if layout := debugPrintHeap(&heap); strings.TrimSpace(expected) != strings.TrimSpace(layout) {
		t.Errorf("Unexpected heap layout:\n%v\n", layout)
	}
	l := Len(&heap)
	if l != len(elems) {
		t.Errorf("Expected heap to have length %v, got %v\n", len(elems), l)
	}
}

func TestExpectedMaxHeapLayout2(t *testing.T) {
	elems := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	var heap Heap[int, Max]
	for _, elem := range elems {
		Push(&heap, elem)
	}
	const expected = `
              10
      9               5
  6       8       1       4
0   3   2   7
`
	if layout := debugPrintHeap(&heap); strings.TrimSpace(expected) != strings.TrimSpace(layout) {
		t.Errorf("Unexpected heap layout:\n%v\n", layout)
	}
	l := Len(&heap)
	if l != len(elems) {
		t.Errorf("Expected heap to have length %v, got %v\n", len(elems), l)
	}
}

func TestExpectedMaxHeapLayout3(t *testing.T) {
	const len = 30
	var heap Heap[myCustomType, Max]
	for i := 0; i < len; i++ {
		PushOrderable(&heap, myCustomType{Key: 10, Content: fmt.Sprintf("%v", i)})
	}
	const expected = `
                                                                   {10 0}
                               {10 1}                                                                  {10 2}
             {10 3}                              {10 4}                              {10 5}                              {10 6}
    {10 7}            {10 8}            {10 9}            {10 10}           {10 11}           {10 12}           {10 13}           {10 14}
{10 15}  {10 16}  {10 17}  {10 18}  {10 19}  {10 20}  {10 21}  {10 22}  {10 23}  {10 24}  {10 25}  {10 26}  {10 27}  {10 28}  {10 29}
`
	if layout := debugPrintHeap(&heap); strings.TrimSpace(expected) != strings.TrimSpace(layout) {
		t.Errorf("Unexpected heap layout:\n%v\n", layout)
	}
	l := Len(&heap)
	if l != len {
		t.Errorf("Expected heap to have length %v, got %v\n", len, l)
	}
}

func TestLenEmpty(t *testing.T) {
	var heap Heap[int, Min]
	l := Len(&heap)
	if l != 0 {
		t.Errorf("Empty heap should have length 0, got %v\n", l)
	}
}

func TestPopEmpty(t *testing.T) {
	var heap Heap[int, Min]
	_, ok := Pop(&heap)
	if ok {
		t.Errorf("Calling Pop on an empty heap [1] should have returned ok=false")
	}
	Push(&heap, 1)
	Pop(&heap)
	_, ok = Pop(&heap)
	if ok {
		t.Errorf("Calling Pop on an empty heap [1] should have returned ok=false")
	}
}

func TestPopMin(t *testing.T) {
	elems := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	var heap Heap[int, Min]
	for _, elem := range elems {
		Push(&heap, elem)
	}
	v1, ok1 := Pop(&heap)
	v2, ok2 := Pop(&heap)
	if !ok1 || !ok2 || v1 != 0 || v2 != 1 {
		t.Errorf("Unepxected pop results: %v %v, %v %v\n", v1, ok1, v2, ok2)
	}
	const expected = `
              2
      3               5
  7       4       9       6
10  8
`
	if layout := debugPrintHeap(&heap); strings.TrimSpace(expected) != strings.TrimSpace(layout) {
		t.Errorf("Unexpected heap layout:\n%v\n", layout)
	}
}

func TestPopMax(t *testing.T) {
	elems := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	var heap Heap[int, Max]
	for _, elem := range elems {
		Push(&heap, elem)
	}
	v1, ok1 := Pop(&heap)
	v2, ok2 := Pop(&heap)
	if !ok1 || !ok2 || v1 != 10 || v2 != 9 {
		t.Errorf("Unepxected pop results: %v %v, %v %v\n", v1, ok1, v2, ok2)
	}
	const expected = `
          8
    7           5
 6     2     1     4
0  3
`
	if layout := debugPrintHeap(&heap); strings.TrimSpace(expected) != strings.TrimSpace(layout) {
		t.Errorf("Unexpected heap layout:\n%v\n", layout)
	}
}

func TestFilter1(t *testing.T) {
	elems := []int{1, 20, 3, 40, 5, 60, 7, 80, 9, 100}
	var heap Heap[int, Max]
	for _, elem := range elems {
		Push(&heap, elem)
	}
	t.Logf("Before:\n%v\n", debugPrintHeap(&heap))
	i := -1
	Filter(&heap, func(elem *int) (bool, BreakOrContinue) {
		i++
		if i >= 8 {
			return true, Break
		}
		if *elem%2 == 0 {
			return true, Continue
		}
		return false, Continue
	})
	const expected = `
       100
  80        40
20   60    9
`
	if layout := debugPrintHeap(&heap); strings.TrimSpace(expected) != strings.TrimSpace(layout) {
		t.Errorf("Unexpected heap layout:\n%v\n", layout)
	}
}

func TestFilter2(t *testing.T) {
	elems := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}
	var heap Heap[int, Max]
	for _, elem := range elems {
		Push(&heap, elem)
	}
	t.Logf("Before:\n%v\n", debugPrintHeap(&heap))
	Filter(&heap, func(elem *int) (bool, BreakOrContinue) {
		return (*elem)%7 == 0, Continue
	})
	const expected = `
14
7
`
	if layout := debugPrintHeap(&heap); strings.TrimSpace(expected) != strings.TrimSpace(layout) {
		t.Errorf("Unexpected heap layout:\n%v\n", layout)
	}
}

func TestFilter3(t *testing.T) {
	elems := []int{1, 2, 3, 4, 5}
	var heap Heap[int, Max]
	for _, elem := range elems {
		Push(&heap, elem)
	}
	t.Logf("Before:\n%v\n", debugPrintHeap(&heap))
	Filter(&heap, func(elem *int) (bool, BreakOrContinue) {
		return true, Continue
	})
	const expected = `
    5
 4     2
1  3 
`
	if layout := debugPrintHeap(&heap); strings.TrimSpace(expected) != strings.TrimSpace(layout) {
		t.Errorf("Unexpected heap layout:\n%v\n", layout)
	}
}

func TestFilter4(t *testing.T) {
	elems := []int{1}
	var heap Heap[int, Max]
	for _, elem := range elems {
		Push(&heap, elem)
	}
	t.Logf("Before:\n%v\n", debugPrintHeap(&heap))
	Filter(&heap, func(elem *int) (bool, BreakOrContinue) {
		return false, Continue
	})
	const expected = ""
	if layout := debugPrintHeap(&heap); strings.TrimSpace(expected) != strings.TrimSpace(layout) {
		t.Errorf("Unexpected heap layout:\n%v\n", layout)
	}
}

func TestFilter5(t *testing.T) {
	elems := []int{1, 2, 3, 4}
	var heap Heap[int, Max]
	for _, elem := range elems {
		Push(&heap, elem)
	}
	t.Logf("Before:\n%v\n", debugPrintHeap(&heap))
	Filter(&heap, func(elem *int) (bool, BreakOrContinue) {
		return *elem == 3, Continue
	})
	const expected = `
3
`
	if layout := debugPrintHeap(&heap); strings.TrimSpace(expected) != strings.TrimSpace(layout) {
		t.Errorf("Unexpected heap layout:\n%v\n", layout)
	}
}

func TestFilter6(t *testing.T) {
	elems := []int{7, 12, 13, 15}
	var heap Heap[int, Max]
	for _, elem := range elems {
		Push(&heap, elem)
	}
	t.Logf("Before:\n%v\n", debugPrintHeap(&heap))
	Filter(&heap, func(elem *int) (bool, BreakOrContinue) {
		return *elem == 7, Continue
	})
	const expected = `
7
`
	if layout := debugPrintHeap(&heap); strings.TrimSpace(expected) != strings.TrimSpace(layout) {
		t.Errorf("Unexpected heap layout:\n%v\n", layout)
	}
}

func TestFilter7(t *testing.T) {
	elems := []int{12, 7, 13, 15}
	var heap Heap[int, Max]
	for _, elem := range elems {
		Push(&heap, elem)
	}
	t.Logf("Before:\n%v\n", debugPrintHeap(&heap))
	Filter(&heap, func(elem *int) (bool, BreakOrContinue) {
		return *elem == 7, Continue
	})
	const expected = `
7
`
	if layout := debugPrintHeap(&heap); strings.TrimSpace(expected) != strings.TrimSpace(layout) {
		t.Errorf("Unexpected heap layout:\n%v\n", layout)
	}
}

func TestFilter8(t *testing.T) {
	elems := []int{12, 13, 15, 7}
	var heap Heap[int, Max]
	for _, elem := range elems {
		Push(&heap, elem)
	}
	t.Logf("Before:\n%v\n", debugPrintHeap(&heap))
	Filter(&heap, func(elem *int) (bool, BreakOrContinue) {
		return *elem == 13, Continue
	})
	const expected = `
13
`
	if layout := debugPrintHeap(&heap); strings.TrimSpace(expected) != strings.TrimSpace(layout) {
		t.Errorf("Unexpected heap layout:\n%v\n", layout)
	}
}

func TestFilter9(t *testing.T) {
	elems := []int{1, 2, 3, 4, 100, 5, 2, 2, 2, 1, 1, 1, 1, 1, 1, 6, 2, 2, 7, 2, 2, 2, 8, 9, 10, 11, 12, 13, 14}
	var heap Heap[int, Min]
	for _, elem := range elems {
		Push(&heap, elem)
	}
	t.Logf("Before:\n%v\n", debugPrintHeap(&heap))
	Filter(&heap, func(elem *int) (bool, BreakOrContinue) {
		return (*elem)%3 == 0, Continue
	})
	const expected = `
  3
6   9
12 
`
	if layout := debugPrintHeap(&heap); strings.TrimSpace(expected) != strings.TrimSpace(layout) {
		t.Errorf("Unexpected heap layout:\n%v\n", layout)
	}
}

type myCustomType struct {
	Key     int
	Content string
}

func (c1 *myCustomType) Cmp(c2 *myCustomType) int {
	return c1.Key - c2.Key
}

func TestCustomKeyTypes(t *testing.T) {
	var heap Heap[myCustomType, Min]
	PushOrderable(&heap, myCustomType{Key: 1, Content: "foo"})
	PushOrderable(&heap, myCustomType{Key: 5, Content: "bar"})
	PushOrderable(&heap, myCustomType{Key: 2, Content: "amp"})
	PushOrderable(&heap, myCustomType{Key: 17, Content: "flub"})
	PushOrderable(&heap, myCustomType{Key: 0, Content: "zero"})
	const expected = `
               {0 zero}
     {1 foo}               {2 amp}
{17 flub}   {5 bar}
`
	if layout := debugPrintHeap(&heap); strings.TrimSpace(layout) != strings.TrimSpace(expected) {
		t.Errorf("Unexpected heap layout:\n%v\n", layout)
	}
}

// Fuzz tests a randomly generated sequence of operations against the same set
// of operations performed on a sorted slice.
func TestMinHeapFuzz(t *testing.T) {
	src := rand.NewSource(123)

	var realHeap Heap[int, Min]
	var naiveHeap []int

	for i := 0; i < 10000; i++ {
		rnd := src.Int63()
		if rnd%13 == 0 { // we'll add more than we pop or filter on average
			t.Logf("Pop")
			v1, ok1 := naiveHeapPop(&naiveHeap)
			v2, ok2 := Pop(&realHeap)
			t.Logf("(%v, %v) (%v, %v)\n", v1, ok1, v2, ok2)
		} else if rnd%17 == 0 {
			t.Logf("Filter %%7 == 0\n")
			f := func(v *int) (bool, BreakOrContinue) {
				return (*v)%7 == 0, Continue
			}
			naiveHeapFilter(&naiveHeap, f)
			Filter(&realHeap, f)
		} else {
			v := int(rnd % 100)
			t.Logf("Push %v\n", v)
			naiveMinHeapPush(&naiveHeap, v)
			Push(&realHeap, v)
		}

		if !checkMinHeapProperty(&realHeap, 0) {
			t.Fatalf("Real heap does not have min heap property:\n%v\n", debugPrintHeap(&realHeap))
		}

		if !slicesHaveSameElems(naiveHeap, realHeap.sl) {
			t.Fatalf("Elements not the same:\n%+v\n\n%v\n", naiveHeap, debugPrintHeap(&realHeap))
		}
	}

	for {
		v1, ok1 := Pop(&realHeap)
		v2, ok2 := naiveHeapPop(&naiveHeap)

		if v1 != v2 || ok1 != ok2 {
			dbg := debugPrintHeap(&realHeap)
			t.Errorf("Oh no! Got %v,%v, expected %v,%v.\n\n%v\n\n%+v\n", v1, ok1, v2, ok2, dbg, naiveHeap)
			break
		}

		if !ok1 {
			break
		}
	}
}

func TestMaxHeapFuzz(t *testing.T) {
	src := rand.NewSource(123)

	var realHeap Heap[int, Max]
	var naiveHeap []int

	for i := 0; i < 10000; i++ {
		rnd := src.Int63()
		if rnd%12 == 0 { // we'll add more than we pop on average
			naiveHeapPop(&naiveHeap)
			Pop(&realHeap)
		} else {
			v := int(rnd % 100)
			naiveMaxHeapPush(&naiveHeap, v)
			Push(&realHeap, v)
		}

		if !checkMaxHeapProperty(&realHeap, 0) {
			t.Fatalf("Real heap does not have min heap property:\n%v\n", debugPrintHeap(&realHeap))
		}

		if !slicesHaveSameElems(naiveHeap, realHeap.sl) {
			t.Fatalf("Elements not the same:\n%+v\n\n%v\n", naiveHeap, debugPrintHeap(&realHeap))
		}
	}

	for {
		v1, ok1 := Pop(&realHeap)
		v2, ok2 := naiveHeapPop(&naiveHeap)

		if v1 != v2 || ok1 != ok2 {
			t.Errorf("Oh no:\n%v\n", debugPrintHeap(&realHeap))
			break
		}

		if !ok1 {
			break
		}
	}
}

// If Push is on average O(1) then we should see times increase linearly with
// the number of elements pushed onto the heap.
func BenchmarkPush10(b *testing.B) {
	benchmarkPush(b, 10)
}

func BenchmarkPush100(b *testing.B) {
	benchmarkPush(b, 100)
}

func BenchmarkPush1000(b *testing.B) {
	benchmarkPush(b, 1000)
}

func BenchmarkPush10000(b *testing.B) {
	benchmarkPush(b, 10000)
}

func BenchmarkPush100000(b *testing.B) {
	benchmarkPush(b, 100000)
}

func benchmarkPush(b *testing.B, nElements int) {
	src := rand.NewSource(456)

	var h Heap[int, Min]
	// preconstruct the sequence of elements to insert so that we don't benchmark
	// the rng
	elems := make([]int, nElements)
	for i := range elems {
		elems[i] = int(src.Int63())
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, e := range elems {
			Push(&h, e)
		}
	}
}

// For Pop, we benchmark how long it takes to push 10 new elements and then pop
// ten times for heaps from of different sizes. Plotting log n against
// the time should give a straight line.

func BenchmarkPop100(b *testing.B) {
	benchmarkPop(b, 100)
}

func BenchmarkPop200(b *testing.B) {
	benchmarkPop(b, 200)
}

func BenchmarkPop400(b *testing.B) {
	benchmarkPop(b, 400)
}

func BenchmarkPop800(b *testing.B) {
	benchmarkPop(b, 800)
}

func BenchmarkPop1600(b *testing.B) {
	benchmarkPop(b, 1600)
}

func BenchmarkPop3200(b *testing.B) {
	benchmarkPop(b, 3200)
}

func BenchmarkPop6400(b *testing.B) {
	benchmarkPop(b, 6400)
}

func benchmarkPop(b *testing.B, nElements int) {
	src := rand.NewSource(789)
	var h Heap[int, Min]
	for i := 0; i < nElements; i++ {
		Push(&h, int(src.Int63()))
	}

	elems := make([]int, 10)
	for i := range elems {
		elems[i] = int(src.Int63())
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, e := range elems {
			Push(&h, e)
			Pop(&h)
		}
	}
}
