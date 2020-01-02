package Collections

import (
	"GiveMeAnOfferGo/Objects"
)

func BubbleSort(collections *[]Objects.Comparable) {
	array := *collections
	if len(array) < 2 {
		return
	}
	for end := len(array) - 1; end > 0; end-- {
		swapped := false
		for current := 0; current < end; current++ {
			if array[current].Compare(array[current+1]) == Objects.OrderedDescending {
				array[current+1], array[current] = array[current], array[current+1]
				swapped = true
			}
		}
		if !swapped {
			return
		}
	}
}

func SelectionSort(collections *[]Objects.Comparable) {
	array := *collections
	if len(array) < 2 {
		return
	}
	for current := 0; current < len(array)-1; current++ {
		lowest := current
		for other := current + 1; other < len(array); other++ {
			if array[lowest].Compare(array[other]) == Objects.OrderedDescending {
				lowest = other
			}
		}
		if lowest != current {
			array[lowest], array[current] = array[current], array[lowest]
		}
	}
}

func InsertionSort(collections *[]Objects.Comparable) {
	array := *collections
	if len(array) < 2 {
		return
	}
	for current := 1; current < len(array); current++ {
		for shifting := current; shifting > 0; shifting-- {
			if array[shifting].Compare(array[shifting-1]) == Objects.OrderedAscending {
				array[shifting], array[shifting-1] = array[shifting-1], array[shifting]
			} else {
				break
			}
		}
	}
}

func MergeSort(array []Objects.Comparable) []Objects.Comparable {
	if len(array) <= 1 {
		return array
	}

	middle := len(array) / 2
	left := MergeSort(array[:middle])
	right := MergeSort(array[middle:])

	return merge(left, right)
}

func merge(left []Objects.Comparable, right []Objects.Comparable) (result []Objects.Comparable) {
	leftIndex := 0
	rightIndex := 0

	result = make([]Objects.Comparable, 0)

	for leftIndex < len(left) && rightIndex < len(right) {
		leftElement := left[leftIndex]
		rightElement := right[rightIndex]

		if leftElement.Compare(rightElement) == Objects.OrderedAscending {
			result = append(result, leftElement)
			leftIndex += 1
		} else if leftElement.Compare(rightElement) == Objects.OrderedDescending {
			result = append(result, rightElement)
			rightIndex += 1
		} else {
			result = append(result, leftElement)
			leftIndex += 1
			result = append(result, rightElement)
			rightIndex += 1
		}
	}

	if leftIndex < len(left) {
		result = append(result, left[leftIndex:]...)
	}

	if rightIndex < len(right) {
		result = append(result, right[rightIndex:]...)
	}

	return
}

func RadixSort(array []int) []int {
	base := 10
	done := false
	digits := 1

	for !done {
		buckets := make([][]int, base)
		for _, number := range array {
			done = true

			remainingPart := number / digits
			digit := remainingPart % base
			buckets[digit] = append(buckets[digit], number)

			if remainingPart > 0 {
				done = false
			}
		}

		digits *= base
		array = make([]int, 0)
		for _, bucket := range buckets {
			array = append(array, bucket...)
		}
	}

	return array
}

func QuickSortNative(collection []Objects.Comparable) []Objects.Comparable {
	if len(collection) <= 1 {
		return collection
	}

	pivot := collection[len(collection)/2]
	var less, equal, greater []Objects.Comparable
	for _, element := range collection {
		result := element.Compare(pivot)
		if result == Objects.OrderedAscending {
			less = append(less, element)
		} else if result == Objects.OrderedSame {
			equal = append(equal, element)
		} else {
			greater = append(greater, element)
		}
	}

	ret := append(QuickSortNative(less), equal...)
	return append(ret, QuickSortNative(greater)...)
}

func PartitionLomuto(array *[]Objects.Comparable, low int, high int) int {
	a := *array
	pivot := a[high]
	i := low
	for j := low; j < high; j++ {
		if a[j].Compare(pivot) != Objects.OrderedDescending {
			a[i], a[j] = a[j], a[i]
			i += 1
		}
	}

	a[i], a[high] = a[high], a[i]
	return i
}

func QuickSortLomuto(array *[]Objects.Comparable, low int, high int) {
	if low < high {
		pivot := PartitionLomuto(array, low, high)
		QuickSortLomuto(array, low, pivot-1)
		QuickSortLomuto(array, pivot+1, high)
	}
}

func PartitionHoare(array *[]Objects.Comparable, low int, high int) int {
	a := *array
	pivot := a[low]
	i := low - 1
	j := high + 1

	for true {
		for j -= 1; a[j].Compare(pivot) == Objects.OrderedDescending; j-- {
		}

		for i += 1; a[i].Compare(pivot) == Objects.OrderedAscending; i++ {
		}

		if i < j {
			a[i], a[j] = a[j], a[i]
		} else {
			return j
		}
	}

	return j
}

func QuickSortHoare(array *[]Objects.Comparable, low int, high int) {
	if low < high {
		p := PartitionHoare(array, low, high)
		QuickSortHoare(array, low, p)
		QuickSortHoare(array, p+1, high)
	}
}

func MedianOfThree(array *[]Objects.Comparable, low int, high int) int {
	a := *array
	center := (low + high) / 2
	if a[low].Compare(a[center]) == Objects.OrderedDescending {
		a[low], a[center] = a[center], a[low]
	}

	if a[low].Compare(a[high]) == Objects.OrderedDescending {
		a[low], a[high] = a[high], a[low]
	}

	if a[center].Compare(a[high]) == Objects.OrderedDescending {
		a[center], a[high] = a[high], a[center]
	}

	return center
}

func QuickSortMedian(array *[]Objects.Comparable, low int, high int) {
	a := *array
	if low < high {
		pivotIndex := MedianOfThree(array, low, high)
		a[pivotIndex], a[high] = a[high], a[pivotIndex]
		pivot := PartitionLomuto(array, low, high)
		QuickSortLomuto(array, low, pivot-1)
		QuickSortLomuto(array, pivot+1, high)
	}
}

func PartitionDutchFlag(array *[]Objects.Comparable, low int, high int, pivotIndex int) (smaller int, larger int) {
	a := *array
	pivot := a[pivotIndex]
	smaller = low
	larger = high
	equal := low
	for equal <= larger {
		ret := a[equal].Compare(pivot)
		if ret == Objects.OrderedAscending {
			a[smaller], a[equal] = a[equal], a[smaller]
			smaller += 1
			equal += 1
		} else if ret == Objects.OrderedSame {
			equal += 1
		} else {
			a[equal], a[larger] = a[larger], a[equal]
			larger -= 1
		}
	}

	return
}

func QuickSortDutchFlag(array *[]Objects.Comparable, low int, high int) {
	if low < high {
		middleFirst, middleLast := PartitionDutchFlag(array, low, high, high)
		QuickSortDutchFlag(array, low, middleFirst-1)
		QuickSortDutchFlag(array, middleLast+1, high)
	}
}
