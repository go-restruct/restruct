package typing

import (
	"errors"
	"fmt"
	"reflect"
)

// ErrReturnArity is the error returned when function return arity is invalid.
var ErrReturnArity = errors.New("functions may only return exactly 1 value")

// FromValue takes an interface value and returns a type.
func FromValue(v interface{}) (Type, error) {
	switch v.(type) {
	case uint, uint8, uint16, uint32, uint64:
		return PrimitiveType(Uint, TypeInfoFromVal(v)), nil
	case int, int8, int16, int32, int64:
		return PrimitiveType(Int, TypeInfoFromVal(v)), nil
	case string:
		return PrimitiveType(String, TypeInfoFromVal(v)), nil
	case float64:
		return PrimitiveType(Float, TypeInfoFromVal(v)), nil
	case bool:
		return PrimitiveType(Boolean, TypeInfoFromVal(v)), nil
	default:
		return FromReflectType(reflect.ValueOf(v).Type())
	}
}

// FromReflectType returns a type from a reflect.Type.
func FromReflectType(t reflect.Type) (Type, error) {
	switch t.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return PrimitiveType(Uint, TypeInfoFromType(t)), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return PrimitiveType(Int, TypeInfoFromType(t)), nil
	case reflect.String:
		return PrimitiveType(String, TypeInfoFromType(t)), nil
	case reflect.Float32, reflect.Float64:
		return PrimitiveType(Float, TypeInfoFromType(t)), nil
	case reflect.Bool:
		return PrimitiveType(Boolean, TypeInfoFromType(t)), nil
	case reflect.Slice, reflect.Array:
		elem, err := FromReflectType(t.Elem())
		if err != nil {
			return nil, err
		}
		return ArrayType(elem, TypeInfoFromType(t)), nil
	case reflect.Func:
		if t.NumOut() != 1 {
			return nil, ErrReturnArity
		}
		nparams := t.NumIn()
		params := make([]Type, 0, nparams)
		for i := 0; i < nparams; i++ {
			param, err := FromReflectType(t.In(i))
			if err != nil {
				return nil, err
			}
			params = append(params, param)
		}
		returns, err := FromReflectType(t.Out(0))
		if err != nil {
			return nil, err
		}
		return FuncType(params, returns, t.IsVariadic(), TypeInfoFromType(t)), nil
	case reflect.Map:
		key, err := FromReflectType(t.Key())
		if err != nil {
			return nil, err
		}
		elem, err := FromReflectType(t.Elem())
		if err != nil {
			return nil, err
		}
		return MapType(key, elem, TypeInfoFromType(t)), nil
	case reflect.Struct:
		nfields := t.NumField()
		fields := make([]StructField, 0, nfields)
		for i := 0; i < nfields; i++ {
			sf := t.Field(i)
			field, err := FromReflectType(sf.Type)
			if err != nil {
				return nil, err
			}
			fields = append(fields, StructField{Name: sf.Name, Type: field})
		}
		return StructType(fields, TypeInfoFromType(t)), nil
	case reflect.Ptr:
		return FromReflectType(t.Elem())
	default:
		return nil, fmt.Errorf("unexpected type %s", t.Kind())
	}
}
