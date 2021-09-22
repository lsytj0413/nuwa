package nuwa

import (
	"fmt"
	"reflect"
)

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

// IsPrimitiveType return nil error if type is primitive
func IsPrimitiveType(typ reflect.Type) error {
	if v, ok := primitiveTypeMaps[typ.Kind()]; ok {
		if v {
			return nil
		}

		return fmt.Errorf("Kind '%v' is not primitive type", typ.Kind())
	}

	return fmt.Errorf("Kind '%v' not found in primitive types", typ.Kind())
}

var primitiveTypeMaps map[reflect.Kind]bool = map[reflect.Kind]bool{
	reflect.Invalid:       false,
	reflect.Bool:          true,
	reflect.Int:           true,
	reflect.Int8:          true,
	reflect.Int16:         true,
	reflect.Int32:         true,
	reflect.Int64:         true,
	reflect.Uint:          true,
	reflect.Uint8:         true,
	reflect.Uint16:        true,
	reflect.Uint32:        true,
	reflect.Uint64:        true,
	reflect.Uintptr:       true,
	reflect.Float32:       true,
	reflect.Float64:       true,
	reflect.Complex64:     true,
	reflect.Complex128:    true,
	reflect.Array:         false,
	reflect.Chan:          false,
	reflect.Func:          false,
	reflect.Interface:     false,
	reflect.Map:           false,
	reflect.Ptr:           true,
	reflect.Slice:         false,
	reflect.String:        true,
	reflect.Struct:        true,
	reflect.UnsafePointer: false,
}
