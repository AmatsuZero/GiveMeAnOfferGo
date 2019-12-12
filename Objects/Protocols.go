package Objects

import "fmt"

type CompareResult int

const InvalidResult = iota - 1986

const (
	OrderedAscending CompareResult = iota
	OrderedSame
	OrderedDescending
)

type Equatable interface {
	IsEqualTo(obj interface{}) bool
}

type Comparable interface {
	Equatable
	Compare(obj interface{}) CompareResult
}

type ObjectProtocol interface {
	Comparable
	fmt.Stringer
	IsNil() bool
}
