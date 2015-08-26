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
	Type    reflect.Type
	DefType reflect.Type
	Order   binary.ByteOrder
	SIndex  int
	Skip    int
	Trivial bool
}

// Fields represents a structure.
type Fields []Field

var fieldCache = map[reflect.Type][]Field{}

// Elem constructs a transient field representing an element of an array, slice,
// or pointer.
func (f *Field) Elem() Field {
	// Special cases for string types, grumble grumble.
	t := f.Type
	if t.Kind() == reflect.String {
		t = reflect.TypeOf([]byte{})
	}

	dt := f.DefType
	if dt.Kind() == reflect.String {
		dt = reflect.TypeOf([]byte{})
	}

	return Field{
		Name:    "*" + f.Name,
		Index:   -1,
		Type:    t.Elem(),
		DefType: dt.Elem(),
		Order:   f.Order,
		SIndex:  -1,
		Skip:    0,
		Trivial: f.Trivial,
	}
}

// FieldFromType returns a field from a reflected type.
func FieldFromType(typ reflect.Type) Field {
	return Field{
		Index:   -1,
		Type:    typ,
		DefType: typ,
		Order:   nil,
		SIndex:  -1,
		Skip:    0,
		Trivial: IsTypeTrivial(typ),
	}
}

// FieldsFromStruct returns a slice of fields for binary packing and unpacking.
func FieldsFromStruct(typ reflect.Type) (result Fields) {
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

		// SizeOf
		sindex := -1
		if opts.SizeOf != "" {
			count := typ.NumField()
			for j := i + 1; j < count; j++ {
				val := typ.Field(j)
				if opts.SizeOf == val.Name {
					sindex = j
				}
			}
			if sindex == -1 {
				panic(fmt.Errorf("couldn't find SizeOf field %s", opts.SizeOf))
			}
		}

		result = append(result, Field{
			Name:    val.Name,
			Index:   i,
			Type:    ftyp,
			DefType: val.Type,
			Order:   opts.Order,
			SIndex:  sindex,
			Skip:    opts.Skip,
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
	alen := 1
	switch f.Type.Kind() {
	case reflect.Int8, reflect.Uint8:
		return 1 + f.Skip
	case reflect.Int16, reflect.Uint16:
		return 2 + f.Skip
	case reflect.Int, reflect.Int32,
		reflect.Uint, reflect.Uint32,
		reflect.Bool, reflect.Float32:
		return 4 + f.Skip
	case reflect.Int64, reflect.Uint64,
		reflect.Float64, reflect.Complex64:
		return 8 + f.Skip
	case reflect.Complex128:
		return 16 + f.Skip
	case reflect.Slice, reflect.String:
		alen = val.Len()
		fallthrough
	case reflect.Array, reflect.Ptr:
		size += f.Skip

		// If array type, get length from type.
		if f.Type.Kind() == reflect.Array {
			alen = f.Type.Len()
		}

		// Optimization: if the array/slice is empty, bail now.
		if alen == 0 {
			return size
		}

		// Optimization: if the type is trivial, we only need to check the
		// first element.
		switch f.DefType.Kind() {
		case reflect.Slice, reflect.String, reflect.Array, reflect.Ptr:
			elem := f.Elem()
			if f.Trivial {
				size += elem.SizeOf(reflect.Zero(f.Type.Elem())) * alen
			} else {
				for i := 0; i < alen; i++ {
					size += elem.SizeOf(val.Index(i))
				}
			}
		}
		return size
	case reflect.Struct:
		size += f.Skip
		for _, field := range FieldsFromStruct(f.Type) {
			size += field.SizeOf(val.Field(field.Index))
		}
		return size
	default:
		return 0
	}
}

// SizeOf returns the size of a struct.
func (fields Fields) SizeOf(val reflect.Value) (size int) {
	for _, field := range fields {
		size += field.SizeOf(val.Field(field.Index))
	}
	return
}
