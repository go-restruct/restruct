package restruct

import (
	"encoding/binary"
	"fmt"
	"reflect"
)

// Field represents a structure field, similar to reflect.StructField.
type Field struct {
	Name    string
	CanSet  bool
	Type    reflect.Type
	DefType reflect.Type
	Order   binary.ByteOrder
}

// FieldsFromStruct returns a slice of fields for binary packing and unpacking.
func FieldsFromStruct(typ reflect.Type) (result []Field) {
	if typ.Kind() != reflect.Struct {
		panic(fmt.Errorf("tried to get fields from non-struct type %s", typ.Kind().String()))
	}

	count := typ.NumField()
	for i := 0; i < count; i++ {
		val := typ.Field(i)

		// Parse struct tag
		opts := MustParseTag(val.Tag.Get("struct"))
		if opts.Ignore {
			continue
		}

		// Derive type
		ftyp := val.Type
		if opts.Type != nil {
			ftyp = opts.Type
		}

		result = append(result, Field{
			Name:    val.Name,
			CanSet:  true,
			Type:    ftyp,
			DefType: val.Type,
			Order:   opts.Order,
		})
	}

	return result
}

// IsTypeTrivial determines if a given type is constant-size.
func IsTypeTrivial(typ reflect.Type) bool {
	switch typ.Kind() {
	case reflect.Bool,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr,
		reflect.Float32,
		reflect.Float64,
		reflect.Complex64,
		reflect.Complex128:
		return true
	case reflect.Array, reflect.Ptr, reflect.Slice:
		return IsTypeTrivial(typ.Elem())
	case reflect.Struct:
		for _, field := range FieldsFromStruct(typ) {
			if !IsTypeTrivial(field.Type) {
				return false
			}
		}
		return true
	default:
		return true
	}
}

// Elem synthesizes a field from an element of a slice, array, or pointer.
// If Field is not a slice, array, or pointer, FieldFromArray panics.
func (f *Field) Elem() *Field {
	return &Field{
		Name:    "*" + f.Name,
		CanSet:  f.CanSet,
		Type:    f.Type.Elem(),
		DefType: f.Type.Elem(),
		Order:   f.Order,
	}
}

// SizeOf determines what the binary size of the field should be.
func (f *Field) SizeOf(val reflect.Value) {

}
