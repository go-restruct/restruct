package restruct

import (
	"encoding/binary"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTag(t *testing.T) {
	tests := []struct {
		input  string
		opts   tagOptions
		errstr string
	}{
		// Blank
		{"", tagOptions{}, ""},

		// Gibberish
		{"!#%,1245^#df,little,&~@~@~@~@", tagOptions{}, "parsing error"},
		{"&paREZysLu&@,83D9I!,9OsQ56BLD", tagOptions{}, "parsing error"},
		{"B7~,H0IDSxDlJ,#xa$kgDEL%Ts,88", tagOptions{}, "parsing error"},
		{"fio&eQ8xwbhAWR*!CRlL2XBDG$45s", tagOptions{}, "parsing error"},
		{"IcPyRJ#EV@a4QAb9wENk4Zq9MpX$p", tagOptions{}, "parsing error"},

		// Conflicting byte order
		{"little,big", tagOptions{Order: binary.BigEndian}, ""},
		{"big,little", tagOptions{Order: binary.LittleEndian}, ""},

		// Byte order
		{"msb", tagOptions{Order: binary.BigEndian}, ""},
		{"lsb", tagOptions{Order: binary.LittleEndian}, ""},
		{"network", tagOptions{Order: binary.BigEndian}, ""},
		{"big little", tagOptions{}, "tag: expected comma"},

		// Ignore
		{"-", tagOptions{Ignore: true}, ""},
		{"-,test", tagOptions{}, "extra options on ignored field"},

		// Bad types
		{"invalid", tagOptions{}, "unknown type invalid"},
		{"chan int8", tagOptions{}, "channel type not allowed"},
		{"map[byte]byte", tagOptions{}, "map type not allowed"},
		{`map[`, tagOptions{}, "struct type: unexpected eof in expr"},
		{`map[]`, tagOptions{}, "struct type: parsing error"},

		// Types
		{"uint8", tagOptions{Type: reflect.TypeOf(uint8(0))}, ""},
		{"uint16", tagOptions{Type: reflect.TypeOf(uint16(0))}, ""},
		{"uint32", tagOptions{Type: reflect.TypeOf(uint32(0))}, ""},
		{"int8", tagOptions{Type: reflect.TypeOf(int8(0))}, ""},
		{"int16", tagOptions{Type: reflect.TypeOf(int16(0))}, ""},
		{"int32", tagOptions{Type: reflect.TypeOf(int32(0))}, ""},

		// Bitfields
		{"uint8:3", tagOptions{Type: reflect.TypeOf(uint8(0)), BitSize: 3}, ""},
		{"uint16:15", tagOptions{Type: reflect.TypeOf(uint16(0)), BitSize: 15}, ""},
		{"uint32:31", tagOptions{Type: reflect.TypeOf(uint32(0)), BitSize: 31}, ""},
		{"uint64:63", tagOptions{Type: reflect.TypeOf(uint64(0)), BitSize: 63}, ""},
		{"uint8:1", tagOptions{Type: reflect.TypeOf(uint8(0)), BitSize: 1}, ""},
		{"uint16:1", tagOptions{Type: reflect.TypeOf(uint16(0)), BitSize: 1}, ""},
		{"uint32:1", tagOptions{Type: reflect.TypeOf(uint32(0)), BitSize: 1}, ""},
		{"uint64:1", tagOptions{Type: reflect.TypeOf(uint64(0)), BitSize: 1}, ""},
		{"[]uint8:1", tagOptions{}, "struct type bits specified on non-bitwise type []uint8"},

		// Wrong bitfields
		{"uint8:0", tagOptions{}, "bit size 0 out of range (1 to 7)"},
		{"uint16:0", tagOptions{}, "bit size 0 out of range (1 to 15)"},
		{"uint32:0", tagOptions{}, "bit size 0 out of range (1 to 31)"},
		{"uint64:0", tagOptions{}, "bit size 0 out of range (1 to 63)"},
		{"int8:0", tagOptions{}, "bit size 0 out of range (1 to 7)"},
		{"int16:0", tagOptions{}, "bit size 0 out of range (1 to 15)"},
		{"int32:0", tagOptions{}, "bit size 0 out of range (1 to 31)"},
		{"int64:0", tagOptions{}, "bit size 0 out of range (1 to 63)"},

		{"uint8:8", tagOptions{BitSize: 0}, "bit size 8 out of range (1 to 7)"},
		{"uint16:16", tagOptions{BitSize: 0}, "bit size 16 out of range (1 to 15)"},
		{"uint32:32", tagOptions{BitSize: 0}, "bit size 32 out of range (1 to 31)"},
		{"uint64:64", tagOptions{BitSize: 0}, "bit size 64 out of range (1 to 63)"},
		{"int8:8", tagOptions{BitSize: 0}, "bit size 8 out of range (1 to 7)"},
		{"int16:16", tagOptions{BitSize: 0}, "bit size 16 out of range (1 to 15)"},
		{"int32:32", tagOptions{BitSize: 0}, "bit size 32 out of range (1 to 31)"},
		{"int64:64", tagOptions{BitSize: 0}, "bit size 64 out of range (1 to 63)"},

		{"uint8:XX", tagOptions{BitSize: 0}, "struct type bits: invalid integer syntax"},
		{"uint16:XX", tagOptions{BitSize: 0}, "struct type bits: invalid integer syntax"},
		{"uint32:XX", tagOptions{BitSize: 0}, "struct type bits: invalid integer syntax"},
		{"uint64:XX", tagOptions{BitSize: 0}, "struct type bits: invalid integer syntax"},
		{"int8:XX", tagOptions{BitSize: 0}, "struct type bits: invalid integer syntax"},
		{"int16:XX", tagOptions{BitSize: 0}, "struct type bits: invalid integer syntax"},
		{"int32:XX", tagOptions{BitSize: 0}, "struct type bits: invalid integer syntax"},
		{"int64:XX", tagOptions{BitSize: 0}, "struct type bits: invalid integer syntax"},

		{"uint8:X:X", tagOptions{BitSize: 0}, "struct type bits: invalid integer syntax"},
		{"uint16:X:X", tagOptions{BitSize: 0}, "struct type bits: invalid integer syntax"},
		{"uint32:X:X", tagOptions{BitSize: 0}, "struct type bits: invalid integer syntax"},
		{"uint64:X:X", tagOptions{BitSize: 0}, "struct type bits: invalid integer syntax"},
		{"int8:X:X", tagOptions{BitSize: 0}, "struct type bits: invalid integer syntax"},
		{"int16:X:X", tagOptions{BitSize: 0}, "struct type bits: invalid integer syntax"},
		{"int32:X:X", tagOptions{BitSize: 0}, "struct type bits: invalid integer syntax"},
		{"int64:X:X", tagOptions{BitSize: 0}, "struct type bits: invalid integer syntax"},

		// Sizeof
		{"sizefrom=OtherField", tagOptions{SizeFrom: "OtherField"}, ""},
		{"sizefrom=日本", tagOptions{SizeFrom: "日本"}, ""},
		{"sizefrom=日本,variantbool", tagOptions{SizeFrom: "日本", VariantBoolFlag: true}, ""},
		{"sizefrom=0", tagOptions{}, "sizefrom: invalid identifier character 0"},

		// Sizeof
		{"sizeof=OtherField", tagOptions{SizeOf: "OtherField"}, ""},
		{"sizeof=日本", tagOptions{SizeOf: "日本"}, ""},
		{"sizeof=日本,variantbool", tagOptions{SizeOf: "日本", VariantBoolFlag: true}, ""},
		{"sizeof=0", tagOptions{}, "sizeof: invalid identifier character 0"},

		// Skip
		{"skip=4", tagOptions{Skip: 4}, ""},
		{"skip=字", tagOptions{}, "skip: invalid integer character 字"},

		// Expressions
		{"if=true", tagOptions{IfExpr: "true"}, ""},
		{"if=call(0,1)", tagOptions{IfExpr: "call(0,1)"}, ""},
		{`if=call(")Test)"),sizeof=Test`, tagOptions{IfExpr: `call(")Test)")`, SizeOf: "Test"}, ""},
		{"size=4", tagOptions{SizeExpr: "4"}, ""},
		{"size={,", tagOptions{}, "size: unexpected eof in expr"},
		{"bits=4", tagOptions{BitsExpr: "4"}, ""},
		{"bits={,", tagOptions{}, "bits: unexpected eof in expr"},
		{"in=42", tagOptions{InExpr: "42"}, ""},
		{"out=struct{}{}", tagOptions{OutExpr: "struct{}{}"}, ""},
		{"out=struct{}{},variantbool", tagOptions{OutExpr: "struct{}{}", VariantBoolFlag: true}, ""},
		{"while=true", tagOptions{WhileExpr: "true"}, ""},
		{`if="`, tagOptions{}, "if: unexpected eof in literal"},
		{`while="\"`, tagOptions{}, "while: unexpected eof in literal"},
		{`in="\"\"""`, tagOptions{}, "in: unexpected eof in literal"},
		{`out="\"test`, tagOptions{}, "out: unexpected eof in literal"},

		// Root
		{"root", tagOptions{RootFlag: true}, ""},

		// Parent
		{"parent", tagOptions{ParentFlag: true}, ""},

		// Composite
		{"uint16,little,sizeof=test,skip=5",
			tagOptions{
				Type:   reflect.TypeOf(uint16(0)),
				Order:  binary.LittleEndian,
				SizeOf: "test",
				Skip:   5,
			},
			"",
		},
	}

	for _, test := range tests {
		opts, err := parseTag(test.input)
		assert.Equal(t, test.opts, opts)
		if err != nil {
			assert.NotEmpty(t, test.errstr)
			assert.Contains(t, err.Error(), test.errstr)
		}
	}
}

func TestMustParseTagPanicsOnError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Invalid tag did not panic.")
		}
	}()
	mustParseTag("???")
}

func TestMustParseTagReturnsOnSuccess(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Error("Valid tag panicked.")
		}
	}()
	mustParseTag("[128]byte,little,sizeof=Test")
}
