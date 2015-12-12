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

		// Ignore
		{"-", tagOptions{Ignore: true}, ""},
		{"-,test", tagOptions{}, "extra options on ignored field"},

		// Bad types
		{"invalid", tagOptions{}, "unknown type invalid"},
		{"chan int8", tagOptions{}, "channel type not allowed"},
		{"map[byte]byte", tagOptions{}, "map type not allowed"},

		// Types
		{"uint8", tagOptions{Type: reflect.TypeOf(uint8(0))}, ""},
		{"uint16", tagOptions{Type: reflect.TypeOf(uint16(0))}, ""},
		{"uint32", tagOptions{Type: reflect.TypeOf(uint32(0))}, ""},
		{"int8", tagOptions{Type: reflect.TypeOf(int8(0))}, ""},
		{"int16", tagOptions{Type: reflect.TypeOf(int16(0))}, ""},
		{"int32", tagOptions{Type: reflect.TypeOf(int32(0))}, ""},

		// Sizeof
		{"sizeof=OtherField", tagOptions{SizeOf: "OtherField"}, ""},
		{"sizeof=日本", tagOptions{SizeOf: "日本"}, ""},

		// Skip
		{"skip=4", tagOptions{Skip: 4}, ""},
		{"skip=字", tagOptions{}, "bad skip amount"},

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
