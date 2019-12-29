package queue

import (
	"GiveMeAnOfferGo/Collections"
	"GiveMeAnOfferGo/Objects"
	"fmt"
	"testing"
)

func TestBuildQueue(t *testing.T) {
	priorityQueue := Collections.NewPriorityQueue(func(lhs Objects.ComparableObject, rhs Objects.ComparableObject) bool {
		return lhs.Compare(rhs) == Objects.OrderedDescending
	}, getInt(1), getInt(12), getInt(3), getInt(4),
		getInt(1), getInt(6), getInt(8), getInt(7))

	for !priorityQueue.IsEmpty() {
		fmt.Println(priorityQueue.Dequeue())
	}
}
