package tree

import (
	"GiveMeAnOfferGo/Collections"
	"GiveMeAnOfferGo/Objects"
	"fmt"
	"testing"
)

func TestBuildHeap(t *testing.T) {
	input := []Objects.ComparableObject{
		getInt(1),
		getInt(12),
		getInt(3),
		getInt(4),
		getInt(1),
		getInt(6),
		getInt(8),
		getInt(7),
	}

	heap := Collections.NewHeap(func(obj1 Objects.ComparableObject, obj2 Objects.ComparableObject) bool {
		return obj1.Compare(obj2) == Objects.OrderedDescending
	}, input)

	for !heap.IsEmpty() {
		fmt.Println(heap.Remove())
	}
}
