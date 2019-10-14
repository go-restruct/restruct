package restruct

import (
	"encoding/binary"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

var intType = reflect.TypeOf(int(0))
var boolType = reflect.TypeOf(false)
var strType = reflect.TypeOf(string(""))

func TestFieldsFromStruct(t *testing.T) {
	tests := []struct {
		input  interface{}
		fields fields
	}{
		{
			struct {
				Simple int
			}{},
			fields{
				field{
					Name:       "Simple",
					Index:      0,
					BinaryType: intType,
					NativeType: intType,
					Order:      nil,
					SIndex:     -1,
					TIndex:     -1,
					Skip:       0,
					Trivial:    true,
					BitSize:    0,
					Flags:      0,
				},
			},
		},
		{
			struct {
				Before int
				During string `struct:"-"`
				After  bool
			}{},
			fields{
				field{
					Name:       "Before",
					Index:      0,
					BinaryType: intType,
					NativeType: intType,
					Order:      nil,
					SIndex:     -1,
					TIndex:     -1,
					Skip:       0,
					Trivial:    true,
					BitSize:    0,
					Flags:      0,
				},
				field{
					Name:       "After",
					Index:      2,
					BinaryType: boolType,
					NativeType: boolType,
					Order:      nil,
					SIndex:     -1,
					TIndex:     -1,
					Skip:       0,
					Trivial:    true,
					BitSize:    0,
					Flags:      0,
				},
			},
		},
		{
			struct {
				VariantBool         bool `struct:"variantbool"`
				InvertedBool        bool `struct:"invertedbool"`
				InvertedVariantBool bool `struct:"variantbool,invertedbool"`
			}{},
			fields{
				field{
					Name:       "VariantBool",
					Index:      0,
					BinaryType: boolType,
					NativeType: boolType,
					Order:      nil,
					SIndex:     -1,
					TIndex:     -1,
					Skip:       0,
					Trivial:    true,
					BitSize:    0,
					Flags:      VariantBoolFlag,
				},
				field{
					Name:       "InvertedBool",
					Index:      1,
					BinaryType: boolType,
					NativeType: boolType,
					Order:      nil,
					SIndex:     -1,
					TIndex:     -1,
					Skip:       0,
					Trivial:    true,
					BitSize:    0,
					Flags:      InvertedBoolFlag,
				},
				field{
					Name:       "InvertedVariantBool",
					Index:      2,
					BinaryType: boolType,
					NativeType: boolType,
					Order:      nil,
					SIndex:     -1,
					TIndex:     -1,
					Skip:       0,
					Trivial:    true,
					BitSize:    0,
					Flags:      VariantBoolFlag | InvertedBoolFlag,
				},
			},
		},
		{
			struct {
				FixedStr string `struct:"[64]byte,skip=4"`
				LSBInt   int    `struct:"uint32,little"`
			}{},
			fields{
				field{
					Name:       "FixedStr",
					Index:      0,
					BinaryType: reflect.TypeOf([64]byte{}),
					NativeType: strType,
					Order:      nil,
					SIndex:     -1,
					TIndex:     -1,
					Skip:       4,
					Trivial:    true,
					BitSize:    0,
					Flags:      0,
				},
				field{
					Name:       "LSBInt",
					Index:      1,
					BinaryType: reflect.TypeOf(uint32(0)),
					NativeType: intType,
					Order:      binary.LittleEndian,
					SIndex:     -1,
					TIndex:     -1,
					Skip:       0,
					Trivial:    true,
					BitSize:    0,
					Flags:      0,
				},
			},
		},
		{
			struct {
				NumColors int32 `struct:"sizeof=Colors"`
				Colors    [][4]uint8
			}{},
			fields{
				field{
					Name:       "NumColors",
					Index:      0,
					BinaryType: reflect.TypeOf(int32(0)),
					NativeType: reflect.TypeOf(int32(0)),
					SIndex:     -1,
					TIndex:     1,
					Skip:       0,
					Trivial:    true,
				},
				field{
					Name:       "Colors",
					Index:      1,
					BinaryType: reflect.TypeOf([][4]uint8{}),
					NativeType: reflect.TypeOf([][4]uint8{}),
					SIndex:     0,
					TIndex:     -1,
					Skip:       0,
					Trivial:    false,
				},
			},
		},
		{
			struct {
				NumColors int32
				Colors    [][4]uint8 `struct:"sizefrom=NumColors"`
			}{},
			fields{
				field{
					Name:       "NumColors",
					Index:      0,
					BinaryType: reflect.TypeOf(int32(0)),
					NativeType: reflect.TypeOf(int32(0)),
					SIndex:     -1,
					TIndex:     1,
					Skip:       0,
					Trivial:    true,
				},
				field{
					Name:       "Colors",
					Index:      1,
					BinaryType: reflect.TypeOf([][4]uint8{}),
					NativeType: reflect.TypeOf([][4]uint8{}),
					SIndex:     0,
					TIndex:     -1,
					Skip:       0,
					Trivial:    false,
				},
			},
		},
	}

	for _, test := range tests {
		fields := fieldsFromStruct(reflect.TypeOf(test.input))
		assert.Equal(t, test.fields, fields)
	}
}

func TestFieldsFromNonStructPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Non-struct did not panic.")
		}
	}()
	fieldsFromStruct(reflect.TypeOf(0))
}

func TestFieldsFromBrokenSizeOf(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Error("Broken struct did not panic.")
		}
		assert.Equal(t, "couldn't find SizeOf field Nonexistant", r.(error).Error())
	}()

	badSize := struct {
		Test int64 `struct:"sizeof=Nonexistant"`
	}{}
	fieldsFromStruct(reflect.TypeOf(badSize))
}

func TestFieldsFromBrokenSizeFrom(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Error("Broken struct did not panic.")
		}
		assert.Equal(t, "couldn't find SizeFrom field Nonexistant", r.(error).Error())
	}()

	badSize := struct {
		Test string `struct:"sizefrom=Nonexistant"`
	}{}
	fieldsFromStruct(reflect.TypeOf(badSize))
}

func TestIsTypeTrivial(t *testing.T) {
	tests := []struct {
		input   interface{}
		trivial bool
	}{
		{int8(0), true},
		{int16(0), true},
		{int32(0), true},
		{int64(0), true},
		{[0]int8{}, true},
		{[]int8{}, false},
		{struct{}{}, true},
		{struct{ int8 }{}, true},
		{struct{ A []int8 }{[]int8{}}, false},
		{struct{ A [0]int8 }{[0]int8{}}, true},
		{(*interface{})(nil), false},
	}

	for _, test := range tests {
		assert.Equal(t, test.trivial, isTypeTrivial(reflect.TypeOf(test.input)))
	}
}

func BenchmarkFieldsFromStruct(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fieldsFromStruct(reflect.TypeOf(TestStruct{}))
	}
}
