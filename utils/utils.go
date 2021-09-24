package utils

import (
	"fmt"
	"reflect"
)

// IndirectToInterface returns the interface's value
// If the interface is reflect.Value, return the reflect.Value.Interface()
func IndirectToInterface(v interface{}) interface{} {
	switch e := v.(type) {
	case reflect.Value:
		return e.Interface()
	}

	return v
}

// IndirectToValue returns the reflect.Value of v
// If the interface is non reflect.Value, return the reflect.ValueOf(v)
func IndirectToValue(v interface{}) reflect.Value {
	switch e := v.(type) {
	case reflect.Value:
		return e
	}

	return reflect.ValueOf(v)
}

// NewValue return the value of type:
// 1. If type is Primitive type, it will return zero value
// 2. If type is Ptr type, it will return ptr to zero value
// 3. Other wise, it will return an error
func NewValue(typ reflect.Type) (reflect.Value, error) {
	if typ.Kind() == reflect.Ptr {
		if typ.Elem().Kind() == reflect.Ptr {
			return reflect.Value{}, fmt.Errorf("Unsupport type '%v'", typ.String())
		}

		ret := reflect.New(typ.Elem())
		ret.Elem().Set(reflect.Zero(typ.Elem()))
		return ret, nil
	}
	return reflect.Zero(typ), nil
}
