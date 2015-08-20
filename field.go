package restruct

import (
	"encoding/binary"
	"fmt"
	"reflect"
)

// Field represents a structure field, similar to reflect.StructField.
type Field struct {
	Name    string
	Index   int
	CanSet  bool
	Type    reflect.Type
	DefType reflect.Type
	Order   binary.ByteOrder
	Trivial bool
}

var fieldCache = map[reflect.Type][]Field{}

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
			Index:   i,
			CanSet:  true,
			Type:    ftyp,
			DefType: val.Type,
			Order:   opts.Order,
			Trivial: IsTypeTrivial(ftyp),
		})
	}

	fieldCache[typ] = result
	return
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
	case reflect.Array, reflect.Ptr:
		return IsTypeTrivial(typ.Elem())
	case reflect.Struct:
		for _, field := range FieldsFromStruct(typ) {
			if !IsTypeTrivial(field.Type) {
				return false
			}
		}
		return true
	default:
		return false
	}
}

// SizeOf determines what the binary size of the field should be.
func (f *Field) SizeOf(val reflect.Value) (size int) {
	switch f.Type.Kind() {
	case reflect.Int8, reflect.Uint8:
		return 1
	case reflect.Int16, reflect.Uint16:
		return 2
	case reflect.Int, reflect.Int32,
		reflect.Uint, reflect.Uint32,
		reflect.Bool, reflect.Float32:
		return 4
	case reflect.Int64, reflect.Uint64,
		reflect.Float64, reflect.Complex64:
		return 8
	case reflect.Complex128:
		return 16
	case reflect.Array, reflect.Ptr, reflect.Slice:
		// Optimization: if the type is trivial, we only need to check the
		// first element.
		if f.Trivial {
			size += f.SizeOf(val.Index(0)) * val.Len()
		} else {
			count := val.Len()
			for i := 0; i < count; i++ {
				size += f.SizeOf(val.Index(i)) * val.Len()
			}
		}
		return size
	case reflect.Struct:
		if f.Trivial {
			for _, field := range FieldsFromStruct(f.Type) {
				size += field.SizeOf(val.Field(field.Index))
			}
		}
		return size
	default:
		return 0
	}
}
