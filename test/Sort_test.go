package test

import (
	"fmt"
	"github.com/AmatsuZero/GiveMeAnOfferGo/Collections"
	"github.com/AmatsuZero/GiveMeAnOfferGo/Objects"
	"github.com/AmatsuZero/GiveMeAnOfferGo/test/Utils"
	"testing"
)

func TestBubbleSort(t *testing.T) {
	array := []Objects.Comparable{Utils.GetInt(9), Utils.GetInt(4), Utils.GetInt(10), Utils.GetInt(3)}
	fmt.Printf("Original: %v\n", array)
	Collections.BubbleSort(&array)
	fmt.Printf("Bubble sorted: %v\n", array)
}

func TestSelectionSort(t *testing.T) {
	array := []Objects.Comparable{Utils.GetInt(9), Utils.GetInt(4), Utils.GetInt(10), Utils.GetInt(3)}
	fmt.Printf("Original: %v\n", array)
	Collections.SelectionSort(&array)
	fmt.Printf("Selection sorted: %v\n", array)
}

func TestInsertionSort(t *testing.T) {
	array := []Objects.Comparable{Utils.GetInt(9), Utils.GetInt(4), Utils.GetInt(10), Utils.GetInt(3)}
	fmt.Printf("Original: %v\n", array)
	Collections.InsertionSort(&array)
	fmt.Printf("Insertion sorted: %v\n", array)
}

func TestMergeSort(t *testing.T) {
	array := []Objects.Comparable{
		Utils.GetInt(7),
		Utils.GetInt(2),
		Utils.GetInt(6),
		Utils.GetInt(3),
		Utils.GetInt(9),
	}
	fmt.Printf("Original: %v\n", array)
	fmt.Printf("Merge sorted: %v\n", Collections.MergeSort(array))
}

func TestRadixSort(t *testing.T) {
	array := []int{88, 410, 1772, 20}
	fmt.Printf("Original: %v\n", array)
	fmt.Printf("Radix sorted: %v\n", Collections.RadixSort(array))
}

func TestHeapSort(t *testing.T) {
	array := []Objects.ComparableObject{
		Utils.GetInt(6),
		Utils.GetInt(12),
		Utils.GetInt(2),
		Utils.GetInt(26),
		Utils.GetInt(8),
		Utils.GetInt(18),
		Utils.GetInt(21),
		Utils.GetInt(9),
		Utils.GetInt(5),
	}

	heap := Collections.NewHeap(func(obj1 Objects.ComparableObject, obj2 Objects.ComparableObject) bool {
		return obj1.Compare(obj2) == Objects.OrderedDescending
	}, array)

	fmt.Println(heap.Sorted())
}

func TestQuickSort(t *testing.T) {
	array := []Objects.Comparable{
		Utils.GetInt(12),
		Utils.GetInt(0),
		Utils.GetInt(3),
		Utils.GetInt(9),
		Utils.GetInt(2),
		Utils.GetInt(21),
		Utils.GetInt(8),
		Utils.GetInt(18),
		Utils.GetInt(27),
		Utils.GetInt(1),
		Utils.GetInt(5),
		Utils.GetInt(8),
		Utils.GetInt(-1),
	}

	fmt.Println(Collections.QuickSortNative(array))
}

func TestQuickSortLomuto(t *testing.T) {
	array := []Objects.Comparable{
		Utils.GetInt(12),
		Utils.GetInt(0),
		Utils.GetInt(3),
		Utils.GetInt(9),
		Utils.GetInt(2),
		Utils.GetInt(21),
		Utils.GetInt(18),
		Utils.GetInt(27),
		Utils.GetInt(1),
		Utils.GetInt(5),
		Utils.GetInt(8),
		Utils.GetInt(-1),
		Utils.GetInt(8),
	}
	Collections.QuickSortLomuto(&array, 0, len(array)-1)
	fmt.Println(array)
}

func TestQuickSortHoare(t *testing.T) {
	array := []Objects.Comparable{
		Utils.GetInt(12),
		Utils.GetInt(0),
		Utils.GetInt(3),
		Utils.GetInt(9),
		Utils.GetInt(2),
		Utils.GetInt(21),
		Utils.GetInt(18),
		Utils.GetInt(27),
		Utils.GetInt(1),
		Utils.GetInt(5),
		Utils.GetInt(8),
		Utils.GetInt(-1),
		Utils.GetInt(8),
	}
	Collections.QuickSortHoare(&array, 0, len(array)-1)
	fmt.Println(array)
}

func TestQuickSortMedian(t *testing.T) {
	array := []Objects.Comparable{
		Utils.GetInt(12),
		Utils.GetInt(0),
		Utils.GetInt(3),
		Utils.GetInt(9),
		Utils.GetInt(2),
		Utils.GetInt(21),
		Utils.GetInt(18),
		Utils.GetInt(27),
		Utils.GetInt(1),
		Utils.GetInt(5),
		Utils.GetInt(8),
		Utils.GetInt(-1),
		Utils.GetInt(8),
	}
	Collections.QuickSortMedian(&array, 0, len(array)-1)
	t.Log(array)
}

func TestQucikSortDutchFlag(t *testing.T) {
	array := []Objects.Comparable{
		Utils.GetInt(12),
		Utils.GetInt(0),
		Utils.GetInt(3),
		Utils.GetInt(9),
		Utils.GetInt(2),
		Utils.GetInt(21),
		Utils.GetInt(18),
		Utils.GetInt(27),
		Utils.GetInt(1),
		Utils.GetInt(5),
		Utils.GetInt(8),
		Utils.GetInt(-1),
		Utils.GetInt(8),
	}
	Collections.QuickSortDutchFlag(&array, 0, len(array)-1)
	t.Log(array)
}
