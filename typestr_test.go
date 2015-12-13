package restruct

import (
	"go/ast"
	"go/token"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ParseTypeTestCase struct {
	input  string
	typ    reflect.Type
	errstr string
}

func TestParseType(t *testing.T) {
	RegisterArrayType([5]*float32{})

	tests := []struct {
		input  string
		typ    reflect.Type
		errstr string
	}{
		// Bad code
		{"", nil, "parsing error"},
		{"Invalid", nil, "unknown type Invalid"},
		{"[][435w5[43]]]**//!!!!!!!", nil, "parsing error"},
		{"日本語ですか？", nil, "parsing error"},

		// Containers
		{"float32", reflect.TypeOf(float32(0)), ""},
		{"*float32", reflect.TypeOf((*float32)(nil)), ""},
		{"[5]*float32", reflect.TypeOf([5]*float32{}), ""},
		{"[5]*invalid", nil, "unknown type invalid"},
		{"[][][]*float32", reflect.TypeOf([][][]*float32{}), ""},
		{"[][][]*invalid", nil, "unknown type invalid"},

		// Types
		{"bool", reflect.TypeOf(false), ""},

		{"uint8", reflect.TypeOf(uint8(0)), ""},
		{"uint16", reflect.TypeOf(uint16(0)), ""},
		{"uint32", reflect.TypeOf(uint32(0)), ""},
		{"uint64", reflect.TypeOf(uint64(0)), ""},

		{"int8", reflect.TypeOf(int8(0)), ""},
		{"int16", reflect.TypeOf(int16(0)), ""},
		{"int32", reflect.TypeOf(int32(0)), ""},
		{"int64", reflect.TypeOf(int64(0)), ""},

		{"complex64", reflect.TypeOf(complex64(0)), ""},
		{"complex128", reflect.TypeOf(complex128(0)), ""},

		{"byte", reflect.TypeOf(byte(0)), ""},
		{"rune", reflect.TypeOf(rune(0)), ""},

		{"uint", reflect.TypeOf(uint(0)), ""},
		{"int", reflect.TypeOf(int(0)), ""},
		{"uintptr", reflect.TypeOf(uintptr(0)), ""},
		{"string", reflect.TypeOf([]byte{}), ""},

		// Illegal types
		{"chan int", nil, "channel type not allowed"},
		{"*chan int", nil, "channel type not allowed"},
		{"map[string]string", nil, "map type not allowed"},
		{"map[interface{}]interface{}", nil, "map type not allowed"},

		// Disallowed expressions
		{"i + 1", nil, "unexpected expression: *ast.BinaryExpr"},
		{"i()", nil, "unexpected expression: *ast.CallExpr"},
	}

	for _, test := range tests {
		typ, err := parseType(test.input)
		if typ != nil {
			assert.Equal(t, test.typ.String(), typ.String())
		}
		if err != nil {
			assert.Equal(t, test.errstr, err.Error())
		}
	}
}

func TestBadAst(t *testing.T) {
	// typeOfExpr should gracefully handle broken AST structures. Let's
	// construct some.

	// Array with bad length descriptor.
	// [Bad]int32
	badArr := ast.ArrayType{
		Len: ast.NewIdent("Bad"),
		Elt: ast.NewIdent("int32"),
	}
	typ, err := typeOfExpr(&badArr)
	assert.Equal(t, typ, nil)
	assert.Equal(t, err.Error(), "invalid array size expression")

	// Array with bad length descriptor.
	// ["How about that!"]int32
	badArr = ast.ArrayType{
		Len: &ast.BasicLit{Kind: token.STRING, Value: `"How about that!"`},
		Elt: ast.NewIdent("int32"),
	}
	typ, err = typeOfExpr(&badArr)
	assert.Equal(t, typ, nil)
	assert.Equal(t, err.Error(), "invalid array size type")

	// Array with bad length descriptor.
	// [10ii0]int32
	badArr = ast.ArrayType{
		Len: &ast.BasicLit{Kind: token.INT, Value: "10ii0"},
		Elt: ast.NewIdent("int32"),
	}
	typ, err = typeOfExpr(&badArr)
	assert.Equal(t, typ, nil)
	assert.Equal(t, err.Error(), "strconv.ParseInt: parsing \"10ii0\": invalid syntax")
}
