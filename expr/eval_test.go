package expr

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestStruct1 struct {
	A int
	B float64
}

func TestEvalSimple(t *testing.T) {
	tests := []struct {
		struc  interface{}
		expr   string
		result interface{}
	}{
		{
			TestStruct1{A: 42},
			"A",
			42,
		},
		{
			TestStruct1{A: 42},
			"A * 2",
			84,
		},
		{
			TestStruct1{B: 10.5},
			"B * 2",
			21.0,
		},
		{
			TestStruct1{B: 10.5},
			"-(B * 2)",
			-21.0,
		},
		{
			TestStruct1{A: 0xf0},
			"^0xf0 | A",
			-1,
		},
		{
			TestStruct1{},
			"2 << 2",
			8,
		},
		{
			TestStruct1{},
			"true",
			true,
		},
		{
			TestStruct1{},
			"false",
			false,
		},
		{
			TestStruct1{},
			"true ? 1.0 : 0.0",
			1.0,
		},
		{
			TestStruct1{},
			"false ? 1.0 : 0.0",
			0.0,
		},
		{
			TestStruct1{},
			`"string value!"`,
			"string value!",
		},
		{
			TestStruct1{},
			`"equal" == "equal"`,
			true,
		},
		{
			TestStruct1{},
			`"equal" == "not equal"`,
			false,
		},
		{
			TestStruct1{},
			`"equal" != "not equal"`,
			true,
		},
		{
			TestStruct1{},
			`"equal" != "equal"`,
			false,
		},
		{
			TestStruct1{},
			`"equal"[1] == 'q'`,
			true,
		},
	}

	for _, test := range tests {
		resolver := NewStructResolver(reflect.ValueOf(test.struc))
		result, err := Eval(resolver, test.expr)
		assert.Nil(t, err)
		assert.Equal(t, test.result, result)
	}
}

func TestError(t *testing.T) {
	tests := []struct {
		struc interface{}
		expr  string
		err   string
	}{
		{
			TestStruct1{A: 42},
			"!A",
			"invalid operation: operator ! not defined for 42 (int)",
		},
		{
			TestStruct1{},
			"!42",
			"invalid operation: operator ! not defined for 42 (untyped int constant)",
		},
		{
			TestStruct1{A: 1, B: 1.0},
			"A == B",
			"cannot convert int to float64",
		},
		{
			TestStruct1{A: 1, B: 1.0},
			"A == true",
			"cannot convert int to untyped bool constant",
		},
		{
			TestStruct1{A: 1, B: 1.0},
			"A > true",
			"cannot convert int to untyped bool constant",
		},
	}

	for _, test := range tests {
		resolver := NewStructResolver(reflect.ValueOf(test.struc))
		_, err := Eval(resolver, test.expr)
		assert.EqualError(t, err, test.err)
	}
}
