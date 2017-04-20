package restruct

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"sync"
)

// Sizer is a type which has a defined size in binary. The SizeOf function
// returns how many bytes the type will consume in memory. This is used during
// encoding for allocation and therefore must equal the exact number of bytes
// the encoded form needs. You may use a pointer receiver even if the type is
// used by value.
type Sizer interface {
	SizeOf() int
}

// field represents a structure field, similar to reflect.StructField.
type field struct {
	Name    string
	Index   int
	Type    reflect.Type
	DefType reflect.Type
	Order   binary.ByteOrder
	SIndex  int
	Skip    int
	Trivial bool
	BitSize uint8
}

// fields represents a structure.
type fields []field

var fieldCache = map[reflect.Type][]field{}
var cacheMutex = sync.RWMutex{}

// Elem constructs a transient field representing an element of an array, slice,
// or pointer.
func (f *field) Elem() field {
	// Special cases for string types, grumble grumble.
	t := f.Type
	if t.Kind() == reflect.String {
		t = reflect.TypeOf([]byte{})
	}

	dt := f.DefType
	if dt.Kind() == reflect.String {
		dt = reflect.TypeOf([]byte{})
	}

	return field{
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

// fieldFromType returns a field from a reflected type.
func fieldFromType(typ reflect.Type) field {
	return field{
		Index:   -1,
		Type:    typ,
		DefType: typ,
		Order:   nil,
		SIndex:  -1,
		Skip:    0,
		Trivial: isTypeTrivial(typ),
	}
}

// fieldsFromStruct returns a slice of fields for binary packing and unpacking.
func fieldsFromStruct(typ reflect.Type) (result fields) {
	if typ.Kind() != reflect.Struct {
		panic(fmt.Errorf("tried to get fields from non-struct type %s", typ.Kind().String()))
	}

	count := typ.NumField()

	for i := 0; i < count; i++ {
		val := typ.Field(i)

		// Skip unexported names (except _)
		if val.PkgPath != "" && val.Name != "_" {
			continue
		}

		// Parse struct tag
		opts := mustParseTag(val.Tag.Get("struct"))
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

		result = append(result, field{
			Name:    val.Name,
			Index:   i,
			Type:    ftyp,
			DefType: val.Type,
			Order:   opts.Order,
			SIndex:  sindex,
			Skip:    opts.Skip,
			Trivial: isTypeTrivial(ftyp),
			BitSize: opts.BitSize,
		})
	}

	return
}

func cachedFieldsFromStruct(typ reflect.Type) (result fields) {
	cacheMutex.RLock()
	result, ok := fieldCache[typ]
	cacheMutex.RUnlock()

	if ok {
		return
	}

	result = fieldsFromStruct(typ)

	cacheMutex.Lock()
	fieldCache[typ] = result
	cacheMutex.Unlock()

	return
}

// isTypeTrivial determines if a given type is constant-size.
func isTypeTrivial(typ reflect.Type) bool {
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
		return isTypeTrivial(typ.Elem())
	case reflect.Struct:
		for _, field := range cachedFieldsFromStruct(typ) {
			if !isTypeTrivial(field.Type) {
				return false
			}
		}
		return true
	default:
		return false
	}
}

func (f *field) sizer(v reflect.Value) (Sizer, bool) {
	if s, ok := v.Interface().(Sizer); ok {
		return s, true
	}

	if !v.CanAddr() {
		return nil, false
	}

	if s, ok := v.Addr().Interface().(Sizer); ok {
		return s, true
	}

	return nil, false
}

// SizeOf determines what the binary size of the field should be.
func (f *field) SizeOf(val reflect.Value) (size int) {
	if f.Name != "_" {
		if s, ok := f.sizer(val); ok {
			return s.SizeOf()
		}
	} else {
		// Non-trivial, unnamed fields do not make sense. You can't set a field
		// with no name, so the elements can't possibly differ.
		// N.B.: Though skip will still work, use struct{} instead for skip.
		if !isTypeTrivial(val.Type()) {
			return f.Skip
		}
	}

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
		switch f.DefType.Kind() {
		case reflect.Slice, reflect.String, reflect.Array, reflect.Ptr:
			alen = val.Len()
		default:
			return 0
		}
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
		var bitSize uint64
		for _, field := range cachedFieldsFromStruct(f.Type) {
			if field.BitSize != 0 {
				bitSize += uint64(field.BitSize)
			} else {
				size += field.SizeOf(val.Field(field.Index))
			}
		}
		size += int(bitSize / 8)
		if bitSize%8 > 0 {
			size++
		}
		return size
	default:
		return 0
	}
}

// SizeOf returns the size of a struct.
func (fields fields) SizeOf(val reflect.Value) (size int) {
	for _, field := range fields {
		size += field.SizeOf(val.Field(field.Index))
	}
	return
}
