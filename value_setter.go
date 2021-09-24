package nuwa

import (
	"fmt"
	"reflect"
)

// ValueSetter is the setter for value
type ValueSetter interface {
	Set(v interface{})
}

var _ ValueSetter = &valueSetter{}

type valueSetter struct {
	val reflect.Value
}

func (s *valueSetter) Set(v interface{}) {
	vv, ok := v.(reflect.Value)
	if ok {
		s.val.Set(vv)
	}

	s.val.Set(reflect.ValueOf(v))
}

// NewValueSetter return the setter for value
func NewValueSetter(v reflect.Value) (ValueSetter, error) {
	if v.Type().Kind() == reflect.Ptr {
		return &valueSetter{
			val: v.Elem(),
		}, nil
	}

	if !v.CanSet() {
		return nil, fmt.Errorf("Cannot set value for typ '%v'", v.Type().String())
	}

	return &valueSetter{
		val: v,
	}, nil
}

// NewFieldValueSetter return the setter for struct field value
func NewFieldValueSetter(v reflect.Value, idx int) (ValueSetter, error) {
	if v.Type().Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Type().Kind() != reflect.Struct {
		return nil, fmt.Errorf("Cannot get field '%v' value setter for typ '%v': It must be struct", idx, v.Type().String())
	}

	if idx < 0 || idx >= v.Type().NumField() {
		return nil, fmt.Errorf("Cannot get field '%v' value setter for typ '%v': idx must in range [0,%v)", idx, v.Type().String(), v.Type().NumField())
	}

	// NOTE: we cannot use NewValueSetter, because it will redirect to set ptr value with derefence,
	// If the ptr is nil, it will cause Set on zero value
	// return NewValueSetter(v.Field(idx))

	v = v.Field(idx)
	if !v.CanSet() {
		return nil, fmt.Errorf("Cannot get field '%v' value setter for typ '%v': It must be setable", idx, v.Type().String())
	}

	return &valueSetter{
		val: v,
	}, nil
}
