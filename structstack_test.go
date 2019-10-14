package restruct

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
		{int8(0), 8},
		{int16(0), 16},
		{int32(0), 32},
		{int64(0), 64},
		{uint8(0), 8},
		{uint16(0), 16},
		{uint32(0), 32},
		{uint64(0), 64},
		{float32(0), 32},
		{float64(0), 64},
		{complex64(0), 64},
		{complex128(0), 128},
		{[0]int8{}, 0},
		{[1]int8{1}, 8},
		{[]int8{1, 2}, 16},
		{[]int32{1, 2}, 64},
		{[2][3]int8{}, 48},
		{struct{}{}, 0},
		{struct{ A int8 }{}, 8},
		{struct{ A []int8 }{[]int8{}}, 0},
		{struct{ A [0]int8 }{[0]int8{}}, 0},
		{struct{ A []int8 }{[]int8{1}}, 8},
		{struct{ A [1]int8 }{[1]int8{1}}, 8},
		{TestStruct{}, 17040},
		{interface{}(struct{}{}), 0},
		{struct{ Test interface{} }{}, 0},

		// Unexported fields test
		{struct{ a int8 }{}, 0},
		{struct{ a []int8 }{[]int8{}}, 0},
		{struct{ a [0]int8 }{[0]int8{}}, 0},
		{struct{ a []int8 }{[]int8{1}}, 0},
		{struct{ a [1]int8 }{[1]int8{1}}, 0},

		// Trivial unnamed fields test
		{struct{ _ [1]int8 }{}, 8},
		{struct {
			_ [1]int8 `struct:"skip=4"`
		}{}, 40},

		// Non-trivial unnamed fields test
		{struct{ _ []interface{} }{}, 0},
		{struct{ _ [1]interface{} }{}, 0},
		{struct {
			_ [1]interface{} `struct:"skip=4"`
		}{}, 32},
		{struct {
			_ [4]struct {
				_ [4]struct{} `struct:"skip=4"`
			} `struct:"skip=4"`
		}{}, 160},
		{struct{ T string }{"yeehaw"}, 48},

		// Byte-misaligned structures
		{[10]struct {
			_ int8 `struct:"uint8:1"`
		}{}, 10},
		{[4]struct {
			_ bool `struct:"uint8:1,variantbool"`
			_ int8 `struct:"uint8:4"`
			_ int8 `struct:"uint8:4"`
		}{}, 36},
	}

	ss := structstack{}
	for _, test := range tests {
		field := fieldFromType(reflect.TypeOf(test.input))
		assert.Equal(t, test.size, ss.fieldbits(field, reflect.ValueOf(test.input)),
			"bad size for input: %#v", test.input)
	}
}

var (
	simpleFields  = fieldsFromStruct(reflect.TypeOf(TestElem{}))
	complexFields = fieldsFromStruct(reflect.TypeOf(TestStruct{}))
)

func TestSizeOfFields(t *testing.T) {
	ss := structstack{}
	assert.Equal(t, 72, ss.fieldsbits(simpleFields, reflect.ValueOf(TestElem{})))
	assert.Equal(t, 17040, ss.fieldsbits(complexFields, reflect.ValueOf(TestStruct{})))
}

func BenchmarkSizeOfSimple(b *testing.B) {
	ss := structstack{}
	for i := 0; i < b.N; i++ {
		ss.fieldsbits(simpleFields, reflect.ValueOf(TestElem{}))
	}
}

func BenchmarkSizeOfComplex(b *testing.B) {
	ss := structstack{}
	for i := 0; i < b.N; i++ {
		ss.fieldsbits(complexFields, reflect.ValueOf(TestStruct{}))
	}
}
