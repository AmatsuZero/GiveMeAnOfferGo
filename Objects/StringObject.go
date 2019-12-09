package Objects

import (
	"fmt"
	"strings"
)

type StringObject struct {
	GoString string
}

func (s *StringObject) ToSlice() (ret []string) {
	for _, s := range s.GoString {
		ret = append(ret, fmt.Sprintf("%c", s))
	}
	return
}

func (s *StringObject) ToObjectSlice() (ret []*StringObject) {
	for _, s := range s.ToSlice() {
		obj := &StringObject{GoString: s}
		ret = append(ret, obj)
	}
	return
}

func (s *StringObject) String() string {
	return s.GoString
}

func (s *StringObject) Compare(obj interface{}) CompareResult {
	object, ok := obj.(*StringObject)
	if !ok {
		return InvalidResult
	}

	switch strings.Compare(s.GoString, object.GoString) {
	case 0:
		return OrderedSame
	case -1:
		return OrderedAscending
	case 1:
		return OrderedDescending
	default:
		return InvalidResult
	}
}

func (s *StringObject) IsEqualTo(obj interface{}) bool {
	return s.Compare(obj) == OrderedSame
}
