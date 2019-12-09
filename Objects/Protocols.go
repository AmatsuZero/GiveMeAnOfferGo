package Objects

import (
	"fmt"
)

type CompareResult int

const InvalidResult = -1

const (
	OrderedAscending CompareResult = iota
	OrderedSame
	OrderedDescending
)

type Equatable interface {
	IsEqualTo(obj interface{}) bool
}

type Comparable interface {
	fmt.Stringer
	Equatable
	Compare(obj interface{}) CompareResult
}
