package property

import (
	"fmt"
	"reflect"

	"github.com/spf13/cast"

	"github.com/lsytj0413/nuwa/xerrors"
)

// NewProperties return the Properties impl
func NewProperties() Properties {
	return propertiesImpl(make(map[string]string))
}

type propertiesImpl map[string]string

func (p propertiesImpl) Get(key string) (string, error) {
	val, ok := p[key]
	if ok {
		return val, nil
	}

	return "", xerrors.WrapNotFound("property with key='%v' not found", key)
}

func (p propertiesImpl) Retrive(key string, i interface{}) error {
	vstr, err := p.Get(key)
	if err != nil {
		return err
	}

	var v reflect.Value
	switch e := i.(type) {
	case reflect.Value:
		v = e
	default:
		v = reflect.ValueOf(i)
	}

	// If the kind is ptr, we will set the value to it's elem
	// Because Ptr is not setable
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if !v.CanSet() {
		return xerrors.Errorf("The '%T' cannot been set, it must setable", i)
	}

	switch v.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u, err := cast.ToUint64E(vstr)
		if err != nil {
			return xerrors.Wrapf(err, "Cannot convert value '%v' to uint with key '%v'", vstr, key)
		}

		// NOTE: this value will be zero if the u overflow. FIX IT
		v.SetUint(u)
		return nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		u, err := cast.ToInt64E(vstr)
		if err != nil {
			return xerrors.Wrapf(err, "Cannot convert value '%v' to int with key '%v'", vstr, key)
		}

		// NOTE: this value will be zero if the u overflow. FIX IT
		v.SetInt(u)
		return nil
	case reflect.Float32, reflect.Float64:
		u, err := cast.ToFloat64E(vstr)
		if err != nil {
			return xerrors.Wrapf(err, "Cannot convert value '%v' to float with key '%v'", vstr, key)
		}

		// NOTE: this value will be zero if the u overflow. FIX IT
		v.SetFloat(u)
		return nil
	case reflect.Bool:
		u, err := cast.ToBoolE(vstr)
		if err != nil {
			return xerrors.Wrapf(err, "Cannot convert value '%v' to bool with key '%v'", vstr, key)
		}
		v.SetBool(u)
		return nil
	case reflect.String:
		v.SetString(vstr)
		return nil
	}

	return xerrors.Errorf("Cannot retrive value for key '%v', unsupported target type '%v'", key, v.Kind())
}

func (p propertiesImpl) Set(key string, val interface{}) error {
	switch v := reflect.ValueOf(val); v.Kind() {
	case reflect.Map:
		// If the val is a map, we expand the val with keys and set it recursive
		for _, k := range v.MapKeys() {
			kstr, err := cast.ToStringE(k.Interface())
			if err != nil {
				return xerrors.Wrapf(err, "Cannot convert map's key '%v' to string", k)
			}

			kstr = fmt.Sprintf("%s.%s", key, kstr)
			kvalue := v.MapIndex(k).Interface()
			err = p.Set(kstr, kvalue)
			if err != nil {
				return xerrors.Wrapf(err, "Cannot set val for map's key '%v'", kstr)
			}
		}
	case reflect.Array, reflect.Slice:
		// If the val is a array/slice, we expand the val with index and set it recursive
		for i := 0; i < v.Len(); i++ {
			kstr := fmt.Sprintf("%s[%d]", key, i)
			kvalue := v.Index(i).Interface()
			err := p.Set(kstr, kvalue)
			if err != nil {
				return xerrors.Wrapf(err, "Cannot set val for array/slice index's key '%v'", kstr)
			}
		}
	default:
		value, err := cast.ToStringE(val)
		if err != nil {
			return xerrors.Wrapf(err, "Cannot convert value to string")
		}
		p[key] = value
	}

	return nil
}
