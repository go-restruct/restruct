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
				field{"Simple", 0, intType, intType, nil, -1, 0, true, 0},
			},
		},
		{
			struct {
				Before int
				During string `struct:"-"`
				After  bool
			}{},
			fields{
				field{"Before", 0, intType, intType, nil, -1, 0, true, 0},
				field{"After", 2, boolType, boolType, nil, -1, 0, true, 0},
			},
		},
		{
			struct {
				FixedStr string `struct:"[64]byte,skip=4"`
				LSBInt   int    `struct:"uint32,little"`
			}{},
			fields{
				field{"FixedStr", 0, reflect.TypeOf([64]byte{}), strType, nil, -1, 4, true, 0},
				field{"LSBInt", 1, reflect.TypeOf(uint32(0)), intType, binary.LittleEndian, -1, 0, true, 0},
			},
		},
		{
			struct {
				NumColors int32 `struct:"sizeof=Colors"`
				Colors    [][4]uint8
			}{},
			fields{
				field{
					Name:    "NumColors",
					Index:   0,
					Type:    reflect.TypeOf(int32(0)),
					DefType: reflect.TypeOf(int32(0)),
					SIndex:  1,
					Skip:    0,
					Trivial: true,
				},
				field{
					Name:    "Colors",
					Index:   1,
					Type:    reflect.TypeOf([][4]uint8{}),
					DefType: reflect.TypeOf([][4]uint8{}),
					SIndex:  -1,
					Skip:    0,
					Trivial: false,
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

func TestFieldsFromBrokenStruct(t *testing.T) {
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

type TestElem struct {
	Test1 int64
	Test2 int8
}

type TestStruct struct {
	Sub [10]struct {
		Sub2 struct {
			Size  int `struct:"uint32,sizeof=Elems"`
			Elems []TestElem
		} `struct:"skip=4"`
	} `struct:"skip=2"`
	Numbers  [128]int64
	Numbers2 []float64 `struct:"[256]float32"`
}

func TestSizeOf(t *testing.T) {
	tests := []struct {
		input interface{}
		size  int
	}{
		{int8(0), 1},
		{int16(0), 2},
		{int32(0), 4},
		{int64(0), 8},
		{uint8(0), 1},
		{uint16(0), 2},
		{uint32(0), 4},
		{uint64(0), 8},
		{float32(0), 4},
		{float64(0), 8},
		{complex64(0), 8},
		{complex128(0), 16},
		{[0]int8{}, 0},
		{[1]int8{1}, 1},
		{[]int8{1, 2}, 2},
		{[]int32{1, 2}, 8},
		{[2][3]int8{}, 6},
		{struct{}{}, 0},
		{struct{ A int8 }{}, 1},
		{struct{ A []int8 }{[]int8{}}, 0},
		{struct{ A [0]int8 }{[0]int8{}}, 0},
		{struct{ A []int8 }{[]int8{1}}, 1},
		{struct{ A [1]int8 }{[1]int8{1}}, 1},
		{TestStruct{}, 2130},
		{interface{}(struct{}{}), 0},
		{struct{ Test interface{} }{}, 0},

		// Unexported fields test
		{struct{ a int8 }{}, 0},
		{struct{ a []int8 }{[]int8{}}, 0},
		{struct{ a [0]int8 }{[0]int8{}}, 0},
		{struct{ a []int8 }{[]int8{1}}, 0},
		{struct{ a [1]int8 }{[1]int8{1}}, 0},

		// Trivial unnamed fields test
		{struct{ _ [1]int8 }{}, 1},
		{struct {
			_ [1]int8 `struct:"skip=4"`
		}{}, 5},

		// Non-trivial unnamed fields test
		{struct{ _ []interface{} }{}, 0},
		{struct{ _ [1]interface{} }{}, 0},
		{struct {
			_ [1]interface{} `struct:"skip=4"`
		}{}, 4},
		{struct {
			_ [4]struct {
				_ [4]struct{} `struct:"skip=4"`
			} `struct:"skip=4"`
		}{}, 20},
		{struct{ T string }{"yeehaw"}, 6},
	}

	for _, test := range tests {
		field := fieldFromType(reflect.TypeOf(test.input))
		assert.Equal(t, test.size, field.SizeOf(reflect.ValueOf(test.input)),
			"bad size for input: %#v", test.input)
	}
}

var simpleFields fields
var complexFields fields

func init() {
	RegisterArrayType([256]float32{})
	simpleFields = fieldsFromStruct(reflect.TypeOf(TestElem{}))
	complexFields = fieldsFromStruct(reflect.TypeOf(TestStruct{}))
}

func TestSizeOfFields(t *testing.T) {
	assert.Equal(t, simpleFields.SizeOf(reflect.ValueOf(TestElem{})), 9)
	assert.Equal(t, complexFields.SizeOf(reflect.ValueOf(TestStruct{})), 2130)
}

func BenchmarkFieldsFromStruct(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fieldsFromStruct(reflect.TypeOf(TestStruct{}))
	}
}

func BenchmarkSizeOfSimple(b *testing.B) {
	for i := 0; i < b.N; i++ {
		simpleFields.SizeOf(reflect.ValueOf(TestElem{}))
	}
}

func BenchmarkSizeOfComplex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		complexFields.SizeOf(reflect.ValueOf(TestStruct{}))
	}
}
