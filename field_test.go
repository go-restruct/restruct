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
		fields Fields
	}{
		{
			struct {
				Simple int
			}{},
			Fields{
				Field{"Simple", 0, true, intType, intType, nil, 0, true},
			},
		},
		{
			struct {
				Before int
				During string `struct:"-"`
				After  bool
			}{},
			Fields{
				Field{"Before", 0, true, intType, intType, nil, 0, true},
				Field{"After", 2, true, boolType, boolType, nil, 0, true},
			},
		},
		{
			struct {
				FixedStr string `struct:"[64]byte,skip=4"`
				LSBInt   int    `struct:"uint32,little"`
			}{},
			Fields{
				Field{"FixedStr", 0, true, reflect.TypeOf([64]byte{}), strType, nil, 4, true},
				Field{"LSBInt", 1, true, reflect.TypeOf(uint32(0)), intType, binary.LittleEndian, 0, true},
			},
		},
	}

	for _, test := range tests {
		fields := FieldsFromStruct(reflect.TypeOf(test.input))
		assert.Equal(t, fields, test.fields)
	}
}

func TestFieldsFromNonStructPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Non-struct did not panic.")
		}
	}()
	FieldsFromStruct(reflect.TypeOf(0))
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
		{struct{ a []int8 }{[]int8{}}, false},
		{struct{ a [0]int8 }{[0]int8{}}, true},
		{(*interface{})(nil), false},
	}

	for _, test := range tests {
		assert.Equal(t, test.trivial, IsTypeTrivial(reflect.TypeOf(test.input)))
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
		{struct{ int8 }{}, 1},
		{struct{ a []int8 }{[]int8{}}, 0},
		{struct{ a [0]int8 }{[0]int8{}}, 0},
		{struct{ a []int8 }{[]int8{1}}, 1},
		{struct{ a [1]int8 }{[1]int8{1}}, 1},
		{TestStruct{}, 2130},
		{interface{}(struct{}{}), 0},
	}

	for _, test := range tests {
		field := FieldFromType(reflect.TypeOf(test.input))
		assert.Equal(t, test.size, field.SizeOf(reflect.ValueOf(test.input)),
			"bad size for input: %#v", test.input)
	}
}

var simpleFields = FieldsFromStruct(reflect.TypeOf(TestElem{}))
var complexFields = FieldsFromStruct(reflect.TypeOf(TestStruct{}))

func TestSizeOfFields(t *testing.T) {
	assert.Equal(t, simpleFields.SizeOf(reflect.ValueOf(TestElem{})), 9)
	assert.Equal(t, complexFields.SizeOf(reflect.ValueOf(TestStruct{})), 2130)
}

func BenchmarkFieldsFromStruct(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FieldsFromStruct(reflect.TypeOf(TestStruct{}))
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
