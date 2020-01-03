package Objects

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
)

var floatType = reflect.TypeOf(float64(0))
var stringType = reflect.TypeOf("")

type NumberObject struct {
	value interface{}
}

func NewNumberWithInt(integer int) *NumberObject {
	return &NumberObject{value: integer}
}

func getFloat(unk interface{}) (float64, error) {
	switch i := unk.(type) {
	case float64:
		return i, nil
	case float32:
		return float64(i), nil
	case int64:
		return float64(i), nil
	case int32:
		return float64(i), nil
	case int:
		return float64(i), nil
	case uint64:
		return float64(i), nil
	case uint32:
		return float64(i), nil
	case uint:
		return float64(i), nil
	case string:
		return strconv.ParseFloat(i, 64)
	default:
		v := reflect.ValueOf(unk)
		v = reflect.Indirect(v)
		if v.Type().ConvertibleTo(floatType) {
			fv := v.Convert(floatType)
			return fv.Float(), nil
		} else if v.Type().ConvertibleTo(stringType) {
			sv := v.Convert(stringType)
			s := sv.String()
			return strconv.ParseFloat(s, 64)
		} else {
			return math.NaN(), fmt.Errorf("can't convert %v to float64", v.Type())
		}
	}
}

func (number *NumberObject) FloatValue() float64 {
	if val, err := getFloat(number.value); err == nil {
		return val
	}
	return 0
}

func (number *NumberObject) Compare(obj interface{}) CompareResult {
	if obj == number {
		return OrderedSame
	}
	object, ok := obj.(NumberObject)
	if !ok {
		tmp, ok := obj.(*NumberObject)
		if ok {
			object = *tmp
		} else {
			return InvalidResult
		}
	}
	lhs := number.FloatValue()
	rhs := object.FloatValue()
	if lhs < rhs {
		return OrderedAscending
	} else if lhs == rhs {
		return OrderedSame
	} else {
		return OrderedDescending
	}
}

func (number *NumberObject) IsEqualTo(obj interface{}) bool {
	if obj == nil {
		return false
	}
	return number.Compare(obj) == OrderedSame
}

func (number *NumberObject) String() string {
	return fmt.Sprint(number.value)
}

func (number *NumberObject) IsNil() bool {
	return number.value == nil
}
