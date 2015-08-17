package restruct

import (
	"encoding/binary"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ParseTagTestCase struct {
	input  string
	opts   TagOptions
	errstr string
}

var parseTagTestCases = [...]ParseTagTestCase{
	// Blank
	ParseTagTestCase{"", TagOptions{}, ""},

	// Gibberish
	ParseTagTestCase{"!#%,1245^#df,little,&~@~@~@~@", TagOptions{}, "parsing error"},
	ParseTagTestCase{"&paREZysLu&@,83D9I!,9OsQ56BLD", TagOptions{}, "parsing error"},
	ParseTagTestCase{"B7~,H0IDSxDlJ,#xa$kgDEL%Ts,88", TagOptions{}, "parsing error"},
	ParseTagTestCase{"fio&eQ8xwbhAWR*!CRlL2XBDG$45s", TagOptions{}, "parsing error"},
	ParseTagTestCase{"IcPyRJ#EV@a4QAb9wENk4Zq9MpX$p", TagOptions{}, "parsing error"},

	// Conflicting byte order
	ParseTagTestCase{"little,big", TagOptions{Order: binary.BigEndian}, ""},
	ParseTagTestCase{"big,little", TagOptions{Order: binary.LittleEndian}, ""},

	// Byte order
	ParseTagTestCase{"msb", TagOptions{Order: binary.BigEndian}, ""},
	ParseTagTestCase{"lsb", TagOptions{Order: binary.LittleEndian}, ""},
	ParseTagTestCase{"network", TagOptions{Order: binary.BigEndian}, ""},

	// Ignore
	ParseTagTestCase{"-", TagOptions{Ignore: true}, ""},
	ParseTagTestCase{"-,test", TagOptions{}, "extra options on ignored field"},

	// Bad types
	ParseTagTestCase{"invalid", TagOptions{}, "unknown type invalid"},
	ParseTagTestCase{"chan int8", TagOptions{}, "channel type not allowed"},
	ParseTagTestCase{"map[byte]byte", TagOptions{}, "map type not allowed"},

	// Types
	ParseTagTestCase{"uint8", TagOptions{Type: reflect.TypeOf(uint8(0))}, ""},
	ParseTagTestCase{"uint16", TagOptions{Type: reflect.TypeOf(uint16(0))}, ""},
	ParseTagTestCase{"uint32", TagOptions{Type: reflect.TypeOf(uint32(0))}, ""},
	ParseTagTestCase{"int8", TagOptions{Type: reflect.TypeOf(int8(0))}, ""},
	ParseTagTestCase{"int16", TagOptions{Type: reflect.TypeOf(int16(0))}, ""},
	ParseTagTestCase{"int32", TagOptions{Type: reflect.TypeOf(int32(0))}, ""},

	// Sizeof
	ParseTagTestCase{"sizeof=OtherField", TagOptions{SizeOf: "OtherField"}, ""},
	ParseTagTestCase{"sizeof=日本", TagOptions{SizeOf: "日本"}, ""},

	// Composite
	ParseTagTestCase{"uint16,little,sizeof=test",
		TagOptions{
			Type:   reflect.TypeOf(uint16(0)),
			Order:  binary.LittleEndian,
			SizeOf: "test",
		},
		"",
	},
}

func TestParseTag(t *testing.T) {
	for _, data := range parseTagTestCases {
		opts, err := ParseTag(data.input)
		assert.Equal(t, data.opts, opts)
		if err != nil {
			assert.Equal(t, data.errstr, err.Error())
		}
	}
}
