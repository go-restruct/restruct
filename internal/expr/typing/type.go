package typing

import (
	"errors"
	"fmt"
	"reflect"
)

// ErrInvalidKind occurs when you call an inappropriate method for a given kind.
var ErrInvalidKind = errors.New("invalid kind")

// NoSuchFieldError is returned when an unknown field is accessed.
type NoSuchFieldError struct {
	field string
}

func (err NoSuchFieldError) Error() string {
	return fmt.Sprintf("no such field: %s", err.field)
}

// Kind is the most basic type descriptor.
type Kind int

// Iteration of valid kinds.
const (
	Invalid Kind = iota

	// Primitive types
	Boolean
	Int
	Uint
	Float
	String

	// Composite types
	Array
	Func
	Map
	Struct
)

// Type is the representation of an expr type.
type Type interface {
	Kind() Kind
	Key() (Type, error)
	Elem() (Type, error)
	Field(name string) (Type, error)
	NumParams() (int, error)
	Param(i int) (Type, error)
	IsVariadic() (bool, error)
	Return() (Type, error)
	ToReflectType() (reflect.Type, error)
}

// primitiveType is the type of primitives.
type primitiveType struct {
	kind Kind
}

// PrimitiveType returns a new primitive type.
func PrimitiveType(k Kind) Type {
	if k >= Array {
		panic("not a primitive kind")
	}
	return primitiveType{k}
}

func (t primitiveType) Kind() Kind                    { return t.kind }
func (primitiveType) Key() (Type, error)              { return nil, ErrInvalidKind }
func (primitiveType) Elem() (Type, error)             { return nil, ErrInvalidKind }
func (primitiveType) Field(name string) (Type, error) { return nil, ErrInvalidKind }
func (primitiveType) NumParams() (int, error)         { return 0, ErrInvalidKind }
func (primitiveType) Param(i int) (Type, error)       { return nil, ErrInvalidKind }
func (primitiveType) IsVariadic() (bool, error)       { return false, ErrInvalidKind }
func (primitiveType) Return() (Type, error)           { return nil, ErrInvalidKind }
func (t primitiveType) ToReflectType() (reflect.Type, error) {
	switch t.kind {
	case Boolean:
		return reflect.TypeOf(false), nil
	case Int:
		return reflect.TypeOf(int64(0)), nil
	case Uint:
		return reflect.TypeOf(uint64(0)), nil
	case Float:
		return reflect.TypeOf(float64(0)), nil
	case String:
		return reflect.TypeOf(""), nil
	default:
		return nil, ErrInvalidKind
	}
}

// arrayType is the type of array-like values.
type arrayType struct {
	elemType Type
}

// ArrayType returns a new array type.
func ArrayType(elem Type) Type {
	return arrayType{elem}
}

func (arrayType) Kind() Kind                      { return Array }
func (arrayType) Key() (Type, error)              { return nil, ErrInvalidKind }
func (t arrayType) Elem() (Type, error)           { return t.elemType, nil }
func (arrayType) Field(name string) (Type, error) { return nil, ErrInvalidKind }
func (arrayType) NumParams() (int, error)         { return 0, ErrInvalidKind }
func (arrayType) Param(i int) (Type, error)       { return nil, ErrInvalidKind }
func (arrayType) IsVariadic() (bool, error)       { return false, ErrInvalidKind }
func (arrayType) Return() (Type, error)           { return nil, ErrInvalidKind }
func (t arrayType) ToReflectType() (reflect.Type, error) {
	val, err := t.elemType.ToReflectType()
	if err != nil {
		return nil, err
	}
	return reflect.SliceOf(val), nil
}

// funcType is the type of function values.
type funcType struct {
	params   []Type
	returns  Type
	variadic bool
}

// FuncType returns a new function type.
func FuncType(params []Type, returns Type, variadic bool) Type {
	return funcType{params, returns, variadic}
}

func (funcType) Kind() Kind                      { return Func }
func (funcType) Key() (Type, error)              { return nil, ErrInvalidKind }
func (funcType) Elem() (Type, error)             { return nil, ErrInvalidKind }
func (funcType) Field(name string) (Type, error) { return nil, ErrInvalidKind }
func (t funcType) NumParams() (int, error)       { return len(t.params), nil }
func (t funcType) Param(i int) (Type, error)     { return t.params[i], nil }
func (t funcType) IsVariadic() (bool, error)     { return t.variadic, nil }
func (t funcType) Return() (Type, error)         { return t.returns, nil }
func (t funcType) ToReflectType() (reflect.Type, error) {
	ret, err := t.returns.ToReflectType()
	if err != nil {
		return nil, err
	}
	out := []reflect.Type{ret}
	nin := len(t.params)
	in := make([]reflect.Type, 0, nin)
	for i := 0; i < nin; i++ {
		param, err := t.params[i].ToReflectType()
		if err != nil {
			return nil, err
		}
		in = append(in, param)
	}
	return reflect.FuncOf(in, out, t.variadic), nil
}

// mapType is the type of maps.
type mapType struct {
	keyType Type
	valType Type
}

// MapType returns a new map type.
func MapType(key Type, val Type) Type {
	return mapType{key, val}
}

func (mapType) Kind() Kind                      { return Map }
func (t mapType) Key() (Type, error)            { return t.keyType, nil }
func (t mapType) Elem() (Type, error)           { return t.valType, nil }
func (mapType) Field(name string) (Type, error) { return nil, ErrInvalidKind }
func (mapType) NumParams() (int, error)         { return 0, ErrInvalidKind }
func (mapType) Param(i int) (Type, error)       { return nil, ErrInvalidKind }
func (mapType) IsVariadic() (bool, error)       { return false, ErrInvalidKind }
func (mapType) Return() (Type, error)           { return nil, ErrInvalidKind }
func (t mapType) ToReflectType() (reflect.Type, error) {
	key, err := t.keyType.ToReflectType()
	if err != nil {
		return nil, err
	}
	val, err := t.valType.ToReflectType()
	if err != nil {
		return nil, err
	}
	return reflect.MapOf(key, val), nil
}

// StructField represents a struct field.
type StructField struct {
	Name string
	Type Type
}

// structType is the type of struct values.
type structType struct {
	fields   []StructField
	fieldMap map[string]StructField
}

// StructType returns a new struct type.
func StructType(fields []StructField) Type {
	fieldMap := map[string]StructField{}
	for _, field := range fields {
		fieldMap[field.Name] = field
	}
	return structType{fields, fieldMap}
}

func (structType) Kind() Kind          { return Struct }
func (structType) Key() (Type, error)  { return nil, ErrInvalidKind }
func (structType) Elem() (Type, error) { return nil, ErrInvalidKind }
func (t structType) Field(name string) (Type, error) {
	field, ok := t.fieldMap[name]
	if !ok {
		return nil, NoSuchFieldError{name}
	}
	return field.Type, nil
}
func (structType) NumParams() (int, error)   { return 0, ErrInvalidKind }
func (structType) Param(i int) (Type, error) { return nil, ErrInvalidKind }
func (structType) IsVariadic() (bool, error) { return false, ErrInvalidKind }
func (structType) Return() (Type, error)     { return nil, ErrInvalidKind }
func (t structType) ToReflectType() (reflect.Type, error) {
	fields := make([]reflect.StructField, 0, len(t.fields))
	for _, field := range t.fields {
		ftype, err := field.Type.ToReflectType()
		if err != nil {
			return nil, err
		}
		fields = append(fields, reflect.StructField{
			Name: field.Name,
			Type: ftype,
		})
	}
	return reflect.StructOf(fields), nil
}
