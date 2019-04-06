package value

import (
	"fmt"
	"reflect"
)

// Value is the interface for all constant values.
type Value interface {
	fmt.Stringer
	Value() interface{}
}

// FromValue takes an interface value and creates a constant.
func FromValue(v interface{}) (Value, error) {
	switch t := v.(type) {
	case uint8:
		return Uint{uint64(t)}, nil
	case int8:
		return Int{int64(t)}, nil
	case uint16:
		return Uint{uint64(t)}, nil
	case int16:
		return Int{int64(t)}, nil
	case uint32:
		return Uint{uint64(t)}, nil
	case int32:
		return Int{int64(t)}, nil
	case uint64:
		return Uint{uint64(t)}, nil
	case int64:
		return Int{int64(t)}, nil
	case uint:
		return Uint{uint64(t)}, nil
	case int:
		return Int{int64(t)}, nil
	case string:
		return String{t}, nil
	case float64:
		return Float{t}, nil
	case bool:
		return Boolean{t}, nil
	default:
		return fromReflectValue(reflect.ValueOf(v))
	}
}

func fromReflectValue(v reflect.Value) (Value, error) {
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		return NewArray(v.Interface()), nil
	case reflect.Func:
		return NewFunc(v.Interface()), nil
	case reflect.Map:
		return NewMap(v.Interface()), nil
	case reflect.Struct:
		return NewStruct(v.Interface()), nil
	case reflect.Ptr:
		return fromReflectValue(v.Elem())
	default:
		return nil, fmt.Errorf("unexpected type %s", v.Kind())
	}
}
