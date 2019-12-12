package Objects

import (
	"fmt"
	"strings"
)

type StringObject struct {
	GoString *string
}

func (s *StringObject) ToSlice() (ret []string) {
	str := *s.GoString
	for _, s := range str {
		ret = append(ret, fmt.Sprintf("%c", s))
	}
	return
}

func (s *StringObject) ToObjectSlice() (ret []*StringObject) {
	for _, s := range s.ToSlice() {
		obj := &StringObject{GoString: &s}
		ret = append(ret, obj)
	}
	return
}

func (s *StringObject) String() string {
	return *s.GoString
}

func (s *StringObject) Compare(obj interface{}) CompareResult {
	object, ok := obj.(*StringObject)
	if !ok {
		return InvalidResult
	}
	if s.GoString == nil && object.GoString == nil {
		return OrderedSame
	} else if s.GoString == nil || object.GoString == nil {
		if s.GoString != nil {
			return OrderedDescending
		} else {
			return OrderedAscending
		}
	}
	lhs := *s.GoString
	rhs := *object.GoString
	switch strings.Compare(lhs, rhs) {
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

func (s *StringObject) IsNil() bool {
	return s.GoString == nil
}
