package restruct

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"sync"
)

// FieldFlags is a type for flags that can be applied to fields individually.
type FieldFlags uint64

const (
	// VariantBoolFlag causes the true value of a boolean to be ~0 instead of
	// just 1 (all bits are set.) This emulates the behavior of VARIANT_BOOL.
	VariantBoolFlag FieldFlags = 1 << iota

	// InvertedBoolFlag causes the true and false states of a boolean to be
	// flipped in binary.
	InvertedBoolFlag
)

// Sizer is a type which has a defined size in binary. The SizeOf function
// returns how many bytes the type will consume in memory. This is used during
// encoding for allocation and therefore must equal the exact number of bytes
// the encoded form needs. You may use a pointer receiver even if the type is
// used by value.
type Sizer interface {
	SizeOf() int
}

// BitSizer is an interface for types that need to specify their own size in
// bit-level granularity. It has the same effect as Sizer.
type BitSizer interface {
	BitSize() int
}

// field represents a structure field, similar to reflect.StructField.
type field struct {
	Name       string
	Index      int
	BinaryType reflect.Type
	NativeType reflect.Type
	Order      binary.ByteOrder
	SIndex     int // Index of size field for a slice/string.
	TIndex     int // Index of target of sizeof field.
	Skip       int
	Trivial    bool
	BitSize    uint8
	Flags      FieldFlags
}

// fields represents a structure.
type fields []field

var fieldCache = map[reflect.Type][]field{}
var cacheMutex = sync.RWMutex{}

// Elem constructs a transient field representing an element of an array, slice,
// or pointer.
func (f *field) Elem() field {
	// Special cases for string types, grumble grumble.
	t := f.BinaryType
	if t.Kind() == reflect.String {
		t = reflect.TypeOf([]byte{})
	}

	dt := f.NativeType
	if dt.Kind() == reflect.String {
		dt = reflect.TypeOf([]byte{})
	}

	return field{
		Name:       "*" + f.Name,
		Index:      -1,
		BinaryType: t.Elem(),
		NativeType: dt.Elem(),
		Order:      f.Order,
		TIndex:     -1,
		SIndex:     -1,
		Skip:       0,
		Trivial:    f.Trivial,
	}
}

// fieldFromType returns a field from a reflected type.
func fieldFromType(typ reflect.Type) field {
	return field{
		Index:      -1,
		BinaryType: typ,
		NativeType: typ,
		Order:      nil,
		TIndex:     -1,
		SIndex:     -1,
		Skip:       0,
		Trivial:    isTypeTrivial(typ),
	}
}

// fieldsFromStruct returns a slice of fields for binary packing and unpacking.
func fieldsFromStruct(typ reflect.Type) (result fields) {
	if typ.Kind() != reflect.Struct {
		panic(fmt.Errorf("tried to get fields from non-struct type %s", typ.Kind().String()))
	}

	count := typ.NumField()

	sizeOfMap := map[string]int{}

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
		tindex := -1
		if j, ok := sizeOfMap[val.Name]; ok {
			sindex = j
			result[sindex].TIndex = i
			delete(sizeOfMap, val.Name)
		} else if opts.SizeOf != "" {
			sizeOfMap[opts.SizeOf] = i
		}

		// SizeFrom
		if opts.SizeFrom != "" {
			for j := 0; j < i; j++ {
				val := result[j]
				if opts.SizeFrom == val.Name {
					sindex = j
					result[sindex].TIndex = i
				}
			}
			if sindex == -1 {
				panic(fmt.Errorf("couldn't find SizeFrom field %s", opts.SizeFrom))
			}
		}

		// Flags
		flags := FieldFlags(0)
		if opts.VariantBoolFlag {
			flags |= VariantBoolFlag
		}
		if opts.InvertedBoolFlag {
			flags |= InvertedBoolFlag
		}

		result = append(result, field{
			Name:       val.Name,
			Index:      i,
			BinaryType: ftyp,
			NativeType: val.Type,
			Order:      opts.Order,
			SIndex:     sindex,
			TIndex:     tindex,
			Skip:       opts.Skip,
			Trivial:    isTypeTrivial(ftyp),
			BitSize:    opts.BitSize,
			Flags:      flags,
		})
	}

	for fieldName := range sizeOfMap {
		panic(fmt.Errorf("couldn't find SizeOf field %s", fieldName))
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
			if !isTypeTrivial(field.BinaryType) {
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

func (f *field) bitSizer(v reflect.Value) (BitSizer, bool) {
	if s, ok := v.Interface().(BitSizer); ok {
		return s, true
	}

	if !v.CanAddr() {
		return nil, false
	}

	if s, ok := v.Addr().Interface().(BitSizer); ok {
		return s, true
	}

	return nil, false
}

func (f *field) bitSizeUsingInterface(val reflect.Value) (int, bool) {
	if s, ok := f.bitSizer(val); ok {
		return s.BitSize(), true
	}

	if s, ok := f.sizer(val); ok {
		return s.SizeOf() * 8, true
	}

	return 0, false
}

// SizeOfBits determines the encoded size of a field in bits.
func (f *field) SizeOfBits(val reflect.Value) (size int) {
	skipBits := f.Skip * 8

	if f.Name != "_" {
		if s, ok := f.bitSizeUsingInterface(val); ok {
			return s
		}
	} else {
		// Non-trivial, unnamed fields do not make sense. You can't set a field
		// with no name, so the elements can't possibly differ.
		// N.B.: Though skip will still work, use struct{} instead for skip.
		if !isTypeTrivial(val.Type()) {
			return skipBits
		}
	}

	alen := 1
	switch f.BinaryType.Kind() {
	case reflect.Int8, reflect.Uint8, reflect.Bool:
		return 8 + skipBits
	case reflect.Int16, reflect.Uint16:
		return 16 + skipBits
	case reflect.Int, reflect.Int32,
		reflect.Uint, reflect.Uint32,
		reflect.Float32:
		return 32 + skipBits
	case reflect.Int64, reflect.Uint64,
		reflect.Float64, reflect.Complex64:
		return 64 + skipBits
	case reflect.Complex128:
		return 128 + skipBits
	case reflect.Slice, reflect.String:
		switch f.NativeType.Kind() {
		case reflect.Slice, reflect.String, reflect.Array, reflect.Ptr:
			alen = val.Len()
		default:
			return 0
		}
		fallthrough
	case reflect.Array, reflect.Ptr:
		size += skipBits

		// If array type, get length from type.
		if f.BinaryType.Kind() == reflect.Array {
			alen = f.BinaryType.Len()
		}

		// Optimization: if the array/slice is empty, bail now.
		if alen == 0 {
			return size
		}

		// Optimization: if the type is trivial, we only need to check the
		// first element.
		switch f.NativeType.Kind() {
		case reflect.Slice, reflect.String, reflect.Array, reflect.Ptr:
			elem := f.Elem()
			if f.Trivial {
				size += elem.SizeOfBits(reflect.Zero(f.BinaryType.Elem())) * alen
			} else {
				for i := 0; i < alen; i++ {
					size += elem.SizeOfBits(val.Index(i))
				}
			}
		}
		return size
	case reflect.Struct:
		size += skipBits
		for _, field := range cachedFieldsFromStruct(f.BinaryType) {
			if field.BitSize != 0 {
				size += int(field.BitSize)
			} else {
				size += field.SizeOfBits(val.Field(field.Index))
			}
		}
		return size
	default:
		return 0
	}
}

// SizeOfBytes returns the effective size in bytes, for the few cases where
// byte sizes are needed.
func (f *field) SizeOfBytes(val reflect.Value) (size int) {
	return (f.SizeOfBits(val) + 7) / 8
}

// SizeOf returns the size of a struct.
func (fields fields) SizeOfBits(val reflect.Value) (size int) {
	for _, field := range fields {
		size += field.SizeOfBits(val.Field(field.Index))
	}
	return
}
