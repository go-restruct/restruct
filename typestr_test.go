package restruct

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ParseTypeTestCase struct {
	input  string
	typ    reflect.Type
	errstr string
}

var parseTypeTestCases = [...]ParseTypeTestCase{
	// Bad code
	ParseTypeTestCase{"", nil, "parsing error"},
	ParseTypeTestCase{"Invalid", nil, "unknown type Invalid"},
	ParseTypeTestCase{"[][435w5[43]]]**//!!!!!!!", nil, "parsing error"},
	ParseTypeTestCase{"日本語ですか？", nil, "parsing error"},

	// Containers
	ParseTypeTestCase{"float32", reflect.TypeOf(float32(0)), ""},
	ParseTypeTestCase{"*float32", reflect.TypeOf((*float32)(nil)), ""},
	ParseTypeTestCase{"[5]*float32", reflect.TypeOf([5]*float32{}), ""},
	ParseTypeTestCase{"[][][]*float32", reflect.TypeOf([][][]*float32{}), ""},

	// Types
	ParseTypeTestCase{"bool", reflect.TypeOf(false), ""},

	ParseTypeTestCase{"uint8", reflect.TypeOf(uint8(0)), ""},
	ParseTypeTestCase{"uint16", reflect.TypeOf(uint16(0)), ""},
	ParseTypeTestCase{"uint32", reflect.TypeOf(uint32(0)), ""},
	ParseTypeTestCase{"uint64", reflect.TypeOf(uint64(0)), ""},

	ParseTypeTestCase{"int8", reflect.TypeOf(int8(0)), ""},
	ParseTypeTestCase{"int16", reflect.TypeOf(int16(0)), ""},
	ParseTypeTestCase{"int32", reflect.TypeOf(int32(0)), ""},
	ParseTypeTestCase{"int64", reflect.TypeOf(int64(0)), ""},

	ParseTypeTestCase{"complex64", reflect.TypeOf(complex64(0)), ""},
	ParseTypeTestCase{"complex128", reflect.TypeOf(complex128(0)), ""},

	ParseTypeTestCase{"byte", reflect.TypeOf(byte(0)), ""},
	ParseTypeTestCase{"rune", reflect.TypeOf(rune(0)), ""},

	ParseTypeTestCase{"uint", reflect.TypeOf(uint(0)), ""},
	ParseTypeTestCase{"int", reflect.TypeOf(int(0)), ""},
	ParseTypeTestCase{"uintptr", reflect.TypeOf(uintptr(0)), ""},
	ParseTypeTestCase{"string", reflect.TypeOf([]byte{}), ""},

	// Illegal types
	ParseTypeTestCase{"chan int", nil, "channel type not allowed"},
	ParseTypeTestCase{"*chan int", nil, "channel type not allowed"},
	ParseTypeTestCase{"map[string]string", nil, "map type not allowed"},
	ParseTypeTestCase{"map[interface{}]interface{}", nil, "map type not allowed"},

	// Disallowed expressions
	ParseTypeTestCase{"i + 1", nil, "unexpected expression: *ast.BinaryExpr"},
	ParseTypeTestCase{"i()", nil, "unexpected expression: *ast.CallExpr"},
}

func TestParseType(t *testing.T) {
	for _, data := range parseTypeTestCases {
		typ, err := ParseType(data.input)
		if typ != nil {
			assert.Equal(t, data.typ.String(), typ.String())
		}
		if err != nil {
			assert.Equal(t, data.errstr, err.Error())
		}
	}
}
