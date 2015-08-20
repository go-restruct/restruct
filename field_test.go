package restruct

import (
	"encoding/binary"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type FieldsFromStructTestCase struct {
	input  interface{}
	fields []Field
}

var intType = reflect.TypeOf(int(0))
var boolType = reflect.TypeOf(false)
var strType = reflect.TypeOf(string(""))

var fieldsFromStructTestCases = []FieldsFromStructTestCase{
	FieldsFromStructTestCase{
		struct {
			Simple int
		}{},
		[]Field{
			Field{"Simple", 0, true, intType, intType, nil, true},
		},
	},
	FieldsFromStructTestCase{
		struct {
			Before int
			During string `struct:"-"`
			After  bool
		}{},
		[]Field{
			Field{"Before", 0, true, intType, intType, nil, true},
			Field{"After", 2, true, boolType, boolType, nil, true},
		},
	},
	FieldsFromStructTestCase{
		struct {
			FixedStr string `struct:"[64]byte"`
			LSBInt   int    `struct:"uint32,little"`
		}{},
		[]Field{
			Field{"FixedStr", 0, true, reflect.TypeOf([64]byte{}), strType, nil, true},
			Field{"LSBInt", 1, true, reflect.TypeOf(uint32(0)), intType, binary.LittleEndian, true},
		},
	},
}

func TestFieldsFromStruct(t *testing.T) {
	for _, data := range fieldsFromStructTestCases {
		fields := FieldsFromStruct(reflect.TypeOf(data.input))
		assert.Equal(t, fields, data.fields)
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

type TypeTrivialTestCase struct {
	input   interface{}
	trivial bool
}

var typeTrivialTestCases = []TypeTrivialTestCase{
	TypeTrivialTestCase{int8(0), true},
	TypeTrivialTestCase{int16(0), true},
	TypeTrivialTestCase{int32(0), true},
	TypeTrivialTestCase{int64(0), true},
	TypeTrivialTestCase{[0]int8{}, true},
	TypeTrivialTestCase{[]int8{}, false},
	TypeTrivialTestCase{struct{}{}, true},
	TypeTrivialTestCase{struct{ int8 }{}, true},
	TypeTrivialTestCase{struct{ a []int8 }{[]int8{}}, false},
	TypeTrivialTestCase{struct{ a [0]int8 }{[0]int8{}}, true},
	TypeTrivialTestCase{(*interface{})(nil), false},
}

func TestIsTypeTrivial(t *testing.T) {
	for _, data := range typeTrivialTestCases {
		assert.Equal(t, data.trivial, IsTypeTrivial(reflect.TypeOf(data.input)))
	}
}
