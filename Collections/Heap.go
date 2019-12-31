package Collections

import (
	"GiveMeAnOfferGo/Objects"
	"fmt"
)

type Heap struct {
	elements []Objects.ComparableObject
	sort     func(obj1 Objects.ComparableObject, obj2 Objects.ComparableObject) bool
}

func NewHeap(sort func(obj1 Objects.ComparableObject, obj2 Objects.ComparableObject) bool, elements []Objects.ComparableObject) *Heap {
	if elements == nil {
		elements = make([]Objects.ComparableObject, 0)
	}
	heap := &Heap{
		elements: elements,
		sort:     sort,
	}
	if !heap.IsEmpty() {
		for i := heap.Count()/2 - 1; i >= 0; i-- {
			heap.siftDown(i, len(elements))
		}
	}
	return heap
}

func (heap *Heap) IsEmpty() bool {
	return heap.Count() == 0
}

func (heap *Heap) Count() int {
	return len(heap.elements)
}

func (heap *Heap) Peek() Objects.ComparableObject {
	if heap.IsEmpty() {
		return nil
	}
	return heap.elements[0]
}

func (heap *Heap) LeftChildIndex(parentIndex int) int {
	return 2*parentIndex + 1
}

func (heap *Heap) RightChildIndex(parentIndex int) int {
	return 2*parentIndex + 2
}

func (heap *Heap) ParentIndex(childIndex int) int {
	return (childIndex - 1) / 2
}

func (heap *Heap) Remove() (popped Objects.ComparableObject) {
	if heap.IsEmpty() {
		return nil
	}
	heap.elements[0], heap.elements[heap.Count()-1] = heap.elements[heap.Count()-1], heap.elements[0]
	defer func() {
		heap.siftDown(0, heap.Count())
	}()

	popped, heap.elements = heap.elements[len(heap.elements)-1], heap.elements[:len(heap.elements)-1]
	return
}

func (heap *Heap) siftDown(fromIndex int, size int) {
	parent := fromIndex
	for {
		left := heap.LeftChildIndex(parent)
		right := heap.RightChildIndex(parent)
		candidate := parent

		if left < size && heap.sort(heap.elements[left], heap.elements[candidate]) {
			candidate = left
		}

		if right < size && heap.sort(heap.elements[right], heap.elements[candidate]) {
			candidate = right
		}

		if candidate == parent {
			return
		}

		heap.elements[parent], heap.elements[candidate] = heap.elements[candidate], heap.elements[parent]
		parent = candidate
	}
}

func (heap *Heap) Insert(element Objects.ComparableObject) {
	heap.elements = append(heap.elements, element)
	heap.siftUp(heap.Count() - 1)
}

func (heap *Heap) RemoveAt(index int) (last Objects.ComparableObject) {
	if index >= heap.Count() {
		return nil
	}

	if index == heap.Count()-1 {
		last, heap.elements = heap.elements[len(heap.elements)-1], heap.elements[:len(heap.elements)-1]
	} else {
		heap.elements[index], heap.elements[heap.Count()-1] = heap.elements[heap.Count()-1], heap.elements[index]
		defer func() {
			heap.siftDown(index, heap.Count())
			heap.siftUp(index)
		}()
		last, heap.elements = heap.elements[len(heap.elements)-1], heap.elements[:len(heap.elements)-1]
	}
	return
}

func (heap *Heap) siftUp(fromIndex int) {
	child := fromIndex
	parent := heap.ParentIndex(child)
	for child > 0 && heap.sort(heap.elements[child], heap.elements[parent]) {
		heap.elements[child], heap.elements[parent] = heap.elements[parent], heap.elements[child]
		child = parent
		parent = heap.ParentIndex(child)
	}
}

func (heap *Heap) IndexOf(element Objects.ComparableObject, i int) int {
	if i >= heap.Count() {
		return Objects.InvalidResult
	}

	if heap.sort(element, heap.elements[i]) {
		return Objects.InvalidResult
	}

	if element.IsEqualTo(heap.elements[i]) {
		return i
	}

	if j := heap.IndexOf(element, heap.LeftChildIndex(i)); j != Objects.InvalidResult {
		return j
	}

	if j := heap.IndexOf(element, heap.RightChildIndex(i)); j != Objects.InvalidResult {
		return j
	}

	return Objects.InvalidResult
}

func (heap *Heap) String() string {
	str := fmt.Sprintf("Heap %p\n", heap)
	if heap.IsEmpty() {
		return str + "empty"
	}
	for _, element := range heap.elements {
		str += fmt.Sprintln(element)
	}
	return str
}

func (heap *Heap) Sorted() []Objects.ComparableObject {
	if heap.IsEmpty() {
		return heap.elements
	}
	newHeap := NewHeap(heap.sort, heap.elements)
	for index := len(heap.elements) - 1; index >= 0; index-- {
		newHeap.elements[0], newHeap.elements[index] = newHeap.elements[index], newHeap.elements[0]
		newHeap.siftDown(0, index)
	}
	return newHeap.elements
}
