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
		opts   TagOptions
		errstr string
	}{
		// Blank
		{"", TagOptions{}, ""},

		// Gibberish
		{"!#%,1245^#df,little,&~@~@~@~@", TagOptions{}, "parsing error"},
		{"&paREZysLu&@,83D9I!,9OsQ56BLD", TagOptions{}, "parsing error"},
		{"B7~,H0IDSxDlJ,#xa$kgDEL%Ts,88", TagOptions{}, "parsing error"},
		{"fio&eQ8xwbhAWR*!CRlL2XBDG$45s", TagOptions{}, "parsing error"},
		{"IcPyRJ#EV@a4QAb9wENk4Zq9MpX$p", TagOptions{}, "parsing error"},

		// Conflicting byte order
		{"little,big", TagOptions{Order: binary.BigEndian}, ""},
		{"big,little", TagOptions{Order: binary.LittleEndian}, ""},

		// Byte order
		{"msb", TagOptions{Order: binary.BigEndian}, ""},
		{"lsb", TagOptions{Order: binary.LittleEndian}, ""},
		{"network", TagOptions{Order: binary.BigEndian}, ""},

		// Ignore
		{"-", TagOptions{Ignore: true}, ""},
		{"-,test", TagOptions{}, "extra options on ignored field"},

		// Bad types
		{"invalid", TagOptions{}, "unknown type invalid"},
		{"chan int8", TagOptions{}, "channel type not allowed"},
		{"map[byte]byte", TagOptions{}, "map type not allowed"},

		// Types
		{"uint8", TagOptions{Type: reflect.TypeOf(uint8(0))}, ""},
		{"uint16", TagOptions{Type: reflect.TypeOf(uint16(0))}, ""},
		{"uint32", TagOptions{Type: reflect.TypeOf(uint32(0))}, ""},
		{"int8", TagOptions{Type: reflect.TypeOf(int8(0))}, ""},
		{"int16", TagOptions{Type: reflect.TypeOf(int16(0))}, ""},
		{"int32", TagOptions{Type: reflect.TypeOf(int32(0))}, ""},

		// Sizeof
		{"sizeof=OtherField", TagOptions{SizeOf: "OtherField"}, ""},
		{"sizeof=日本", TagOptions{SizeOf: "日本"}, ""},

		// Skip
		{"skip=4", TagOptions{Skip: 4}, ""},
		{"skip=字", TagOptions{}, "bad skip amount"},

		// Composite
		{"uint16,little,sizeof=test,skip=5",
			TagOptions{
				Type:   reflect.TypeOf(uint16(0)),
				Order:  binary.LittleEndian,
				SizeOf: "test",
				Skip:   5,
			},
			"",
		},
	}

	for _, test := range tests {
		opts, err := ParseTag(test.input)
		assert.Equal(t, test.opts, opts)
		if err != nil {
			assert.Equal(t, test.errstr, err.Error())
		}
	}
}

func TestMustParseTagPanicsOnError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Invalid tag did not panic.")
		}
	}()
	MustParseTag("???")
}

func TestMustParseTagReturnsOnSuccess(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Error("Valid tag panicked.")
		}
	}()
	MustParseTag("[128]byte,little,sizeof=Test")
}
